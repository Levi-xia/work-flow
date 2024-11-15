package serctx

import (
	"github.com/jmoiron/sqlx"
)

var SerCtx *ServerContext

type ServerContext struct {
	Db *sqlx.DB
}

func InitServerContext() {
	// 初始化数据库
	db, _ := initDB()

	SerCtx = &ServerContext{
		Db: db,
	}
}
