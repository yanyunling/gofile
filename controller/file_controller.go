package controller

import (
	"gofile/middleware"
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/service"
	"os"
	"path/filepath"

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
		path = filepath.Join(tokenCache.Username, fileInfo.Path)
		// 确保创建用户名目录
		os.MkdirAll(filepath.Join(common.DataPath, common.PrivateDirName, tokenCache.Username), 0755)
	}

	result := service.FileList(common.PrivateDirName, path)
	ctx.JSON(common.NewSuccessData("查询成功", result))
}

// 下载公开文件
func FileDownloadPublic(ctx iris.Context) {
	fileInfo := entity.FileInfo{}
	resolveParam(ctx, &fileInfo)
	service.FileDownload(ctx, common.PublicDirName, fileInfo.Path, fileInfo.Name)
}

// 下载保护文件
func FileDownloadProtected(ctx iris.Context) {
	fileInfo := entity.FileInfo{}
	resolveParam(ctx, &fileInfo)
	service.FileDownload(ctx, common.ProtectedDirName, fileInfo.Path, fileInfo.Name)
}

// 下载私有文件
func FileDownloadPrivate(ctx iris.Context) {
	fileInfo := entity.FileInfo{}
	resolveParam(ctx, &fileInfo)

	// 管理员可以下载所有人的文件
	tokenCache := middleware.CurrentUserCache(ctx)
	var path string
	if tokenCache.IsAdmin {
		path = fileInfo.Path
	} else {
		path = filepath.Join(tokenCache.Username, fileInfo.Path)
	}

	service.FileDownload(ctx, common.PrivateDirName, path, fileInfo.Name)
}

// 创建公开目录
func FileFolderPublic(ctx iris.Context) {
	fileInfo := entity.FileInfo{}
	resolveParam(ctx, &fileInfo)
	service.FileFolder(common.PublicDirName, fileInfo.Path, fileInfo.Name)
	ctx.JSON(common.NewSuccess("创建成功"))
}

// 创建保护目录
func FileFolderProtected(ctx iris.Context) {
	fileInfo := entity.FileInfo{}
	resolveParam(ctx, &fileInfo)
	service.FileFolder(common.ProtectedDirName, fileInfo.Path, fileInfo.Name)
	ctx.JSON(common.NewSuccess("创建成功"))
}

// 创建私有目录
func FileFolderPrivate(ctx iris.Context) {
	fileInfo := entity.FileInfo{}
	resolveParam(ctx, &fileInfo)

	// 管理员可以创建所有目录
	tokenCache := middleware.CurrentUserCache(ctx)
	var path string
	if tokenCache.IsAdmin {
		path = common.PrivateDirName
	} else {
		path = filepath.Join(common.PrivateDirName, tokenCache.Username)
	}

	service.FileFolder(path, fileInfo.Path, fileInfo.Name)
	ctx.JSON(common.NewSuccess("创建成功"))
}

// 上传公开文件
func FileUploadPublic(ctx iris.Context) {
	service.FileUpload(ctx, common.PublicDirName)
	ctx.JSON(common.NewSuccess("上传成功"))
}

// 上传保护文件
func FileUploadProtected(ctx iris.Context) {
	service.FileUpload(ctx, common.ProtectedDirName)
	ctx.JSON(common.NewSuccess("上传成功"))
}

// 上传私有文件
func FileUploadPrivate(ctx iris.Context) {
	// 管理员可以上传到所有人的目录内
	tokenCache := middleware.CurrentUserCache(ctx)
	var path string
	if tokenCache.IsAdmin {
		path = common.PrivateDirName
	} else {
		path = filepath.Join(common.PrivateDirName, tokenCache.Username)
	}

	service.FileUpload(ctx, path)
	ctx.JSON(common.NewSuccess("上传成功"))
}

// 删除公开文件
func FileDeletePublic(ctx iris.Context) {
	fileInfo := entity.FileInfo{}
	resolveParam(ctx, &fileInfo)
	service.FileDelete(common.PublicDirName, fileInfo.Path, fileInfo.Name)
	ctx.JSON(common.NewSuccess("删除成功"))
}

// 删除保护文件
func FileDeleteProtected(ctx iris.Context) {
	fileInfo := entity.FileInfo{}
	resolveParam(ctx, &fileInfo)
	service.FileDelete(common.ProtectedDirName, fileInfo.Path, fileInfo.Name)
	ctx.JSON(common.NewSuccess("删除成功"))
}

// 删除私有文件
func FileDeletePrivate(ctx iris.Context) {
	fileInfo := entity.FileInfo{}
	resolveParam(ctx, &fileInfo)

	// 管理员可以删除所有人的文件
	tokenCache := middleware.CurrentUserCache(ctx)
	var path string
	if tokenCache.IsAdmin {
		path = common.PrivateDirName
	} else {
		path = filepath.Join(common.PrivateDirName, tokenCache.Username)
	}

	service.FileDelete(path, fileInfo.Path, fileInfo.Name)
	ctx.JSON(common.NewSuccess("删除成功"))
}
