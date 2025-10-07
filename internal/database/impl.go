package database

import (
	"database/sql"
	_ "embed"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed "migrations/000 - initial schema.sql"
var migration0 string

func Open(dataSourceName string) *sql.DB {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		panic(err)
	}

	initialize(db)

	return db
}

func initialize(db *sql.DB) {
	row := db.QueryRow("PRAGMA user_version;")
	var userVersion int
	row.Scan(&userVersion)

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	switch userVersion {
	case 0:
		log.Println("upgrading database schema 0 -> 1")

		_, err = tx.Exec(migration0)
		if err != nil {
			panic(err)
		}

		fallthrough
	case 1:
		log.Println("database has the latest schema")
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
