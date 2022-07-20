package main

import (
	"blog/api"
	db "blog/db/sqlc"
	"blog/util"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can not loadconfig", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("connot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot create new server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

}
