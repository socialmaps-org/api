package model

import (
	"context"
	"database/sql"

	"codeberg.org/socialmaps/auth/internal/database"
)

type Review struct {
	Model

	Created      int64
	Updated      int64
	ID           string
	PlaceID      string
	UserID       string
	Liked        bool
	Comment      string
	NLikes       uint64
	DecNLikes    uint64
	DecUpdatedAt *int64
}

const reviewColumns = `
	  created
	, updated
	, id
	, place_id
	, user_id
	, liked
	, comment
	, n_likes
	, dec_n_likes
	, dec_updated_at
`

func scanReview(scn database.Scanner) *Review {
	var rvw Review
	err := scn.Scan(
		&rvw.Created,
		&rvw.Updated,
		&rvw.ID,
		&rvw.PlaceID,
		&rvw.UserID,
		&rvw.Liked,
		&rvw.Comment,
		&rvw.NLikes,
		&rvw.DecNLikes,
		&rvw.DecUpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}

	return &rvw
}

func CreateReview(ctx context.Context, db *sql.DB, placeID, userID string, liked bool, comment string) *Review {
	id := randomID("rvw_")

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	row := tx.QueryRowContext(ctx, `
		INSERT INTO Reviews (
			  id
			, place_id
			, user_id
			, liked
			, comment
		) VALUES (
		 	  @id
		 	, @place_id
			, @user_id
			, @liked
			, @comment
		) RETURNING `+reviewColumns+`;`,
		sql.Named("id", id),
		sql.Named("place_id", placeID),
		sql.Named("user_id", userID),
		sql.Named("liked", liked),
		sql.Named("comment", comment),
	)

	rvw := scanReview(row)
	rvw.DB = db

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return rvw
}

func DeleteReview(ctx context.Context, db *sql.DB, id string) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	_, err = tx.ExecContext(ctx, `
		DELETE FROM Reviews
		WHERE
			id = @id
		;
		`,
		sql.Named("id", id),
	)
	if err != nil {
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func UpdateReview(ctx context.Context, db *sql.DB, id string, liked bool, comment string) *Review {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	row := tx.QueryRowContext(ctx,
		`
		UPDATE Reviews SET
			  liked   = @liked
			, comment = @comment
		WHERE
			id = @id
		RETURNING `+reviewColumns+`;`,
		sql.Named("liked", liked),
		sql.Named("comment", comment),
		sql.Named("id", id),
	)

	rvw := scanReview(row)
	rvw.DB = db

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return rvw
}

func ListLatestReviewsOfPlace(ctx context.Context, db *sql.DB, placeID string) []*Review {
	rows, err := db.QueryContext(ctx, `
		SELECT `+reviewColumns+`
		FROM Reviews
		WHERE
			place_id = @place_id
		ORDER BY
			created DESC
		`,
		sql.Named("place_id", placeID),
	)
	if err != nil {
		panic(err)
	}

	var reviews []*Review
	for rows.Next() {
		rvw := scanReview(rows)
		reviews = append(reviews, rvw)
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}

	return reviews
}

func LoadReview(ctx context.Context, db *sql.DB, id string) *Review {
	row := db.QueryRowContext(ctx, `
		SELECT `+reviewColumns+`
		FROM Reviews
		WHERE
			id = @id
		`,
		sql.Named("id", id),
	)

	rvw := scanReview(row)

	return rvw
}
