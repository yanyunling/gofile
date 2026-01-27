package service

import (
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/util"
	"os"
	"sort"

	"github.com/kataras/golog"
)

// 查询文件列表
func FileList(parentDir, path string) []entity.FileInfo {
	// 读取当前目录
	dirPath := util.ReplaceMultipleSlashes(common.DataPath + parentDir + "/" + path)
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
			Path:       util.ReplaceMultipleSlashes(util.PathCompletion(path)),
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
