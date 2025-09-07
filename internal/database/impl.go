package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

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

		_, err = tx.Exec(`
			CREATE TABLE Users (
				  id              INTEGER PRIMARY KEY
				, openid_provider TEXT NOT NULL
				, openid_subject  TEXT NOT NULL
				, username        TEXT NOT NULL
				, CONSTRAINT unique_identity UNIQUE (openid_provider, openid_subject) 
			) STRICT;
			CREATE TABLE Sessions (
				  id         INTEGER PRIMARY KEY
				, user_id    INTEGER REFERENCES Users NOT NULL
				, created_at INTEGER NOT NULL
				, revoked_at INTEGER
			) STRICT;
		`)
		if err != nil {
			panic(err)
		}
		_, err = tx.Exec("PRAGMA user_version = 1;")
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
