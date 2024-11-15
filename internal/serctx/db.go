package serctx

import "github.com/jmoiron/sqlx"

func initDB() (*sqlx.DB, error) {
	// 初始化数据库连接
	dsn := "root:Mysqlxwy9264@tcp(127.0.0.1:3306)/workflow?charset=utf8mb4&parseTime=True"
	db := sqlx.MustConnect("mysql", dsn)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return db, nil
}
