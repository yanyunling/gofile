package controller

import (
	"gofile/middleware"
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/service"

	"github.com/kataras/iris/v12"
)

// 分页查询用户记录
func UserPage(ctx iris.Context) {
	pageCondition := common.PageCondition[string]{}
	resolveParam(ctx, &pageCondition)
	ctx.JSON(common.NewSuccessData("查询成功", service.UserPage(pageCondition)))
}

// 保存用户记录
func UserSave(ctx iris.Context) {
	user := entity.User{}
	resolveParam(ctx, &user)
	service.UserSave(user)
	ctx.JSON(common.NewSuccess("保存成功"))
}

// 删除用户记录
func UserDelete(ctx iris.Context) {
	user := entity.User{}
	resolveParam(ctx, &user)
	service.UserDelete(user.Id)
	ctx.JSON(common.NewSuccess("删除成功"))
}

// 重置用户密码
func UserResetPassword(ctx iris.Context) {
	user := entity.User{}
	resolveParam(ctx, &user)
	service.UserResetPassword(user)
	ctx.JSON(common.NewSuccess("保存成功"))
}

// 更新用户密码
func UserUpdatePassword(ctx iris.Context) {
	userCondition := entity.UserCondition{}
	resolveParam(ctx, &userCondition)
	userId := middleware.CurrentUserCache(ctx).UserId
	service.UserUpdatePassword(userCondition, userId)
	ctx.JSON(common.NewSuccess("保存成功"))
}
