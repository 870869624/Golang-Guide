package app

import (
	zerologx "git.cloud.top/cbb/infrastructure/zero/logx"
	"git.cloud.top/go/go-zero/core/logx"

	"git.cloud.top/aiop/common/util"

	"git.cloud.top/ngedr/server/rest/api/config"
	rest "git.cloud.top/ngedr/server/rest/app"
	rpc "git.cloud.top/ngedr/server/rpc/app"
)

func New(envPrefix, configFile string) (*rpc.NgEDR, *rest.NgEDR) {
	envPrefix = util.EnvPrefix(envPrefix)
	conf := config.Load(envPrefix, configFile)
	conf.RestConf.EnvPrefix = envPrefix
	logxConf := conf.RestConf.Logger.Parse()
	zeroLogx := zerologx.Logx(logxConf)
	writer, err := zeroLogx.NewWriter(conf.RestConf.Logger.Stacktrace)
	logx.Must(err)
	logx.MustSetup(logxConf)
	logx.SetLevel(zerologx.Level(conf.RestConf.Logger.Level))
	logx.SetWriter(writer)
	return rpc.New(configFile, conf.RestConf), rest.New(conf)
}
