package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/luhaopei/mxshop_api/user-web/middlewares"
	"github.com/luhaopei/mxshop_api/user-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())
	ApiGrop := Router.Group("/v1")
	router.InitUserRouter(ApiGrop)
	router.InitBaseRouter(ApiGrop)

	return Router
}
