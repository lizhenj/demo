package database

import (
	"database/sql"
	"demo/log"
	"fmt"
)

const (
	actorSql = `CREATE TABLE IF NOT EXISTS actor (
    id int(11) default 0,
    basedata varchar(100)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8`
)

var (
	gamedb *sql.DB

	stmtInsertActor *sql.Stmt
)

const (
	DBConnectParam = "?loc=Local&parseTime=true&readTimeout=30s&writeTimeout=30s"
)

func InitDb(cfg string) {
	var err error
	if gamedb, err = sql.Open("mysql", fmt.Sprintf("%s%s", cfg, DBConnectParam)); err != nil {
		panic(err)
	} else if err = gamedb.Ping(); err != nil {
		panic(err)
	}

	gamedb.SetConnMaxLifetime(0)

	//设置表
	for _, tbsql := range []string{actorSql} {
		if _, err = gamedb.Exec(tbsql); err != nil {
			panic(err)
		}
	}

	//预处理
	stmtInsertActor, err = gamedb.Prepare("Insert into actor(id,basedata) values (?,?)")
	if err != nil {
		panic(err)
	}
}

func InsertActor(args ...interface{}) (res bool) {
	if len(args) <= 0 {
		return res
	}

	if _, err := stmtInsertActor.Query(args[0].(int), args[1].(string)); err != nil {
		log.Errorf("InsertActor failed,err: %v", err)
		return res
	}
	return true
}
