package controller

import (
	"gofile/middleware"
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/service"

	"github.com/kataras/iris/v12"
)

// 分页查询分享记录
func SharePage(ctx iris.Context) {
	pageCondition := common.PageCondition[entity.Share]{}
	resolveParam(ctx, &pageCondition)

	// 非管理员只能查看自己的分享记录
	tokenCache := middleware.CurrentUserCache(ctx)
	if !tokenCache.IsAdmin {
		pageCondition.Condition.Username = tokenCache.Username
	}

	ctx.JSON(common.NewSuccessData("查询成功", service.SharePage(pageCondition)))
}

// 删除分享记录
func ShareDelete(ctx iris.Context) {
	share := entity.Share{}
	resolveParam(ctx, &share)

	// 非管理员只能删除自己的分享记录
	tokenCache := middleware.CurrentUserCache(ctx)
	if !tokenCache.IsAdmin {
		share.Username = tokenCache.Username
	} else {
		share.Username = ""
	}

	service.ShareDelete(share.Id, share.Username)
	ctx.JSON(common.NewSuccess("删除成功"))
}
