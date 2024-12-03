package serctx

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"workflow/config"
)

func initDB() (*sqlx.DB, error) {
	// 初始化数据库连接
	dbConf := config.Conf.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True", dbConf.Username, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Database, dbConf.Charset)
	db := sqlx.MustConnect("mysql", dsn)
	db.SetMaxOpenConns(dbConf.MaxOpenConns)
	db.SetMaxIdleConns(dbConf.MaxIdleConns)
	return db, nil
}
