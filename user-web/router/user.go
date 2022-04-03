package router

import (
	"github.com/gin-gonic/gin"
	"github.com/luhaopei/mxshop_api/user-web/api"
)

func InitUserRouter(router *gin.RouterGroup) {
	UserRouter := router.Group("user")
	{
		UserRouter.GET("list", api.GetUserList)
	}

}
