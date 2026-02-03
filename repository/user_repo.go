package repository

import (
	"errors"
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/util"

	"github.com/jmoiron/sqlx"
)

// 根据id查询用户
func UserGetById[T *sqlx.Tx | *sqlx.DB](db T, id string) (entity.User, error) {
	// 超级管理员用户特殊处理
	adminUser := entity.GetAdminUser()
	if adminUser.Id == id {
		return adminUser, nil
	}

	sql := `select * from t_user where id=$1`
	result := entity.User{}
	var err error
	switch v := any(db).(type) {
	case *sqlx.Tx:
		err = v.Get(&result, sql, id)
	case *sqlx.DB:
		err = v.Get(&result, sql, id)
	}

	return result, err
}

// 根据用户名查询用户
func UserGetByUsername(db *sqlx.DB, username string) (entity.User, error) {
	// 超级管理员用户特殊处理
	adminUser := entity.GetAdminUser()
	if adminUser.Username == username {
		return adminUser, nil
	}

	sql := `select * from t_user where username=$1`
	result := entity.User{}
	err := db.Get(&result, sql, username)

	return result, err
}

// 分页查询用户记录列表
func UserPage(db *sqlx.DB, page common.Page, condition string) ([]entity.User, int, error) {
	sqlCompletion := util.SqlCompletion{}
	sqlCompletion.InitSql(`select * from t_user`)

	// 状态筛选特殊处理，根据启用和停用区分
	switch condition {
	case "启用":
		sqlCompletion.EqOr("enable", 1)
	case "停用":
		sqlCompletion.EqOr("enable", 0)
	}

	// 更新权限筛选特殊处理，根据有和无区分
	switch condition {
	case "有":
		sqlCompletion.EqOr("public_auth", 1)
		sqlCompletion.EqOr("protected_auth", 1)
		sqlCompletion.EqOr("private_auth", 1)
	case "无":
		sqlCompletion.EqOr("public_auth", 0)
		sqlCompletion.EqOr("protected_auth", 0)
		sqlCompletion.EqOr("private_auth", 0)
	}

	sqlCompletion.LikeOr([]string{"username"}, condition)
	sqlCompletion.Order("username", true)
	sqlCompletion.Limit(page.Current, page.Size)

	// 查询分页数据
	result := []entity.User{}
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

// 添加用户记录
func UserAdd(tx *sqlx.Tx, user entity.User) error {
	// 超级管理员用户特殊处理
	adminUser := entity.GetAdminUser()
	if adminUser.Username == user.Username {
		return errors.New("用户已存在")
	}

	sql := `insert into t_user 
	(id,
	username,
	password,
	enable,
	public_auth,
	protected_auth,
	private_auth)
	values 
	(:id,
	:username,
	:password,
	:enable,
	:public_auth,
	:protected_auth,
	:private_auth)`
	_, err := tx.NamedExec(sql, user)
	return err
}

// 修改用户记录
func UserUpdate(tx *sqlx.Tx, user entity.User) error {
	sql := `update t_user set 
	enable=:enable,
	public_auth=:public_auth,
	protected_auth=:protected_auth,
	private_auth=:private_auth
	where id=:id`
	_, err := tx.NamedExec(sql, user)
	return err
}

// 修改用户密码
func UserUpdatePassword(tx *sqlx.Tx, user entity.User) error {
	sql := `update t_user set 
	password=:password
	where id=:id`
	_, err := tx.NamedExec(sql, user)
	return err
}

// 根据id删除用户记录
func UserDeleteById(tx *sqlx.Tx, id string) error {
	sql := `delete from t_user where id=$1`
	_, err := tx.Exec(sql, id)
	return err
}
