package models

import (
	"database/sql"
	"fmt"

	"github.com/k88t76/CodeArchives-server/config"

	// db
	_ "github.com/go-sql-driver/mysql"
)

const (
	tableNameArchives = "archives"
	tableNameUsers    = "users"
	tableNameSessions = "sessions"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open(config.Config.SQLDriver, fmt.Sprintf("root:%s@unix(/cloudsql/%s)/code_archives?parseTime=true", config.Config.Dbpass, config.Config.CloudSQL))

	/*db, err = sql.Open(config.Config.SQLDriver, config.Config.DbAccess+"?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		fmt.Println(err)
	}
	*/
	if err != nil {
		panic(err)
	}

	name := config.Config.DbName

	cmd := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", name)
	_, err = db.Exec(cmd)
	if err != nil {
		fmt.Println(err)
	}

	/*
		cmd = fmt.Sprintf("USE %s", name)
		_, err = db.Exec(cmd)
		if err != nil {
			fmt.Println(err)
		}
	*/

	// create notebooksTable
	cmd = fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
			id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
			uuid VARCHAR(36) NOT NULL,
			content TEXT,
			title VARCHAR(255),
			author VARCHAR(255),
			language VARCHAR(255),
			created_at DATETIME)`, tableNameArchives)
	db.Exec(cmd)

	// create usersTable
	cmd = fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
			id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
			uuid VARCHAR(36) NOT NULL,
			name VARCHAR(255),
			password VARCHAR(255),
			created_at DATETIME)`, tableNameUsers)
	db.Exec(cmd)

	cmd = fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
			id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
			uuid VARCHAR(36) NOT NULL,
			user_id VARCHAR(36),
			user_name VARCHAR(255),
			created_at DATETIME)`, tableNameSessions)
	db.Exec(cmd)
}
