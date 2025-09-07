package model

import (
	"database/sql"
	"time"
)

type Session struct {
	Model

	ID        int64
	UserID    int64
	CreatedAt time.Time
	RevokedAt *time.Time
}

func CreateSession(db *sql.DB, userID int64) *Session {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	row := tx.QueryRow(
		`
		INSERT INTO Sessions (
			  user_id
			, created_at
		) VALUES (
		 	  @user_id
			, unixepoch()
		) RETURNING
			  id
			, user_id
			, created_at
			, revoked_at
		;
		`,
		sql.Named("user_id", userID),
	)

	ses, err := scanSession(row)
	if err != nil {
		panic(err)
	}
	ses.DB = db

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return ses
}

func LoadActiveSession(db *sql.DB, id int64) *Session {
	row := db.QueryRow(
		`
		SELECT
			  id
			, user_id
			, created_at
			, revoked_at
		FROM
			Sessions
		WHERE
			id = @id
		;
		`,
		sql.Named("id", id),
	)
	ses, err := scanSession(row)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	if ses == nil {
		return nil
	}
	if ses.RevokedAt != nil && ses.RevokedAt.Before(time.Now()) {
		return nil
	}

	ses.DB = db

	return ses
}

func (ses *Session) Revoke() {
	tx, err := ses.DB.Begin()
	if err != nil {
		panic(err)
	}

	ses.DB.Exec(
		`
		UPDATE
			  revoked_at = unixepoch()
		FROM
			Sessions
		WHERE
			id = @id
		;
		`,
		sql.Named("id", ses.ID),
	)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func scanSession(row *sql.Row) (*Session, error) {
	var createdAt, revokedAt *int64
	var ses Session
	err := row.Scan(&ses.ID, &ses.UserID, &createdAt, &revokedAt)
	if err != nil {
		return nil, err
	}

	ses.CreatedAt = time.Unix(*createdAt, 0).UTC()
	if revokedAt != nil {
		t := time.Unix(*revokedAt, 0).UTC()
		ses.RevokedAt = &t
	}

	return &ses, nil
}
