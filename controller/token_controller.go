package controller

import (
	"gofile/model/common"
	"gofile/service"

	"github.com/kataras/iris/v12"
)

// 登录
func SignIn(ctx iris.Context) {
	signIn := common.SignIn{}
	resolveParam(ctx, &signIn)
	tokenResult := service.SignIn(signIn)
	ctx.JSON(common.NewSuccessData("登录成功", tokenResult))
}

// 退出登录
func SignOut(ctx iris.Context) {
	tokenResult := common.TokenResult{}
	resolveParam(ctx, &tokenResult)
	service.SignOut(tokenResult)
	ctx.JSON(common.NewSuccess("已退出登录"))
}

// 刷新token
func TokenRefresh(ctx iris.Context) {
	tokenResult := common.TokenResult{}
	resolveParam(ctx, &tokenResult)
	tokenResult = service.TokenRefresh(tokenResult.RefreshToken)
	ctx.JSON(common.NewSuccessData("token刷新成功", tokenResult))
}

// 获取验证码
func Captcha(ctx iris.Context) {
	signIn := common.SignIn{}
	resolveParam(ctx, &signIn)
	captchaResult := service.Captcha(signIn)
	ctx.JSON(common.NewSuccessData("验证码获取成功", captchaResult))
}

// 校验验证码
func CaptchaValidate(ctx iris.Context) {
	signIn := common.SignIn{}
	resolveParam(ctx, &signIn)
	ok := service.CaptchaValidate(signIn)
	if ok {
		ctx.JSON(common.NewSuccessData("验证码校验成功", true))
	} else {
		ctx.JSON(common.NewSuccessData("验证码校验失败", false))
	}
}
