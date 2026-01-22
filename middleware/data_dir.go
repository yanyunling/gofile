package middleware

import (
	"os"

	"github.com/kataras/golog"
)

// 初始化数据目录
func InitDataDir(dataPath, privateDirName, protectedDirName, publicDirName string) error {
	path := dataPath + privateDirName
	err := os.MkdirAll(path, 0777)
	if err != nil {
		golog.Error("创建数据目录失败：", err)
		return err
	}

	path = dataPath + protectedDirName
	err = os.MkdirAll(path, 0777)
	if err != nil {
		golog.Error("创建数据目录失败：", err)
		return err
	}

	path = dataPath + publicDirName
	err = os.MkdirAll(path, 0777)
	if err != nil {
		golog.Error("创建数据目录失败：", err)
		return err
	}

	return nil
}
