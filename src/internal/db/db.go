package db

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

//go:embed schema.sql
var dbSchema embed.FS

func Init(dbPath string) error {
	connStr := fmt.Sprintf("file:%s?mode=rwc&_foreign_keys=true&", dbPath)

	var err error
	db, err = sql.Open("sqlite3", connStr)
	if err != nil {
		return err
	}

	//Test DB connection
	if err = db.Ping(); err != nil {
		return err
	}

	//Create tables if not exist
	sql, err := dbSchema.ReadFile("schema.sql")
	if err != nil {
		return err
	}
	if _, err = db.Exec(string(sql)); err != nil {
		return err
	}

	//Create initial user if not exists
	insertSql := "INSERT OR IGNORE INTO users (id, username, password) VALUES (?,?,?)"
	initialPw, err := bcrypt.GenerateFromPassword([]byte("admin"), 14)

	if err != nil {
		return err
	}
	_, err = db.Exec(insertSql, 1, "admin", string(initialPw))
	if err != nil {
		return err
	}

	return nil
}
