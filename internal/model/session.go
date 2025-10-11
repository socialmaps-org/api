package model

import (
	"context"
	"database/sql"
	"time"
)

type Session struct {
	Model

	ID        string
	UserID    string
	Created   time.Time
	Updated   time.Time
	RevokedAt *time.Time
}

const sessionColumns = `
	  created
	, updated
	, id
	, user_id
	, revoked_at
`

func scanSession(row *sql.Row) *Session {
	var ses Session
	var created, updated int64
	var revokedAt sql.NullInt64

	err := row.Scan(
		&created,
		&updated,
		&ses.ID,
		&ses.UserID,
		&revokedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}

	// Convert Unix timestamps to time.Time
	ses.Created = time.Unix(created, 0)
	ses.Updated = time.Unix(updated, 0)
	if revokedAt.Valid {
		t := time.Unix(revokedAt.Int64, 0)
		ses.RevokedAt = &t
	}

	return &ses
}

func CreateSession(ctx context.Context, db *sql.DB, userID string) *Session {
	id := randomID("ses_")

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	row := tx.QueryRowContext(ctx,
		`
		INSERT INTO Sessions (
			  id
			, user_id
		) VALUES (
		 	  @id
		 	, @user_id
		) RETURNING `+sessionColumns+`;`,
		sql.Named("id", id),
		sql.Named("user_id", userID),
	)

	ses := scanSession(row)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return ses
}

func LoadActiveSession(ctx context.Context, db *sql.DB, id string) *Session {
	row := db.QueryRowContext(ctx,
		`
		SELECT `+sessionColumns+`
		FROM
			Sessions
		WHERE
			id = @id
		;
		`,
		sql.Named("id", id),
	)
	ses := scanSession(row)
	if ses == nil {
		return nil
	}
	if ses.RevokedAt != nil && ses.RevokedAt.Before(time.Now()) {
		return nil
	}

	return ses
}

func RevokeSession(ctx context.Context, db *sql.DB, id string) *Session {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	row := db.QueryRowContext(ctx,
		`
		UPDATE Sessions SET
			  revoked_at = unixepoch()			
		WHERE
			id = @id
		RETURNING `+sessionColumns+`;`,
		sql.Named("id", id),
	)

	ses := scanSession(row)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return ses
}
