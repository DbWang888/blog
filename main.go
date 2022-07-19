package main

import (
	"blog/api"
	db "blog/db/sqlc"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbDriver = "mysql"
	// dbSource = "mysql://root:4524@tcp(localhost:3306)/blog"
	dbSource = "root:4524@tcp(127.0.0.1:3306)/blog?charset=utf8&parseTime=True&loc=Local"
)

var querier *db.Queries

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("connot connect to db:", err)
	}

	querier = db.New(conn)
	server, err := api.NewServer(querier)
	if err != nil {
		log.Fatal("cannot create new server", err)
	}

	err = server.Start("0.0.0.0:8080")
	if err != nil {
		log.Fatal("cannot start server", err)
	}

}
