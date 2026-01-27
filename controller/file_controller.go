package controller

import (
	"gofile/middleware"
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/service"
	"os"

	"github.com/kataras/iris/v12"
)

// 查询公开文件列表
func FileListPublic(ctx iris.Context) {
	fileInfo := entity.FileInfo{}
	resolveParam(ctx, &fileInfo)
	result := service.FileList(common.PublicDirName, fileInfo.Path)
	ctx.JSON(common.NewSuccessData("查询成功", result))
}

// 查询保护文件列表
func FileListProtected(ctx iris.Context) {
	fileInfo := entity.FileInfo{}
	resolveParam(ctx, &fileInfo)
	result := service.FileList(common.ProtectedDirName, fileInfo.Path)
	ctx.JSON(common.NewSuccessData("查询成功", result))
}

// 查询私有文件列表
func FileListPrivate(ctx iris.Context) {
	fileInfo := entity.FileInfo{}
	resolveParam(ctx, &fileInfo)

	// 管理员可以查看所有人的文件列表
	tokenCache := middleware.CurrentUserCache(ctx)
	var path string
	if tokenCache.IsAdmin {
		path = fileInfo.Path
	} else {
		path = tokenCache.Username + "/" + fileInfo.Path
		// 确保创建用户名目录
		os.MkdirAll(common.DataPath+common.PrivateDirName+"/"+tokenCache.Username, 0777)
	}

	result := service.FileList(common.PrivateDirName, path)
	ctx.JSON(common.NewSuccessData("查询成功", result))
}
