package repository

import (
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/util"

	"github.com/jmoiron/sqlx"
)

// 添加日志记录
func LogAdd(db *sqlx.DB, log entity.Log) error {
	sql := `insert into t_log 
	(id,
	title,
	content,
	level,
	username,
	create_time)
	values 
	(:id,
	:title,
	:content,
	:level,
	:username,
	:create_time)`
	_, err := db.NamedExec(sql, log)
	return err
}

// 分页查询日志记录
func LogPage(db *sqlx.DB, page common.Page, condition entity.LogCondition) ([]entity.Log, int, error) {
	sqlCompletion := util.SqlCompletion{}
	sqlCompletion.InitSql(`select * from t_log`)
	if condition.Title != "" {
		sqlCompletion.Eq("title", condition.Title)
	}
	if condition.Level != "" {
		sqlCompletion.Eq("level", condition.Level)
	}
	if condition.Username != "" {
		sqlCompletion.Eq("username", condition.Username)
	}
	if condition.StartTime > 0 {
		sqlCompletion.Ge("create_time", condition.StartTime)
	}
	if condition.EndTime > 0 {
		sqlCompletion.Le("create_time", condition.EndTime)
	}
	sqlCompletion.Like("content", condition.Content)
	sqlCompletion.Order("create_time", false)
	sqlCompletion.Limit(page.Current, page.Size)

	// 查询分页数据
	result := []entity.Log{}
	err := db.Select(&result, sqlCompletion.GetSql(), sqlCompletion.GetParams()...)
	if err != nil {
		return result, 0, err
	}

	// 查询总记录数
	countResult := common.CountResult{}
	err = db.Get(&countResult, sqlCompletion.GetCountSql(), sqlCompletion.GetCountParams()...)
	if err != nil {
		return result, 0, err
	}

	return result, countResult.Count, err
}
