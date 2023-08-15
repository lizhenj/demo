package database

import (
	"database/sql"
	"fmt"
)

var (
	gamedb *sql.DB
)

const (
	DBConnectParam = "?loc=Local&parseTime=true&readTimeout=30s&writeTimeout=30s"
)

func Initdb(cfg string) {
	var err error
	if gamedb, err = sql.Open("mysql", fmt.Sprintf("%s%s", cfg, DBConnectParam)); err != nil {
		panic(err)
	} else if err = gamedb.Ping(); err != nil {
		panic(err)
	}

	gamedb.SetConnMaxLifetime(0)

}
