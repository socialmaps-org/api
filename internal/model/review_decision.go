package model

import (
	"context"
	"database/sql"
	"time"

	"codeberg.org/socialmaps/api/internal/database"
)

type ReviewDecision struct {
	Created   time.Time
	ID        string
	ReviewID  string
	Moderator string
	Approved  bool
	Details   string
}

const reviewDecisionColumns = `
	  created
	, id
	, review_id
    , moderator
    , approved
    , details
`

func scanReviewDecision(scn database.Scanner) *ReviewDecision {
	var created int64
	var dec ReviewDecision
	err := scn.Scan(
		&created,
		&dec.ID,
		&dec.ReviewID,
		&dec.Moderator,
		&dec.Approved,
		&dec.Details,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}

	dec.Created = time.Unix(created, 0)

	return &dec
}

func CreateReviewDecision(ctx context.Context, db *sql.DB, reviewID, moderator string, approved bool, details string) *ReviewDecision {
	id := NewRandomID("dec")

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	row := tx.QueryRowContext(ctx, `
		INSERT INTO ReviewDecisions (
			  id
			, review_id
			, moderator
			, approved
			, details
		) VALUES (
		 	  @id
			, @review_id
			, @moderator
			, @approved
			, @details
		) RETURNING `+reviewDecisionColumns+`;`,
		sql.Named("id", id.String()),
		sql.Named("review_id", reviewID),
		sql.Named("moderator", moderator),
		sql.Named("approved", approved),
		sql.Named("details", details),
	)

	dec := scanReviewDecision(row)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return dec
}

func LoadLatestDecisionOfReview(ctx context.Context, db *sql.DB, reviewID string) *ReviewDecision {
	row := db.QueryRowContext(ctx, `
		SELECT `+reviewDecisionColumns+`
		FROM ReviewDecisions
		WHERE
			review_id = @review_id
		ORDER BY
			created DESC
		LIMIT 1
		;`,
		sql.Named("review_id", reviewID),
	)

	dec := scanReviewDecision(row)

	return dec
}
