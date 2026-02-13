package controller

import (
	"gofile/middleware"
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/service"

	"github.com/kataras/iris/v12"
)

// 分页查询日志记录
func LogPage(ctx iris.Context) {
	pageCondition := common.PageCondition[entity.LogCondition]{}
	resolveParam(ctx, &pageCondition)

	// 非管理员只能查看自己的日志
	tokenCache := middleware.CurrentUserCache(ctx)
	if !tokenCache.IsAdmin {
		pageCondition.Condition.Username = tokenCache.Username
	}

	ctx.JSON(common.NewSuccessData("查询成功", service.LogPage(pageCondition)))
}
