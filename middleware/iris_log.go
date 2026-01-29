package middleware

import (
	"gofile/model/common"
	"io"
	"os"
	"path/filepath"

	"github.com/kataras/golog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 初始化日志
func InitLog() {
	golog.SetLevel("info")
	golog.SetTimeFormat("2006-01-02 15:04:05")

	// 文件输出
	logger := &lumberjack.Logger{
		Filename:   filepath.Join(common.LogPath, "app.log"),
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	}
	golog.SetOutput(io.MultiWriter(logger, os.Stdout))
}
