package service

import (
	"gofile/middleware"
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/repository"
	"gofile/util"
	"regexp"
	"strings"
)

// 分页查询用户记录
func UserPage(pageCondition common.PageCondition[string]) common.PageResult[entity.User] {
	records, total, err := repository.UserPage(middleware.Db, pageCondition.Page, pageCondition.Condition)
	if err != nil {
		panic(common.NewErr("查询失败", err))
	}

	// 密码置空
	for i := range records {
		records[i].Password = ""
	}

	pageResult := common.PageResult[entity.User]{Records: records, Total: total}
	return pageResult
}

// 保存用户记录
func UserSave(user entity.User) {
	// 超级管理员用户特殊处理
	adminUser := entity.GetAdminUser()
	if adminUser.Username == user.Username {
		panic(common.NewError("用户名已存在"))
	}

	tx := middleware.Db.MustBegin()
	defer tx.Rollback()

	// 数据校验
	regex := regexp.MustCompile(`^[\w\-\p{Han}]{1,20}$`)
	if !regex.MatchString(user.Username) {
		panic(common.NewError("用户名：只允许20个字符以内的数字、英文、汉字、中划线、下划线"))
	}

	// 更新用户
	var err error
	if user.Id == "" {
		// 新增
		user.Id = util.UUID()
		// 生成新密码
		newPassword, err := util.EncryptBcrypt(user.Password)
		if err != nil {
			panic(common.NewErr("密码加密失败", err))
		}
		user.Password = newPassword
		err = repository.UserAdd(tx, user)
	} else {
		// 修改
		err = repository.UserUpdate(tx, user)
	}

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			panic(common.NewErr("用户名已存在", err))
		}
		panic(common.NewErr("保存失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("保存失败", err))
	}
}

// 删除用户记录
func UserDelete(id string) {
	tx := middleware.Db.MustBegin()
	defer tx.Rollback()

	// 删除用户
	err := repository.UserDeleteById(tx, id)
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}
}

// 重置用户密码
func UserResetPassword(user entity.User) {
	tx := middleware.Db.MustBegin()
	defer tx.Rollback()

	// 生成新密码
	newPassword, err := util.EncryptBcrypt(user.Password)
	if err != nil {
		panic(common.NewErr("密码加密失败", err))
	}
	user.Password = newPassword
	err = repository.UserUpdatePassword(tx, user)

	if err != nil {
		panic(common.NewErr("保存失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("保存失败", err))
	}
}

// 更新用户密码
func UserUpdatePassword(userCondition entity.UserCondition, userId string) {
	tx := middleware.Db.MustBegin()
	defer tx.Rollback()

	// 查询登录用户信息
	loginUser, err := repository.UserGetById(tx, userId)
	if err != nil {
		panic(common.NewErr("保存失败", err))
	}

	// 判断原密码正确性
	if !util.CheckBcrypt(loginUser.Password, userCondition.Password) {
		panic(common.NewErr("原密码不正确", err))
	}

	// 生成新密码
	newPassword, err := util.EncryptBcrypt(userCondition.NewPassword)
	if err != nil {
		panic(common.NewErr("密码加密失败", err))
	}

	user := entity.User{}
	user.Id = userId
	user.Password = newPassword
	err = repository.UserUpdatePassword(tx, user)

	if err != nil {
		panic(common.NewErr("保存失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("保存失败", err))
	}
}
