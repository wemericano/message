package config

import (
	"os"

	"gopkg.in/ini.v1"
)

var cfg *ini.File

func Init() error {
	// 환경 변수가 설정되어 있으면 환경 변수 사용, 없으면 파일 사용
	if os.Getenv("DB_HOST") != "" {
		return nil
	}

	var err error
	cfg, err = ini.Load("config/conf.ini")
	return err
}

func GetPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	if cfg != nil {
		return cfg.Section("web").Key("port").String()
	}
	return "4500"
}

func GetDBHost() string {
	if host := os.Getenv("DB_HOST"); host != "" {
		return host
	}
	if cfg != nil {
		return cfg.Section("database").Key("host").String()
	}
	return ""
}

func GetDBPort() string {
	if port := os.Getenv("DB_PORT"); port != "" {
		return port
	}
	if cfg != nil {
		return cfg.Section("database").Key("port").String()
	}
	return "1433"
}

func GetDBUser() string {
	if user := os.Getenv("DB_USER"); user != "" {
		return user
	}
	if cfg != nil {
		return cfg.Section("database").Key("user").String()
	}
	return ""
}

func GetDBPassword() string {
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		return password
	}
	if cfg != nil {
		return cfg.Section("database").Key("password").String()
	}
	return ""
}

func GetDBName() string {
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		return dbname
	}
	if cfg != nil {
		return cfg.Section("database").Key("dbname").String()
	}
	return ""
}
