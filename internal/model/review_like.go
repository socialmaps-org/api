package model

import "database/sql"

func LikeReview(db *sql.DB, reviewID, userID uint64) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	_, err = tx.Exec(
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

func UnlikeReview(db *sql.DB, reviewID, userID uint64) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	_, err = tx.Exec(
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
