package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

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

	// 디버깅: 설정 값 확인 (비밀번호 제외)
	log.Printf("### DB Config - Host: %s, Port: %s, User: %s, DB: %s", host, port, user, dbname)

	// 빈 값 체크
	if host == "" || user == "" || password == "" || dbname == "" {
		log.Fatalf("### Database configuration is incomplete. Host: %s, User: %s, DB: %s", host, user, dbname)
	}

	// 비밀번호에 특수문자가 있을 수 있으므로 URL 인코딩
	encodedPassword := url.QueryEscape(password)

	// 연결 문자열 형식 개선
	connString := fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=%s;encrypt=disable;connection timeout=30",
		host, port, user, encodedPassword, dbname)

	log.Printf("### Attempting to connect to database...")

	var err error
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("### Failed to open database: %v", err)
	}

	// 연결 풀 설정
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 연결 테스트
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("### Pinging database...")
	if err = db.PingContext(ctx); err != nil {
		log.Fatalf("### Failed to connect to database: %v\n### Connection string: server=%s;port=%s;user id=%s;database=%s",
			err, host, port, user, dbname)
	}

	log.Println("### Database connection successful")
}
