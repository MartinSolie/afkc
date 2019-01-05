package main

import (
	"fmt"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// sqlite3|postgres
	dbType = "sqlite3"

	sqliteDbPath = "./test.db"

	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "afkc"
)

func getConnectionString(db string) (dbInfo string, success bool) {
	if db == "postgres" {
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname,
		), true
	} else if db == "sqlite3" {
		return sqliteDbPath, true
	}
	return "", false
}

func connectToDb() *sql.DB {
	dbInfo, success := getConnectionString(dbType)
	if !success {
		panic(fmt.Sprintf("Unknown dbType: %s", dbType))
	}
	db, err := sql.Open(dbType, dbInfo)
	panicOnError(err)
	return db
}

