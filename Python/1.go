package queue

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"git.cloud.top/aiop/ersvc/rpc/ersvc"
	"git.cloud.top/cbb/infrastructure/queue/message"
	"git.cloud.top/go/go-zero/core/logx"
	"git.cloud.top/ngedr/server/pkg/model"
	"git.cloud.top/ngedr/server/pkg/model/threatresponse"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TLConsumer struct {
	model *model.Model
	erSvc ersvc.Ersvc
}

type VirusRecord struct {
	AgentID         string    `bson:"agent_id"`
	MD5             string    `bson:"md5"`
	FilePath        string    `bson:"file_path"`
	LoggedAt        time.Time `bson:"logged_at"`
	UpdatedAt       time.Time `bson:"updated_at"`
	Result          string    `bson:"result"`
	IsolationResult string    `bson:"isolation_result"`
}

// nolint
func (l *TLConsumer) Read(messages ...*message.ConsumerMessage) {
	logx.Infof("Reading messages...")
	a := time.Now()
	virusTypeMap, err := l.getVirusTypeTranslation()
	if err != nil {
		logx.Errorf("Failed to get virus type translation: %v", err)
		return
	}

	var bulkModels []mongo.WriteModel

	for _, msg := range messages {
		logs := make(map[string]interface{})
		if err := json.Unmarshal(msg.Value, &logs); err != nil {
			logx.Errorf("client log consumer unmarshal failed: %v", err)
			continue
		}
		logx.Infof("client log: %v", logs)

		loggedAt, ok := logs["logged_at"].(float64)
		if !ok {
			logx.Errorf("Invalid logged_at value: %v", logs["logged_at"])
			continue
		}

		logType, ok := logs["log_type"].(string)
		if !ok {
			logx.Errorf("Invalid log_type value: %v", logs["log_type"])
			continue
		}

		detailTypeFloat, ok := logs["detail_type"].(float64)
		if !ok {
			logx.Errorf("Invalid detail_type value: %v", logs["detail_type"])
			continue
		}
		detailType := int(detailTypeFloat)

		agentIP := safeString(logs["agent_ip"])
		agentID := safeString(logs["agent_id"])
		groupID := safeString(logs["group_id"])
		groupName := safeString(logs["group_name"])
		hostname := safeString(logs["hostname"])
		tenantID := safeString(logs["tenant_id"])
		ruleVersion := safeString(logs["rule_version"])
		agentVersion := safeString(logs["agent_version"])
		principalName := safeString(logs["principal_name"])

		logLevelFloat, _ := logs["log_level"].(float64)
		logLevel := int(logLevelFloat)

		relatedType := safeString(logs["related_type"])
		relatedID := safeString(logs["related_id"])
		relatedName := safeString(logs["related_name"])
		description := safeString(logs["description"])
		virusName := safeString(logs["virus_name"])
		virusType := safeString(logs["virus_type"])
		filePath := safeString(logs["file_path"])
		virusPath := safeString(logs["virus_path"])
		handle := safeString(logs["handle"])
		reason := safeString(logs["reason"])
		result := safeString(logs["result"])
		isolationResult := safeString(logs["isolation_result"])
		md5 := safeString(logs["md5"])
		sha1 := safeString(logs["sha1"])
		sha256 := safeString(logs["sha256"])

		logTypeBelongTo, err := l.getLogTypeBelongTo(logType, tenantID)
		if err != nil {
			logx.Errorf("Failed to get log type belong to: %v", err)
			continue
		}

		if logTypeBelongTo == "virus_defense" && detailType != 2 {
			virusTypeCn := virusTypeMap[virusType]

			properties := map[string]interface{}{
				"logged_at":        int64(loggedAt),
				"agent_ip":         agentIP,
				"agent_id":         agentID,
				"group_id":         groupID,
				"group_name":       groupName,
				"hostname":         hostname,
				"tenant_id":        tenantID,
				"rule_version":     ruleVersion,
				"agent_version":    agentVersion,
				"principal_name":   principalName,
				"log_type":         logType,
				"log_level":        logLevel,
				"related_type":     relatedType,
				"related_id":       relatedID,
				"related_name":     relatedName,
				"detail_type":      detailType,
				"description":      description,
				"virus_name":       virusName,
				"virus_type":       virusType,
				"virus_type_cn":    virusTypeCn,
				"file_path":        filePath,
				"file_name":        filepath.Base(filePath),
				"virus_path":       virusPath,
				"handle":           handle,
				"reason":           reason,
				"result":           result,
				"isolation_result": isolationResult,
				"md5":              md5,
				"sha1":             sha1,
				"sha256":           sha256,
			}

			data := VirusRecord{
				AgentID:         agentID,
				FilePath:        filePath,
				IsolationResult: isolationResult,
				LoggedAt:        time.Now(),
				UpdatedAt:       time.Now(),
				MD5:             md5,
				Result:          result,
			}

			model := mongo.NewUpdateOneModel().
				SetFilter(bson.M{
					"agent_id":  data.AgentID,
					"md5":       data.MD5,
					"file_path": data.FilePath,
				}).
				SetUpdate(bson.M{
					"$set": properties,
				}).
				SetUpsert(true)

			bulkModels = append(bulkModels, model)

			if len(bulkModels) >= 100 {
				err := l.bulkInsertIfNotExists(l.model.ThreatResponseModel, bulkModels)
				if err != nil {
					logx.Errorf("Failed to bulk insert threat response logs: %v", err)
				}
				bulkModels = bulkModels[:0]
			}
		}
	}

	if len(bulkModels) > 0 {
		err := l.bulkInsertIfNotExists(l.model.ThreatResponseModel, bulkModels)
		if err != nil {
			logx.Errorf("Failed to bulk insert threat response logs: %v", err)
		}
	}

	logx.Infof("Threat response logs processed in %v, total records: %d", time.Since(a), len(messages))
}

func safeString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}

func (l *TLConsumer) getVirusTypeTranslation() (map[string]interface{}, error) {
	resp, err := l.erSvc.AuditSearchEnumsTranslate(context.Background(), &ersvc.AuditedSearchEnumsReq{
		Field: "virus_type",
	})
	if err != nil {
		return nil, err
	}
	return resp.GetData().AsMap(), nil
}

func (l *TLConsumer) getLogTypeBelongTo(logType, tenantID string) (string, error) {
	resp, err := l.erSvc.AuditSearchLogTypeBelongTo(context.Background(), &ersvc.AuditedSearchEnumsReq{
		Field: logType,
		Properties: &ersvc.ReservedProperty{
			TenantId: tenantID,
		},
	})
	if err != nil {
		return "", err
	}
	result := resp.GetData().AsMap()
	if val, exists := result[logType]; exists {
		return val.(string), nil
	}
	return "", fmt.Errorf("log type %s not found in response", logType)
}

func (l *TLConsumer) bulkInsertIfNotExists(collection threatresponse.ThreatResponseModel, models []mongo.WriteModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := collection.BulkWrite(ctx, models)
	if err != nil {
		return fmt.Errorf("threat response log failed to bulk insert records: %v", err)
	}

	return nil
}
