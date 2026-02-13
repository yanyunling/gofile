package service

import (
	"gofile/middleware"
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/repository"
	"gofile/util"
	"time"

	"github.com/kataras/golog"
)

// 添加日志记录
func LogAdd(level entity.LogLevel, title, content, username string) {
	// 写入日志到文件
	logText := "[" + title + "]"
	if username != "" {
		logText += "[" + username + "]"
	}
	logText += " " + content
	switch level {
	case entity.Debug:
		golog.Debug(logText)
	case entity.Info:
		golog.Info(logText)
	case entity.Warn:
		golog.Warn(logText)
	case entity.Error:
		golog.Error(logText)
	}

	// 写入日志到数据库
	log := entity.Log{Id: util.UUID(), Title: title, Content: content, Level: level, Username: username, CreateTime: time.Now().UnixMilli()}
	err := repository.LogAdd(middleware.Db, log)
	if err != nil {
		golog.Warn("日志保存失败", err)
	}
}

// 分页查询日志记录
func LogPage(pageCondition common.PageCondition[entity.LogCondition]) common.PageResult[entity.Log] {
	records, total, err := repository.LogPage(middleware.Db, pageCondition.Page, pageCondition.Condition)
	if err != nil {
		panic(common.NewErr("查询失败", err))
	}

	pageResult := common.PageResult[entity.Log]{Records: records, Total: total}
	return pageResult
}
