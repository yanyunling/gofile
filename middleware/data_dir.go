package middleware

import (
	"os"
	"path/filepath"

	"github.com/kataras/golog"
)

// 初始化数据目录
func InitDataDir(dataPath, privateDirName, protectedDirName, publicDirName string) error {
	path := filepath.Join(dataPath, privateDirName)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		golog.Error("创建数据目录失败：", err)
		return err
	}

	path = filepath.Join(dataPath, protectedDirName)
	err = os.MkdirAll(path, 0755)
	if err != nil {
		golog.Error("创建数据目录失败：", err)
		return err
	}

	path = filepath.Join(dataPath, publicDirName)
	err = os.MkdirAll(path, 0755)
	if err != nil {
		golog.Error("创建数据目录失败：", err)
		return err
	}

	return nil
}
