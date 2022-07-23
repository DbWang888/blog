package main

import (
	"blog/api"
	db "blog/db/sqlc"
	"blog/e"
	"blog/logger"
	"blog/util"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/natefinch/lumberjack.v2"
)

var config util.Config

//新增日志
func init() {

	err := setupConfig()
	if err != nil {
		log.Fatalf("配置文件加载失败 err: %v", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("日志模块初始化失败 err : %v", err)
	}
}

func setupLogger() error {
	e.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  config.LogSavePath + "/" + config.LogFileName + config.LogFileExt,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}

func setupConfig() error {
	var err error
	config, err = util.LoadConfig(".")
	if err != nil {
		log.Fatal("can not loadconfig", err)
	}
	return nil
}

func main() {

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
