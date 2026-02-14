package middleware

import (
	"gofile/model/common"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/golog"
	_ "modernc.org/sqlite"
)

// 数据库连接
var Db *sqlx.DB

// 建表语句
var createTableSql = `
CREATE TABLE IF NOT EXISTS t_user
(
	id varchar(36) PRIMARY KEY NOT NULL,
	username text NOT NULL,
	password text NOT NULL,
	enable int NOT NULL,
	public_auth int NOT NULL,
	protected_auth int NOT NULL,
	private_auth int NOT NULL
);

CREATE TABLE IF NOT EXISTS t_log
(
	id varchar(36) PRIMARY KEY NOT NULL,
	title text NOT NULL,
	content text NOT NULL,
	level text NOT NULL,
	username text NOT NULL,
	create_time bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS t_share
(
	id varchar(36) PRIMARY KEY NOT NULL,
	parent_dir text NOT NULL,
	path text NOT NULL,
	name text NOT NULL,
	username text NOT NULL,
	share_hours int NOT NULL,
	start_time bigint NOT NULL,
	end_time bigint NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS "user_username"
ON "t_user" (
  "username" ASC
);
`

type pragma struct {
	PageCount     int64 `json:"pageCount" db:"page_count"`
	FreelistCount int64 `json:"freelistCount" db:"freelist_count"`
}

// 初始化数据库连接
func InitDB() error {
	// 开启数据库文件
	var err error
	Db, err = sqlx.Connect("sqlite", filepath.Join(common.DataPath, "gofile.db"))
	if err != nil {
		golog.Error("开启sqlite数据库文件失败：", err)
		return err
	}

	// 开启WAL模式
	_, err = Db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		golog.Error("开启sqlite WAL模式失败：", err)
		return err
	}

	// 创建表结构
	Db.MustExec(createTableSql)

	// 数据库定时任务
	go dbTask()

	return nil
}

// 数据库定时任务
func dbTask() {
	// 每10分钟执行一次
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		// 清理分享表
		shareClear()

		// 清理日志表
		logClear()

		// 回收数据库空间
		vaccum()
	}
}

// 执行VACUUM回收空间
func vaccum() {
	// 查询总页数
	p1 := pragma{}
	err := Db.Get(&p1, "PRAGMA page_count;")
	if err != nil {
		golog.Error("[定时任务][VACUUM 回收数据库空间] PRAGMA page_count 执行失败", err)
		return
	}

	// 查询空闲页数量
	p2 := pragma{}
	err = Db.Get(&p2, "PRAGMA freelist_count;")
	if err != nil {
		golog.Error("[定时任务][VACUUM 回收数据库空间] PRAGMA freelist_count 执行失败", err)
		return
	}

	// 计算空闲页占比，超过30%则回收空间
	pageCount := p1.PageCount
	freelistCount := p2.FreelistCount
	if pageCount == 0 {
		return
	}
	ratio := float64(freelistCount) * 100.0 / float64(pageCount)
	if ratio > 30 {
		// 回收空闲空间
		_, err = Db.Exec("VACUUM;")
		if err != nil {
			golog.Error("[定时任务][VACUUM 回收数据库空间] 执行失败", err)
		} else {
			golog.Info("[定时任务][VACUUM 回收数据库空间] 执行成功，总页数：" + strconv.FormatInt(pageCount, 10) + "，空闲页数：" + strconv.FormatInt(freelistCount, 10))
		}
	}

}

// 清理超过100万行的日志记录，每次执行删除1万行
func logClear() {
	sql := `DELETE FROM t_log WHERE id IN (SELECT id FROM t_log ORDER BY id LIMIT 10000) AND EXISTS (SELECT 1 FROM t_log ORDER BY id LIMIT 1 OFFSET 1000000);`
	result, err := Db.Exec(sql)
	if err != nil {
		golog.Error("[定时任务][清理日志表] 执行失败", err)
	} else {
		count, err := result.RowsAffected()
		if err == nil && count > 0 {
			golog.Info("[定时任务][清理日志表] 执行成功")
		}
	}
}

// 清理过期的分享记录
func shareClear() {
	currentTime := time.Now().UnixMilli()
	sql := `DELETE FROM t_share WHERE end_time<$1`
	result, err := Db.Exec(sql, currentTime)
	if err != nil {
		golog.Error("[定时任务][清理分享表] 执行失败", err)
	} else {
		count, err := result.RowsAffected()
		if err == nil && count > 0 {
			golog.Info("[定时任务][清理分享表] 执行成功，清理行数：" + strconv.FormatInt(count, 10))
		}
	}
}
