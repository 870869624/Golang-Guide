```go
import (
	"context"
	_ "embed"
	"encoding/json"
	"git.cloud.top/aiop/cleansvc/rpc/internal/model/system"
	"git.cloud.top/go/go-zero/core/logx"
	"git.cloud.top/go/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//go:embed ./soap_config.conf
var soapConfFile []byte

func (m *customSoapConfModel) InitSoapConfInfo() (*SoapConf, error) {
	var result SoapConf
	err := json.Unmarshal(soapConfFile, &result)
	if err != nil {
		logx.Errorf("json unmarshal failed.check ag! err=%s\r\n", err)
		return nil, err
	}
	return &result, nil
}
```

可以用这种方式绑定文件
