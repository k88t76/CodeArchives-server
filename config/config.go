package config

import (
	"log"
	"os"

	"github.com/go-ini/ini"
)

type ConfigList struct {
	LogFile   string
	CloudSQL  string
	DbAccess  string
	DbName    string
	Dbpass    string
	SQLDriver string
	Port      string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{
		LogFile:   cfg.Section("application").Key("log_file").String(),
		CloudSQL:  cfg.Section("db").Key("cloudSQL").String(),
		DbAccess:  cfg.Section("db").Key("access").String(),
		DbName:    cfg.Section("db").Key("name").String(),
		Dbpass:    cfg.Section("db").Key("pass").String(),
		SQLDriver: cfg.Section("db").Key("driver").String(),
		Port:      cfg.Section("web").Key("port").String(),
	}
}
