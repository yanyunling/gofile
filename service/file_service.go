package service

import (
	"archive/zip"
	"fmt"
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/util"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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
		// 目录，压缩为 zip 后下载
		displayName := strings.TrimPrefix(fileName, "/")
		if displayName == "" {
			displayName = "download"
		}
		zipFileName := displayName + ".zip"

		// 在系统临时目录创建 zip
		tempZipName := fmt.Sprintf("dl_%s_%d.zip", strings.ReplaceAll(displayName, string(filepath.Separator), "_"), time.Now().UnixNano())
		tempZipPath := filepath.Join(os.TempDir(), tempZipName)

		zipFile, err := os.Create(tempZipPath)
		if err != nil {
			panic(common.NewErr("创建压缩文件失败", err))
		}

		defer os.Remove(tempZipPath)
		defer zipFile.Close()

		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()

		baseName := strings.TrimPrefix(fileName, "/")
		if baseName == "" {
			baseName = "download"
		}

		// 遍历目录，写入 zip（包含根目录本身）
		err = filepath.Walk(filePath, func(walkPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 计算相对路径
			relPath, err := filepath.Rel(filePath, walkPath)
			if err != nil {
				return err
			}

			// 根目录本身映射为 baseName，子内容拼接在其下
			if relPath == "." {
				relPath = baseName
			} else {
				relPath = filepath.Join(baseName, relPath)
			}
			relPath = filepath.ToSlash(relPath)

			if info.IsDir() {
				// 目录条目以 / 结尾，保留空目录结构
				_, err = zipWriter.Create(relPath + "/")
				return err
			}

			// 创建 zip 内的文件条目
			fileInZip, err := zipWriter.Create(relPath)
			if err != nil {
				return err
			}

			// 流式复制（walk 回调里不能用 defer，需立即关闭）
			srcFile, err := os.Open(walkPath)
			if err != nil {
				return err
			}
			_, err = io.Copy(fileInZip, srcFile)
			srcFile.Close()
			return err
		})

		if err != nil {
			panic(common.NewErr("压缩目录失败", err))
		}

		// 必须在 SendFileWithRate 之前手动关闭，确保 central directory 写入并落盘
		if err := zipWriter.Close(); err != nil {
			panic(common.NewErr("写入压缩数据失败", err))
		}
		if err := zipFile.Close(); err != nil {
			panic(common.NewErr("关闭压缩文件失败", err))
		}

		// 下载 zip
		ctx.SendFileWithRate(tempZipPath, zipFileName, float64(common.DownloadLimitKB)*iris.KB, 2*common.DownloadLimitKB*iris.KB)

		// 记录日志
		logFilePath := filepath.Join(parentDir, path, fileName)
		content := "用时：" + strconv.FormatFloat(time.Since(start).Seconds(), 'f', 2, 64) + "秒，路径：" + logFilePath + "（目录压缩下载）"
		LogAdd(entity.Info, "下载目录", content, username)
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
	innerPath := filepath.Join(path, fileName)
	_, err := util.PathIsDir(filePath)
	if err == nil {
		panic(common.NewErr("目录已存在："+innerPath, err))
	}

	// 创建目录
	err = os.Mkdir(filePath, 0755)
	if err != nil {
		panic(common.NewErr("目录创建失败："+innerPath, err))
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
	ctx.SetMaxRequestBodySize(iris.GB)

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
