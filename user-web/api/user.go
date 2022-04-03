package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"github.com/luhaopei/mxshop_api/user-web/middlewares"
	"github.com/luhaopei/mxshop_api/user-web/models"

	"github.com/luhaopei/mxshop_api/user-web/forms"
	"github.com/luhaopei/mxshop_api/user-web/global"
	"github.com/luhaopei/mxshop_api/user-web/global/response"
	"github.com/luhaopei/mxshop_api/user-web/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandlerGrpcErrorToHttp(err error, c *gin.Context) {
	if err == nil {
		return
	}
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"msg": e.Message(),
			})
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "内部错误",
			})
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "参数错误",
			})
		case codes.Unavailable:
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "用户服务不可用",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "其他错误" + string(e.Message()),
			})
		}
		return
	}
}
func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func GetUserList(ctx *gin.Context) {

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户：%d", currentUser.ID)

	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)

	pSize := ctx.DefaultQuery("psize", "10")
	pnSize, _ := strconv.Atoi(pSize)

	resp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pnSize),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList]查询【用户列表】失败")
		HandlerGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range resp.Data {
		user := response.UserResponse{
			value.Id,
			value.NickName,
			response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			value.Gender,
			value.Mobile,
		}

		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}

func PassWordLogin(c *gin.Context) {
	// 表单验证
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList]连接【用户服务失败】",
			"msg", err.Error(),
		)
	}
	userSrvClient := proto.NewUserClient(userConn)

	resp, err := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: passwordLoginForm.Mobile})
	if err != nil {
		zap.S().Errorw("[PassWordLogin]查询【用户列表】失败")
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
			}
		}
		return
	}
	passRsp, err := userSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{Password: passwordLoginForm.PassWord,
		EncryptedPasword: resp.Password})
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"password": "登录失败",
		})
	}
	if passRsp.Success {
		j := middlewares.NewJWT()
		claims := models.CustomClaims{
			ID:          uint(resp.Id),
			NickName:    resp.NickName,
			AuthorityId: uint(resp.Role),
			StandardClaims: jwt.StandardClaims{NotBefore: time.Now().Unix(),
				ExpiresAt: time.Now().Unix() + 60*60*24*30,
				Issuer:    "imooc"},
		}
		token, err := j.CreateToken(claims)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"msg": "生成token失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":         resp.Id,
			"nick_name":  resp.NickName,
			"token":      token,
			"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
		})
	} else {
		c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "登录失败",
		})
	}

}

func Register(c *gin.Context) {
	// 用户注册
	registerForm := forms.RegisterForm{}

	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	// 验证码校验
	rdb := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port)})
	value, err := rdb.Get(registerForm.Mobile).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
		return
	} else {
		if value != registerForm.Code {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": "验证码错误",
			})
			return
		}
	}

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList]连接【用户服务失败】",
			"msg", err.Error(),
		)
	}
	userSrvClient := proto.NewUserClient(userConn)

	user, err := userSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{NickName: registerForm.Mobile,
		PassWord: registerForm.PassWord,
		Mobile:   registerForm.Mobile})

	if err != nil {
		zap.S().Errorf("[Register]查询 【新建用户失败】失败:%s", err.Error())
		HandlerGrpcErrorToHttp(err, c)
		return
	}

	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + 60*60*24*30,
			Issuer:    "imooc"},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"msg": "生成token失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}
