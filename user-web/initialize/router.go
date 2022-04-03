package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/luhaopei/mxshop_api/user-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	ApiGrop := Router.Group("/v1")
	router.InitUserRouter(ApiGrop)

	return Router
}
