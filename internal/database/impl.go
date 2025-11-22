package database

import (
	"database/sql"
	_ "embed"
	"log"

	"codeberg.org/socialmaps/api/internal/mytime"
	"github.com/mattn/go-sqlite3"
)

//go:embed "migrations/000 - initial schema.sql"
var migration0 string

func init() {
	sql.Register("sqlite3_extended", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.RegisterFunc("my_unixepoch", func() int64 {
				return mytime.Now().Unix()
			}, false)
		},
	})
}

func Open(dataSourceName string) *sql.DB {
	db, err := sql.Open("sqlite3_extended", dataSourceName)
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

	tx.Exec("PRAGMA foreign_keys = ON;")

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
