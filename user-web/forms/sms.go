package forms

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` //自定义validate
	Type   uint   `form:"type" json:"type" binding:"required,oneof=1 2"`  // 1代表注册，2代表动态验证码登录
}
