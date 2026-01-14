// 流式调用日志拦截器
package app

import (
	"encoding/json"
	"fmt"
	"time"

	"git.cloud.top/go/go-zero/core/conf"
	"git.cloud.top/go/go-zero/core/logx"
	"git.cloud.top/go/go-zero/core/service"
	"git.cloud.top/go/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"

	rsetconfig "git.cloud.top/aiop/common/rest/config"

	managersvcPb "git.cloud.top/aiop/managersvc/rpc/pb"

	"git.cloud.top/ngdlp/server/rpc/internal/config"
	"git.cloud.top/ngdlp/server/rpc/internal/cronjob"

	dipmanagerServer "git.cloud.top/ngdlp/server/rpc/internal/server/dipmanager"
	ngdlpServer "git.cloud.top/ngdlp/server/rpc/internal/server/ngdlp"
	"git.cloud.top/ngdlp/server/rpc/internal/svc"
	"git.cloud.top/ngdlp/server/rpc/pb"
)

type NgDLP struct {
	envPrefix string
	listenOn  string
	server    *zrpc.RpcServer
}

func New(configFile string, appConfig rsetconfig.Config) *NgDLP {
	ngDLP := &NgDLP{
		envPrefix: appConfig.EnvPrefix,
	}
	var c config.Config
	conf.MustLoad(configFile, &c)
	c.LoadEnv(ngDLP.envPrefix)
	c.Project = appConfig.Project
	c.MongoDB = appConfig.MongoDB
	c.RPCAuth = appConfig.RPCAuth
	c.RPCClient = appConfig.RPCClient
	c.MaxBytes = appConfig.MaxBytes
	c.RPCServerConf.CpuThreshold = 0
	if appConfig.Middlewares.Shedding.Enable {
		c.RPCServerConf.CpuThreshold = appConfig.Middlewares.Shedding.CPUThreshold
	}
	confBytes, _ := json.MarshalIndent(c, "", " ")
	logx.Infof("rpc config: %s", confBytes)
	ctx := svc.NewServiceContext(c)
	s := zrpc.MustNewServer(c.RPCServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterNgdlpServer(grpcServer, ngdlpServer.NewNgdlpServer(ctx))
		managersvcPb.RegisterDipManagersvcServer(grpcServer, dipmanagerServer.NewDipManagersvcServer(ctx))

		if c.RPCServerConf.Mode == service.DevMode || c.RPCServerConf.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	logx.Must(cronjob.RunCroners(ctx))
	s.AddStreamInterceptors(StreamLoggingInterceptor)
	s.AddOptions(grpc.MaxRecvMsgSize(c.MaxBytes))
	ctx.Init(s)
	ngDLP.listenOn = c.RPCServerConf.ListenOn
	ngDLP.server = s
	return ngDLP
}

func (s *NgDLP) Start() {
	fmt.Printf("⇨ rpc server started on %s\n", s.listenOn)
	s.server.Start()
}

func (s *NgDLP) Stop() {
	s.server.Stop()
}

// 重点在此处
func StreamLoggingInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	startTime := time.Now()
	ctx := ss.Context()
	var addr string
	if client, ok := peer.FromContext(ctx); ok {
		addr = client.Addr.String()
	}
	logx.WithContext(ctx).Infof("[RPC] [STREAM-START] %s - %s", addr, info.FullMethod)
	err := handler(srv, ss)
	duration := time.Since(startTime)
	if err != nil {
		logx.WithContext(ctx).Errorf("[RPC] [STREAM-END] %s - %s, duration: %v, error: %v", addr, info.FullMethod, duration, err)
	} else {
		logx.WithContext(ctx).Infof("[RPC] [STREAM-END] %s - %s, duration: %v", addr, info.FullMethod, duration)
	}
	return err
}
