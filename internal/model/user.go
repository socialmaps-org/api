package model

import (
	"context"
	"database/sql"
)

type User struct {
	ID          string
	OIDProvider string
	OIDSubject  string
	Username    string
}

func CreateOrUpdateUser(ctx context.Context, db *sql.DB, oidProvider, oidSubject, username string) *User {
	id := randomID("usr_")

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	row := tx.QueryRowContext(ctx,
		`
		INSERT INTO Users (
			  id
			, openid_provider
			, openid_subject
			, username
		) VALUES (
		 	  @id
			, @openid_provider
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
		sql.Named("id", id),
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

func LoadUser(ctx context.Context, db *sql.DB, id string) *User {
	row := db.QueryRowContext(ctx,
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
