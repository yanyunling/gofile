package middleware

import (
	"gofile/model/common"
	"path/filepath"

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
	has_update int NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS "user_username"
ON "t_user" (
  "username" ASC
);
`

// 初始化数据库连接
func InitDB() error {
	// 开启数据库文件
	var err error
	Db, err = sqlx.Connect("sqlite", filepath.Join(common.DataPath, "gofile.db"))
	if err != nil {
		golog.Error("开启sqlite数据库文件失败：", err)
		return err
	}

	// 执行VACUUM优化数据库碎片
	_, err = Db.Exec("VACUUM;")
	if err != nil {
		golog.Error("VACUUM执行失败：", err)
	}

	// 只有用户管理操作数据库，暂时无需开启WAL模式
	// 开启WAL模式
	// _, err = Db.Exec("PRAGMA journal_mode=WAL;")
	// if err != nil {
	// 	golog.Error("开启sqlite WAL模式失败：", err)
	// 	return err
	// }

	golog.Info("已连接sqlite")

	// 创建表结构
	Db.MustExec(createTableSql)

	return nil
}
