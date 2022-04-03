package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/luhaopei/mxshop_api/user-web/models"
	"net/http"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		curentUser := claims.(*models.CustomClaims)
		if curentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}

}
