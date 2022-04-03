package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/luhaopei/mxshop_api/user-web/global"
	"github.com/luhaopei/mxshop_api/user-web/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	_ "github.com/mbobakov/grpc-consul-resolver"
)

func InitSrvConn() {
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			global.ServerConfig.ConsulInfo.Host,
			global.ServerConfig.ConsulInfo.Port,
			global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InittSrvConn] 连接 【用户服务失败】")
	}
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}

func InitSrvConn2() {
	zap.S().Debug("获取用户列表页")
	// 从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	zap.S().Infof("cfg.Address: %s", cfg.Address)
	userSrvHost := ""
	userSrvPort := 0

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	
	data, err := client.Agent().Services()
	if err != nil {
		panic(err)
	}
	zap.S().Infof("data size: %d", len(data))
	for key, _ := range data {
		zap.S().Info(key)
	}

	filter := fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSrvInfo.Name)
	zap.S().Info(filter)
	data, err = client.Agent().ServicesWithFilter(filter)
	if err != nil {
		panic(err)
	}

	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == "" {
		zap.S().Fatal("[InittSrvConn] 连接 【用户服务失败】")
		return
	}
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost,
		userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList]连接【用户服务失败】",
			"msg", err.Error(),
		)
	}
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}
