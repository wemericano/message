package config

import (
	"gopkg.in/ini.v1"
)

var cfg *ini.File

func Init() error {
	var err error
	cfg, err = ini.Load("config/conf.ini")
	return err
}

func GetPort() string {
	return cfg.Section("web").Key("port").String()
}

func GetDBHost() string {
	return cfg.Section("database").Key("host").String()
}

func GetDBPort() string {
	return cfg.Section("database").Key("port").String()
}

func GetDBUser() string {
	return cfg.Section("database").Key("user").String()
}

func GetDBPassword() string {
	return cfg.Section("database").Key("password").String()
}

func GetDBName() string {
	return cfg.Section("database").Key("dbname").String()
}
