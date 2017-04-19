package util

import (
	"database/sql"

	"sync/cmd/app/options"

	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Db *sql.DB
)

func InitDB() {
	db, err := sql.Open("mysql", options.GlobalConfig.Dsn)
	if err != nil {
		logs.Critical("init the database's connection err:%v", err)
	}
	if err = db.Ping(); err != nil {
		logs.Critical("can't connect to the database, the err is :%v", err)
	}
	db.SetMaxOpenConns(500)
	db.SetMaxIdleConns(5)
	Db = db
}
