package main

import (
	"embed"
	"flag"
	"gofile/controller"
	"gofile/middleware"
	"gofile/model/common"
	"gofile/util"
	"io/fs"
	"net/http"
	"time"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
)

//go:embed web/dist
var web embed.FS

func init() {
	// 解析命令行参数
	flag.StringVar(&common.Port, "p", "9300", "监听端口")
	flag.StringVar(&common.LogPath, "log", "./logs", "日志目录，存放近30天的日志")
	flag.StringVar(&common.DataPath, "data", "./data", "数据目录，存放数据库及文件资源")
	flag.StringVar(&common.AdminPassword, "pass", "Admin123", "超级管理员密码")
	flag.IntVar(&common.DownloadLimitKB, "limit", 10240, "下载限速(KB)，0则不限速")
	flag.Parse()

	// 固定配置
	common.BasicTokenKey = "gofile"
	common.PrivateDirName = "private"
	common.ProtectedDirName = "protected"
	common.PublicDirName = "public"

	// 超级管理员密码加密
	password, err := util.EncryptBcrypt(util.EncryptSHA256([]byte(common.AdminPassword)))
	if err != nil {
		golog.Error("超级管理员密码加密失败")
		return
	}
	common.AdminPassword = password
}

func main() {
	// 创建iris服务
	app := iris.New()

	// 初始化日志
	middleware.InitLog()

	// 全局异常恢复
	app.Use(middleware.GlobalRecover)

	// 初始化缓存
	middleware.InitCache()

	// 初始化验证码组件
	err := middleware.InitCaptcha()
	if err != nil {
		return
	}

	// 初始化数据目录
	err = middleware.InitDataDir(common.DataPath, common.PrivateDirName, common.ProtectedDirName, common.PublicDirName)
	if err != nil {
		return
	}

	// 初始化数据库连接
	err = middleware.InitDB()
	if err != nil {
		return
	}
	defer middleware.Db.Close()

	// 初始化API路由
	controller.InitRouter(app)

	// gzip压缩
	app.Use(iris.Compression)

	// 网页资源路由
	app.Use(iris.StaticCache(time.Hour * 720))
	webFs, err := fs.Sub(web, "web/dist")
	if err != nil {
		golog.Error("初始化网页资源失败：", err)
		return
	}
	app.HandleDir("/", http.FS(webFs))

	// 启动服务
	app.Logger().Error(app.Run(iris.Addr(":" + common.Port)))
}
