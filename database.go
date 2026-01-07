package main

import (
	"database/sql"
	"log"

	"messanger/config"

	_ "github.com/denisenkom/go-mssqldb"
)

var db *sql.DB

// initDB 데이터베이스 초기화 및 연결
func initDB() {
	host := config.GetDBHost()
	port := config.GetDBPort()
	user := config.GetDBUser()
	password := config.GetDBPassword()
	dbname := config.GetDBName()

	connString := "server=" + host + ";port=" + port + ";user id=" + user + ";password=" + password + ";database=" + dbname + ";encrypt=disable"

	var err error
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("### Failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("### Failed to connect to database: %v", err)
	}

	log.Println("### Database connection successful")
}
