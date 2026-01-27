package service

import (
	"fmt"
	"gofile/middleware"
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/repository"
	"gofile/util"
	"time"

	"github.com/muesli/cache2go"
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
	cache2go.Cache(common.CaptchaCache).Delete(signIn.CaptchaId)

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
	tokenResult.NickName = userFromDB.NickName
	tokenResult.HasUpdate = userFromDB.HasUpdate

	// 缓存数据
	tokenCache := common.TokenCache{}
	tokenCache.AccessToken = accessToken
	tokenCache.RefreshToken = refreshToken
	tokenCache.IsAdmin = adminUser.Id == userFromDB.Id
	tokenCache.UserId = userFromDB.Id
	tokenCache.Username = userFromDB.Username
	tokenCache.HasUpdate = userFromDB.HasUpdate

	// 缓存token
	cache2go.Cache(common.AccessTokenCache).Add(tokenResult.AccessToken, AccessTokenExpire, &tokenCache)
	cache2go.Cache(common.RefreshTokenCache).Add(tokenResult.RefreshToken, RefreshTokenExpire, &tokenCache)
	cache2go.Cache(common.SignInTimesCache).Delete(signIn.Username)

	return tokenResult
}

// 退出登录
func SignOut(tokenResult common.TokenResult) {
	res, err := cache2go.Cache(common.RefreshTokenCache).Value(tokenResult.RefreshToken)
	if err == nil {
		tokenCache := res.Data().(*common.TokenCache)
		if tokenCache.RefreshToken != "" {
			cache2go.Cache(common.RefreshTokenCache).Delete(tokenCache.RefreshToken)
		}
		if tokenCache.AccessToken != "" {
			cache2go.Cache(common.AccessTokenCache).Delete(tokenCache.AccessToken)
		}
	}
}

// 刷新token
func TokenRefresh(refreshToken string) common.TokenResult {
	res, err := cache2go.Cache(common.RefreshTokenCache).Value(refreshToken)
	if err != nil {
		panic(common.NewError("认证信息已过期，请重新登录"))
	}
	tokenCache := res.Data().(*common.TokenCache)
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
	tokenResult.NickName = userFromDB.NickName
	tokenResult.HasUpdate = userFromDB.HasUpdate

	// 缓存数据
	newTokenCache := common.TokenCache{}
	newTokenCache.AccessToken = accessToken
	newTokenCache.RefreshToken = refreshToken
	newTokenCache.IsAdmin = adminUser.Id == userFromDB.Id
	newTokenCache.UserId = userFromDB.Id
	newTokenCache.Username = userFromDB.Username
	tokenCache.HasUpdate = userFromDB.HasUpdate

	// 缓存token
	cache2go.Cache(common.AccessTokenCache).Add(newTokenCache.AccessToken, AccessTokenExpire, &newTokenCache)
	cache2go.Cache(common.RefreshTokenCache).Add(newTokenCache.RefreshToken, RefreshTokenExpire, &newTokenCache)

	return tokenResult
}

// 获取验证码
func Captcha(signIn common.SignIn) middleware.CaptchaResult {
	// 传入的id是上一次的验证码，此操作为刷新验证码，移除上一次的缓存
	if signIn.CaptchaId != "" {
		cache2go.Cache(common.CaptchaCache).Delete(signIn.CaptchaId)
	}

	// 生成验证码
	captchaResult, err := middleware.CaptchaGet()
	captchaResult.CaptchaId = util.UUID()
	if err != nil {
		panic(common.NewErr("验证码生成失败", err))
	}

	// 缓存验证信息
	captchaCache := middleware.CaptchaCache{X: captchaResult.X, Y: captchaResult.Y}
	cache2go.Cache(common.CaptchaCache).Add(captchaResult.CaptchaId, CaptchaExpire, &captchaCache)

	return captchaResult
}

// 校验验证码
func CaptchaValidate(signIn common.SignIn) bool {
	// 根据验证码id从缓存获取验证码信息
	res, err := cache2go.Cache(common.CaptchaCache).Value(signIn.CaptchaId)
	if err != nil {
		return false
	}
	captchaCache := res.Data().(*middleware.CaptchaCache)

	// 校验
	ok := slide.Validate(signIn.CaptchaX, signIn.CaptchaY, captchaCache.X, captchaCache.Y, 6)

	// 移除验证码缓存
	if !ok {
		cache2go.Cache(common.CaptchaCache).Delete(signIn.CaptchaId)
	}

	return ok
}

// 校验登录次数，如已超出则抛出异常
func checkSignInTimes(name string) {
	cache := cache2go.Cache(common.SignInTimesCache)
	signInTimes, err := cache.Value(name)
	times := 1
	expireSecond := int64(0)
	if err != nil {
		cache.Add(name, time.Minute*5, true)
	} else {
		expireSecond = 300 - (time.Now().Unix() - signInTimes.CreatedOn().Unix())
		// 有时过期后缓存不会立即删除，手动重置次数
		if expireSecond <= 0 {
			cache.Add(name, time.Minute*5, true)
		} else {
			times = int(signInTimes.AccessCount()) + 1
		}
	}
	if times > 5 {
		panic(common.NewError(fmt.Sprintf("登录次数已达上限，请于%v分钟后再试", (expireSecond/60 + 1))))
	}
}
