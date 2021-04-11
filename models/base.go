package models

import (
	"database/sql"
	"fmt"
	"log"

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

	/*
		// db connection for local
		db, err = sql.Open(config.Config.SQLDriver, config.Config.DbAccess+"code_archives?parseTime=true&loc=Asia%2FTokyo")
		if err != nil {
			fmt.Println(err)
		}
	*/

	name := config.Config.DbName

	cmd := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", name)
	_, err = db.Exec(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	cmd = fmt.Sprintf("USE %s", name)
	_, err = db.Exec(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	// create archivesTable
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
		uuid VARCHAR(36) NOT NULL, 
		content LONGTEXT NOT NULL,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		language VARCHAR(255),
		created_at VARCHAR(255))`, tableNameArchives)
	db.Exec(cmd)

	// create usersTable
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
		uuid VARCHAR(36) NOT NULL, 
		name VARCHAR(255),
		password VARCHAR(255),
		created_at VARCHAR(255))`, tableNameUsers)
	db.Exec(cmd)

	// create sessionsTable
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
		uuid VARCHAR(36) NOT NULL, 
		token VARCHAR(36),
		user_id VARCHAR(36),
		created_at VARCHAR(255))`, tableNameSessions)
	db.Exec(cmd)
}
