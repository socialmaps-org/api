package model

import (
	"database/sql"
)

type User struct {
	ID          int64
	OIDProvider string
	OIDSubject  string
	Username    string
}

func CreateOrUpdateUser(db *sql.DB, oidProvider, oidSubject, username string) *User {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	row := tx.QueryRow(
		`
		INSERT INTO Users (
			  openid_provider
			, openid_subject
			, username
		) VALUES (
			  @openid_provider
			, @openid_subject
			, @username
		) ON CONFLICT DO UPDATE SET
		 	  username = @username
		WHERE 
			openid_provider = @openid_provider AND
			openid_subject = @openid_subject
		RETURNING
			  id
			, openid_provider
			, openid_subject
			, username
		;
		`,
		sql.Named("openid_provider", oidProvider),
		sql.Named("openid_subject", oidSubject),
		sql.Named("username", username),
	)

	usr := scanUser(row)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return usr
}

func LoadUser(db *sql.DB, id int64) *User {
	row := db.QueryRow(
		`
		SELECT
			  id
			, openid_provider
			, openid_subject
			, username
		FROM
			Users
		WHERE
			id = @id
		;
		`,
		sql.Named("id", id),
	)
	usr := scanUser(row)
	return usr
}

func scanUser(row *sql.Row) *User {
	var usr User
	err := row.Scan(&usr.ID, &usr.OIDProvider, &usr.OIDSubject, &usr.Username)
	if err != nil {
		panic(err)
	}

	return &usr
}
