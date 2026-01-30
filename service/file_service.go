package service

import (
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/util"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
)

// 查询文件列表
func FileList(parentDir, path string) []entity.FileInfo {
	// 读取当前目录
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
			Path:       filepath.Join(path),
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
func FileDownload(ctx iris.Context, parentDir, path, fileName string) {
	// 判断是目录还是文件
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
		ctx.SendFile(filePath, fileName)
	}
}

// 创建目录
func FileFolder(parentDir, path, fileName string) {
	// 判断目录是否存在
	filePath := filepath.Join(common.DataPath, parentDir, path, fileName)
	_, err := util.PathIsDir(filePath)
	if err == nil {
		panic(common.NewErr("目录已存在", err))
	}

	// 创建目录
	os.MkdirAll(filePath, 0755)
}

// 上传文件
func FileUpload(ctx iris.Context, parentDir string) {
	// 设置最大上传大小1GB
	ctx.SetMaxRequestBodySize(1 << 30)

	// 获取文件
	path := ctx.FormValue("path")
	file, fileHeader, err := ctx.FormFile("file")
	if err != nil {
		panic(common.NewErr("上传失败", err))
	}
	defer file.Close()

	// 创建目录
	os.MkdirAll(filepath.Join(common.DataPath, parentDir, path), 0755)

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
}

// 删除文件
func FileDelete(parentDir, path, fileName string) {
	err := os.RemoveAll(filepath.Join(common.DataPath, parentDir, path, fileName))
	if err != nil {
		panic(common.NewErr("文件删除失败", err))
	}
}
