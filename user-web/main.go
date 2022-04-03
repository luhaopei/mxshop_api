package main

import (
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/luhaopei/mxshop_api/user-web/global"
	"github.com/luhaopei/mxshop_api/user-web/initialize"
	"github.com/luhaopei/mxshop_api/user-web/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	myvalidator "github.com/luhaopei/mxshop_api/user-web/validator"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	Router := initialize.Routers()
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	initialize.InitSrvConn()

	viper.AutomaticEnv()
	debug := viper.GetBool("MXSHOP_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	// 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	zap.S().Infof("启动服务器，端口：8021")
	if err := Router.Run(":8021"); err != nil {
		zap.S().Panic("启动失败：", err.Error())
	}
}
