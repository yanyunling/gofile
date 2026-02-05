package service

import (
	"fmt"
	"gofile/middleware"
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/repository"
	"gofile/util"
	"time"

	"github.com/wenlng/go-captcha/v2/slide"
)

const AccessTokenExpire = time.Hour
const RefreshTokenExpire = time.Hour * 24 * 30
const CaptchaExpire = time.Minute * 5

// 登录
func SignIn(signIn common.SignIn) common.TokenResult {
	// 去除用户名的空白
	signIn.Username = util.RemoveBlank(signIn.Username)
	if signIn.Username == "" || signIn.Password == "" {
		panic(common.NewError("用户名或密码不可为空"))
	}

	// 校验登录次数
	checkSignInTimes(signIn.Username)

	// 校验验证码
	ok := CaptchaValidate(signIn)
	if !ok {
		panic(common.NewError("验证失败"))
	}
	// 移除验证码缓存
	middleware.Cache.Delete(common.CaptchaCache + signIn.CaptchaId)

	// 根据用户名查询用户
	userFromDB, err := repository.UserGetByUsername(middleware.Db, signIn.Username)
	if err != nil {
		panic(common.NewErr("用户名或密码错误", err))
	}

	// 验证状态
	if userFromDB.Enable != 1 {
		panic(common.NewError("用户已禁用"))
	}

	// 验证密码
	if !util.CheckBcrypt(userFromDB.Password, signIn.Password) {
		panic(common.NewError("用户名或密码错误"))
	}

	// 管理员用户
	adminUser := entity.GetAdminUser()

	// 生成token
	accessToken := util.RandomString(64)
	refreshToken := util.RandomString(64)

	// 返回数据
	tokenResult := common.TokenResult{}
	tokenResult.AccessToken = accessToken
	tokenResult.RefreshToken = refreshToken
	tokenResult.IsAdmin = adminUser.Id == userFromDB.Id
	tokenResult.Username = userFromDB.Username
	tokenResult.PublicAuth = userFromDB.PublicAuth
	tokenResult.ProtectedAuth = userFromDB.ProtectedAuth
	tokenResult.PrivateAuth = userFromDB.PrivateAuth

	// 缓存数据
	tokenCache := common.TokenCache{}
	tokenCache.AccessToken = accessToken
	tokenCache.RefreshToken = refreshToken
	tokenCache.IsAdmin = adminUser.Id == userFromDB.Id
	tokenCache.UserId = userFromDB.Id
	tokenCache.Username = userFromDB.Username
	tokenCache.PublicAuth = userFromDB.PublicAuth
	tokenCache.ProtectedAuth = userFromDB.ProtectedAuth
	tokenCache.PrivateAuth = userFromDB.PrivateAuth

	// 缓存token
	middleware.Cache.Set(common.AccessTokenCache+tokenResult.AccessToken, &tokenCache, AccessTokenExpire)
	middleware.Cache.Set(common.RefreshTokenCache+tokenResult.RefreshToken, &tokenCache, RefreshTokenExpire)
	middleware.Cache.Delete(common.SignInTimesCache + signIn.Username)

	return tokenResult
}

// 退出登录
func SignOut(tokenResult common.TokenResult) {
	res, found := middleware.Cache.Get(common.RefreshTokenCache + tokenResult.RefreshToken)
	if found {
		tokenCache := res.(*common.TokenCache)
		if tokenCache.RefreshToken != "" {
			middleware.Cache.Delete(common.RefreshTokenCache + tokenCache.RefreshToken)
		}
		if tokenCache.AccessToken != "" {
			middleware.Cache.Delete(common.AccessTokenCache + tokenCache.AccessToken)
		}
	}
}

// 刷新token
func TokenRefresh(refreshToken string) common.TokenResult {
	res, found := middleware.Cache.Get(common.RefreshTokenCache + refreshToken)
	if !found {
		panic(common.NewError("认证信息已过期，请重新登录"))
	}
	tokenCache := res.(*common.TokenCache)
	if tokenCache.RefreshToken == "" {
		panic(common.NewError("认证信息已过期，请重新登录"))
	}

	// 查询用户
	userFromDB, err := repository.UserGetById(middleware.Db, tokenCache.UserId)
	if err != nil {
		panic(common.NewErr("用户加载失败", err))
	}

	// 验证状态
	if userFromDB.Enable != 1 {
		panic(common.NewError("用户已禁用"))
	}

	// 管理员用户
	adminUser := entity.GetAdminUser()

	// 生成token
	accessToken := util.RandomString(64)
	refreshToken = util.RandomString(64)

	// 返回数据
	tokenResult := common.TokenResult{}
	tokenResult.AccessToken = accessToken
	tokenResult.RefreshToken = refreshToken
	tokenResult.IsAdmin = adminUser.Id == userFromDB.Id
	tokenResult.Username = userFromDB.Username
	tokenResult.PublicAuth = userFromDB.PublicAuth
	tokenResult.ProtectedAuth = userFromDB.ProtectedAuth
	tokenResult.PrivateAuth = userFromDB.PrivateAuth

	// 缓存数据
	newTokenCache := common.TokenCache{}
	newTokenCache.AccessToken = accessToken
	newTokenCache.RefreshToken = refreshToken
	newTokenCache.IsAdmin = adminUser.Id == userFromDB.Id
	newTokenCache.UserId = userFromDB.Id
	newTokenCache.Username = userFromDB.Username
	newTokenCache.PublicAuth = userFromDB.PublicAuth
	newTokenCache.ProtectedAuth = userFromDB.ProtectedAuth
	newTokenCache.PrivateAuth = userFromDB.PrivateAuth

	// 缓存token
	middleware.Cache.Set(common.AccessTokenCache+newTokenCache.AccessToken, &newTokenCache, AccessTokenExpire)
	middleware.Cache.Set(common.RefreshTokenCache+newTokenCache.RefreshToken, &newTokenCache, RefreshTokenExpire)

	return tokenResult
}

// 获取验证码
func Captcha(signIn common.SignIn) middleware.CaptchaResult {
	// 传入的id是上一次的验证码，此操作为刷新验证码，移除上一次的缓存
	if signIn.CaptchaId != "" {
		middleware.Cache.Delete(common.CaptchaCache + signIn.CaptchaId)
	}

	// 生成验证码
	captchaResult, err := middleware.CaptchaGet()
	captchaResult.CaptchaId = util.UUID()
	if err != nil {
		panic(common.NewErr("验证码生成失败", err))
	}

	// 缓存验证信息
	captchaCache := middleware.CaptchaCache{X: captchaResult.X, Y: captchaResult.Y}
	middleware.Cache.Set(common.CaptchaCache+captchaResult.CaptchaId, &captchaCache, CaptchaExpire)

	return captchaResult
}

// 校验验证码
func CaptchaValidate(signIn common.SignIn) bool {
	// 根据验证码id从缓存获取验证码信息
	res, found := middleware.Cache.Get(common.CaptchaCache + signIn.CaptchaId)
	if !found {
		return false
	}
	captchaCache := res.(*middleware.CaptchaCache)

	// 校验
	ok := slide.Validate(signIn.CaptchaX, signIn.CaptchaY, captchaCache.X, captchaCache.Y, 6)

	// 移除验证码缓存
	if !ok {
		middleware.Cache.Delete(common.CaptchaCache + signIn.CaptchaId)
	}

	return ok
}

// 校验登录次数，如已超出则抛出异常
func checkSignInTimes(name string) {
	times := 1
	signInTimes, exp, found := middleware.Cache.GetWithExpiration(common.SignInTimesCache + name)
	now := time.Now()
	if !found || exp.Before(now) {
		middleware.Cache.Set(common.SignInTimesCache+name, times, time.Minute*5)
	} else {
		times = signInTimes.(int) + 1
		middleware.Cache.Set(common.SignInTimesCache+name, times, exp.Sub(now))
	}
	if times > 5 {
		panic(common.NewError(fmt.Sprintf("登录次数已达上限，请于%v分钟后再试", int(exp.Sub(now).Minutes())+1)))
	}
}
