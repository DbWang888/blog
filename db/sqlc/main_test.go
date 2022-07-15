package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var testQueries *Queries
var testDB *sql.DB

const (
	dbDriver = "mysql"
	// dbSource = "mysql://root:4524@tcp(localhost:3306)/blog"
	dbSource = "root:4524@tcp(127.0.0.1:3306)/blog?charset=utf8&parseTime=True&loc=Local"
)

func TestMain(m *testing.M) {

	var err error

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("connot connect to db:", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
