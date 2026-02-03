package controller

import (
	"gofile/middleware"
	"gofile/model/common"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
)

// 初始化iris路由
func InitRouter(app *iris.Application) {
	// 允许跨域
	app.UseRouter(cors.AllowAll())

	app.PartyFunc("/api", func(api iris.Party) {
		// token接口
		api.PartyFunc("/token", func(token iris.Party) {
			// token接口授权
			token.Use(middleware.TokenAuth)
			token.Post("/sign-in", SignIn)
			token.Post("/sign-out", SignOut)
			token.Post("/refresh", TokenRefresh)
			token.Post("/captcha", Captcha)
			token.Post("/captcha/validate", CaptchaValidate)
		})

		// 数据接口
		api.PartyFunc("/data", func(data iris.Party) {
			// 数据接口授权
			data.Use(middleware.DataAuth)

			// 用户管理
			data.PartyFunc("/user", func(user iris.Party) {
				user.Post("/update-password", UserUpdatePassword)

				// 管理员接口授权
				user.Use(middleware.AdminAuth)
				user.Post("/page", UserPage)
				user.Post("/save", UserSave)
				user.Post("/delete", UserDelete)
				user.Post("/reset-password", UserResetPassword)
			})

			// 文件管理
			data.PartyFunc("/file", func(file iris.Party) {
				file.Post("/public/folder", FileFolderPublic)
				file.Post("/public/upload", FileUploadPublic)
				file.Post("/public/delete", FileDeletePublic)
				file.Post("/protected/list", FileListProtected)
				file.Post("/protected/download", FileDownloadProtected)
				file.Post("/protected/folder", FileFolderProtected)
				file.Post("/protected/upload", FileUploadProtected)
				file.Post("/protected/delete", FileDeleteProtected)
				file.Post("/private/list", FileListPrivate)
				file.Post("/private/download", FileDownloadPrivate)
				file.Post("/private/folder", FileFolderPrivate)
				file.Post("/private/upload", FileUploadPrivate)
				file.Post("/private/delete", FileDeletePrivate)
			})
		})

		// 开放接口
		api.PartyFunc("/open", func(open iris.Party) {
			// 文件管理
			open.PartyFunc("/file", func(file iris.Party) {
				file.Post("/public/list", FileListPublic)
				file.Post("/public/download", FileDownloadPublic)
			})
		})
	})
}

// 解析参数
func resolveParam(ctx iris.Context, con any) {
	err := ctx.ReadJSON(&con)
	if err != nil {
		panic(common.NewErr("参数解析失败", err))
	}
}
