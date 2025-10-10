package model

import (
	"context"
	"database/sql"
)

func LikeReview(ctx context.Context, db *sql.DB, reviewID, userID string) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	_, err = tx.ExecContext(ctx,
		`
		INSERT INTO ReviewLikes (
			  review_id
			, user_id
		
		) VALUES (
			  @review_id
			, @user_id
		)
		ON CONFLICT DO NOTHING
		;
		`,
		sql.Named("review_id", reviewID),
		sql.Named("user_id", userID),
	)
	if err != nil {
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func UnlikeReview(ctx context.Context, db *sql.DB, reviewID, userID string) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	_, err = tx.ExecContext(ctx,
		`
		DELETE FROM ReviewLikes
		WHERE
			review_id = @review_id AND
			user_id = @user_id
		;
		`,
		sql.Named("review_id", reviewID),
		sql.Named("user_id", userID),
	)
	if err != nil {
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
