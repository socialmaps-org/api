package model

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"codeberg.org/socialmaps/api/internal/database"
)

type User struct {
	Created     time.Time
	Updated     time.Time
	ID          string
	OSMID       int
	DisplayName string
}

const userCols = `
	  created
	, updated
	, id
	, osm_id
	, display_name
`

func scanUser(scn database.Scanner) *User {
	var created, updated int64

	var usr User
	err := scn.Scan(
		&created,
		&updated,
		&usr.ID,
		&usr.OSMID,
		&usr.DisplayName,
	)
	if err != nil {
		panic(err)
	}

	usr.Created = time.Unix(created, 0)
	usr.Updated = time.Unix(updated, 0)

	return &usr
}

func UpsertUser(ctx context.Context, db *sql.DB, osmSub string, displayName string) *User {
	id := NewRandomID("usr")

	osmID, err := strconv.Atoi(osmSub)
	if err != nil {
		panic(err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	row := tx.QueryRowContext(ctx,
		`
		INSERT INTO Users (
			  created
			, id
			, osm_id
			, display_name
		) VALUES (
		 	  @created
		 	, @id
			, @osm_id
			, @display_name
		) ON CONFLICT DO UPDATE SET
		 	  display_name = @display_name
		WHERE
			osm_id = @osm_id
		RETURNING `+userCols+`;`,
		sql.Named("created", id.time.Unix()),
		sql.Named("id", id.String()),
		sql.Named("display_name", displayName),
		sql.Named("osm_id", osmID),
	)

	usr := scanUser(row)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return usr
}
