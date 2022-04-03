package main

import (
	"github.com/luhaopei/mxshop_api/user-web/initialize"
	"go.uber.org/zap"
)

func main() {
	initialize.InitLogger()
	Router := initialize.Routers()

	zap.S().Infof("启动服务器，端口：8021")
	if err := Router.Run(":8021"); err != nil {
		zap.S().Panic("启动失败：", err.Error())
	}
}
