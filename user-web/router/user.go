package router

import (
	"github.com/gin-gonic/gin"
	"github.com/luhaopei/mxshop_api/user-web/api"
	"github.com/luhaopei/mxshop_api/user-web/middlewares"
)

func InitUserRouter(router *gin.RouterGroup) {
	UserRouter := router.Group("user")
	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)
		UserRouter.POST("register", api.Register)
	}

}
