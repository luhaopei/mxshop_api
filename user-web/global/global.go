package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/luhaopei/mxshop_api/user-web/config"
	"github.com/luhaopei/mxshop_api/user-web/proto"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	Trans        ut.Translator

	UserSrvClient proto.UserClient
)
