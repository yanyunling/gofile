package repository

import (
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/util"
	"time"

	"github.com/jmoiron/sqlx"
)

// 根据id查询分享记录
func ShareGet(db *sqlx.DB, id string) (entity.Share, error) {
	sql := `select * from t_share where id=$1`
	result := entity.Share{}
	err := db.Get(&result, sql, id)
	return result, err
}

// 添加分享记录
func ShareAdd(tx *sqlx.Tx, share entity.Share) error {
	sql := `insert into t_share 
	(id,
	parent_dir,
	path,
	name,
	username,
	share_hours,
	start_time,
	end_time)
	values 
	(:id,
	:parent_dir,
	:path,
	:name,
	:username,
	:share_hours,
	:start_time,
	:end_time)`
	_, err := tx.NamedExec(sql, share)
	return err
}

// 根据id删除分享记录
func ShareDelete(db *sqlx.DB, id, username string) error {
	var err error
	if username == "" {
		sql := `delete from t_share where id=$1`
		_, err = db.Exec(sql, id)
	} else {
		sql := `delete from t_share where id=$1 and username=$2`
		_, err = db.Exec(sql, id, username)
	}
	return err
}

// 分页查询分享记录
func SharePage(db *sqlx.DB, page common.Page, condition entity.Share) ([]entity.Share, int, error) {
	currentTime := time.Now().UnixMilli()
	sqlCompletion := util.SqlCompletion{}
	sqlCompletion.InitSql(`select * from t_share`)
	sqlCompletion.Ge("end_time", currentTime)
	sqlCompletion.Like("parent_dir", condition.ParentDir)
	sqlCompletion.Like("path", condition.Path)
	sqlCompletion.Like("name", condition.Name)
	if condition.Username != "" {
		sqlCompletion.Eq("username", condition.Username)
	}
	if condition.ShareHours > 0 {
		sqlCompletion.Eq("share_hours", condition.ShareHours)
	}
	if condition.StartTime > 0 {
		sqlCompletion.Ge("start_time", condition.StartTime)
	}
	if condition.EndTime > 0 {
		sqlCompletion.Le("start_time", condition.EndTime)
	}
	sqlCompletion.Order("start_time", false)
	sqlCompletion.Limit(page.Current, page.Size)

	// 查询分页数据
	result := []entity.Share{}
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
