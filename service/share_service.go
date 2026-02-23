package service

import (
	"gofile/middleware"
	"gofile/model/common"
	"gofile/model/entity"
	"gofile/repository"
	"gofile/util"
	"time"
)

// 添加分享记录
func ShareAdd(share entity.Share) string {
	if share.ShareHours <= 0 || share.ShareHours > 720 {
		panic(common.NewError("分享时间最长720小时"))
	}

	tx := middleware.Db.MustBegin()
	defer tx.Rollback()

	startTime := time.Now()
	endTime := startTime.Add(time.Duration(share.ShareHours) * time.Hour)
	share.Id = util.UUIDNoHyphen()
	share.StartTime = startTime.UnixMilli()
	share.EndTime = endTime.UnixMilli()
	err := repository.ShareAdd(tx, share)
	if err != nil {
		panic(common.NewErr("分享失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("分享失败", err))
	}

	return share.Id
}

// 查询未过期的分享记录
func ShareGet(id string) entity.Share {
	// 查询分享记录
	share, err := repository.ShareGet(middleware.Db, id)
	if err != nil {
		panic(common.NewErr("分享链接已失效", err))
	}

	// 判断是否过期，过期则删除
	currentTime := time.Now().UnixMilli()
	if share.EndTime <= currentTime {
		repository.ShareDelete(middleware.Db, id)
		panic(common.NewError("分享链接已失效"))
	}

	return share
}

// 分页查询分享记录
func SharePage(pageCondition common.PageCondition[entity.Share]) common.PageResult[entity.Share] {
	records, total, err := repository.SharePage(middleware.Db, pageCondition.Page, pageCondition.Condition)
	if err != nil {
		panic(common.NewErr("查询失败", err))
	}

	pageResult := common.PageResult[entity.Share]{Records: records, Total: total}
	return pageResult
}

// 删除分享记录
func ShareDelete(id, username string) {
	tx := middleware.Db.MustBegin()
	defer tx.Rollback()

	// 删除分享记录
	err := repository.ShareDeleteById(tx, id, username)
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}
}
