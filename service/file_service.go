package service

import (
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/util"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
)

// 查询文件列表
func FileList(parentDir, path string) []entity.FileInfo {
	// 读取当前目录
	path = filepath.Clean("/" + path)
	dirPath := filepath.Join(common.DataPath, parentDir, path)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		panic(common.NewErr("路径访问出错", err))
	}

	// 遍历目录
	result := []entity.FileInfo{}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			golog.Warn("文件信息访问失败", err)
			continue
		}
		result = append(result, entity.FileInfo{
			Name:       entry.Name(),
			Size:       info.Size(),
			UpdateTime: info.ModTime().UnixMilli(),
			IsDir:      entry.IsDir()})
	}

	// 排序，文件夹在前，再按拼音升序
	sort.Slice(result, func(i, j int) bool {
		if result[i].IsDir != result[j].IsDir {
			return result[i].IsDir && !result[j].IsDir
		}
		return util.StringSort(result[i].Name, result[j].Name)
	})

	return result
}

// 下载文件
func FileDownload(ctx iris.Context, parentDir, path, fileName, username string) {
	start := time.Now()

	// 判断是目录还是文件
	path = filepath.Clean("/" + path)
	fileName = filepath.Clean("/" + fileName)
	filePath := filepath.Join(common.DataPath, parentDir, path, fileName)
	isDir, err := util.PathIsDir(filePath)
	if err != nil {
		panic(common.NewErr("文件不存在", err))
	}

	if isDir {
		// 目录
		panic(common.NewError("暂不支持目录下载"))
	} else {
		// 文件
		ctx.SendFileWithRate(filePath, fileName, float64(common.DownloadLimitKB)*iris.KB, 2*common.DownloadLimitKB*iris.KB)

		// 添加日志
		logFilePath := filepath.Join(parentDir, path, fileName)
		content := "用时：" + strconv.FormatFloat(time.Since(start).Seconds(), 'f', 2, 64) + "秒，路径：" + logFilePath
		LogAdd(entity.Info, "下载文件", content, username)
	}
}

// 分享文件
func FileShare(share entity.Share, username string) string {
	// 判断是目录还是文件
	path := filepath.Clean("/" + share.Path)
	fileName := filepath.Clean("/" + share.Name)
	filePath := filepath.Join(common.DataPath, share.ParentDir, path, fileName)
	isDir, err := util.PathIsDir(filePath)
	if err != nil {
		panic(common.NewErr("文件不存在", err))
	}
	if isDir {
		// 目录
		panic(common.NewError("暂不支持分享目录"))
	}

	// 保存分享信息
	share.Username = username
	id := ShareAdd(share)

	// 添加日志
	logFilePath := filepath.Join(share.ParentDir, path, fileName)
	content := "时长：" + strconv.Itoa(share.ShareHours) + "小时，链接：" + id + "，路径：" + logFilePath
	LogAdd(entity.Info, "分享文件", content, username)

	return id
}

// 下载分享文件
func FileShareDownload(ctx iris.Context, id string) {
	// 查询分享记录
	share := ShareGet(id)

	// 下载
	FileDownload(ctx, share.ParentDir, share.Path, share.Name, " ")
}

// 创建目录
func FileFolder(parentDir, path, fileName, username string) {
	// 判断目录是否存在
	path = filepath.Clean("/" + path)
	fileName = filepath.Clean("/" + fileName)
	filePath := filepath.Join(common.DataPath, parentDir, path, fileName)
	_, err := util.PathIsDir(filePath)
	if err == nil {
		panic(common.NewErr("目录已存在", err))
	}

	// 创建目录
	err = os.Mkdir(filePath, 0755)
	if err != nil {
		panic(common.NewErr("目录创建失败", err))
	}

	// 添加日志
	logFilePath := filepath.Join(parentDir, path, fileName)
	content := "路径：" + logFilePath
	LogAdd(entity.Info, "创建目录", content, username)
}

// 上传文件
func FileUpload(ctx iris.Context, parentDir, username string) {
	start := time.Now()

	// 设置最大上传大小1GB
	ctx.SetMaxRequestBodySize(1 << 30)

	// 获取文件
	path := ctx.FormValue("path")
	path = filepath.Clean("/" + path)
	file, fileHeader, err := ctx.FormFile("file")
	if err != nil {
		panic(common.NewErr("上传失败", err))
	}
	defer file.Close()

	// 判断文件是否存在
	filePath := filepath.Join(common.DataPath, parentDir, path, fileHeader.Filename)
	_, err = util.PathIsDir(filePath)
	if err == nil {
		panic(common.NewErr("文件已存在", err))
	}

	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		panic(common.NewErr("文件保存失败", err))
	}
	defer out.Close()

	// 复制文件
	_, err = io.Copy(out, file)
	if err != nil {
		os.RemoveAll(filePath)
		panic(common.NewErr("文件保存失败", err))
	}

	// 添加日志
	logFilePath := filepath.Join(parentDir, path, fileHeader.Filename)
	content := "用时：" + strconv.FormatFloat(time.Since(start).Seconds(), 'f', 2, 64) + "秒，路径：" + logFilePath
	LogAdd(entity.Info, "上传文件", content, username)
}

// 删除文件
func FileDelete(parentDir, path, fileName, username string) {
	path = filepath.Clean("/" + path)
	fileName = filepath.Clean("/" + fileName)
	filePath := filepath.Join(common.DataPath, parentDir, path, fileName)
	err := os.RemoveAll(filePath)
	if err != nil {
		panic(common.NewErr("文件删除失败", err))
	}

	// 添加日志
	logFilePath := filepath.Join(parentDir, path, fileName)
	content := "路径：" + logFilePath
	LogAdd(entity.Info, "删除文件", content, username)
}
