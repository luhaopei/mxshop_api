package api

import "github.com/gin-gonic/gin"

func GetUserList(ctx *gin.Context) {
	zap.S().Debug("获取用户列表页")
	ip := "127.0.0.1"
	port :=50051
	grpc.Dil(fmt.Sprintg("%s:%d", ip , port), grpc.With)

}
