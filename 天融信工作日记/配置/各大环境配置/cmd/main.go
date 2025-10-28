package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"git.cloud.top/cbb/infrastructure/version"
	"git.cloud.top/go/go-zero/core/service"

	"git.cloud.top/ngedr/server/cmd/app"
)

func main() {
	var configFile = flag.String("f", "etc/ngedr.yaml", "the config file")
	var versionFlag = flag.Bool("v", false, "show version info")
	flag.Parse()

	version.Print()

	if *versionFlag {
		os.Exit(0)
	}

	services := service.NewServiceGroup()
	defer services.Stop()
	rpc, rest := app.New("", *configFile)
	services.Add(rpc)
	services.Add(rest)

	fmt.Printf("⇨ ngedr start at: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	services.Start()
}
