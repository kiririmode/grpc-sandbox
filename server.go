package main

import (
	"github.com/kiririmode/grpc-sandbox/common"
	"github.com/kiririmode/grpc-sandbox/common/conf"
	"github.com/kiririmode/grpc-sandbox/common/log"
	"github.com/kiririmode/grpc-sandbox/router"
	"google.golang.org/grpc"
)

func main() {
	s := grpc.NewServer()

	// リソースの準備
	config := conf.NewConfiguration("stubserver", "development", []string{"conf"})
	logr := log.NewLog(config)
	server := router.NewGrpcServer(s, config, logr)

	// リソースの開始・終了処理
	rm := common.NewResourceManager([]common.Resource{config, logr, server})
	rm.Initialize()
	defer rm.Finalize()

	logr.Logger.Info("initialization succeeds")
	err := server.Serve()
	if err != nil {
		logr.Logger.Fatalf("failed to serve %s", err)
	}
}
