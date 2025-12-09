package moderation

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"codeberg.org/socialmaps/api/internal/model"
)

func Producer(ctx context.Context, db *sql.DB, ch chan<- *model.Review) {
	next := model.EarliestID("plc").String()

	for {
		next = produce(ctx, db, ch, next)
	}
}

func produce(ctx context.Context, db *sql.DB, ch chan<- *model.Review, next string) string {
	reviews := model.ListEarliestUnapprovedReviews(ctx, db, 100, next)
	for _, rvw := range reviews {
		ch <- rvw
	}
	if len(reviews) > 0 {
		next = reviews[len(reviews)-1].ID
		slog.InfoContext(ctx, "CANONICAL-MODERATION-PRODUCER-LINE",
			"n_reviews", len(reviews),
			"next", next,
		)
	} else {
		<-time.Tick(15 * time.Second)
	}

	return next
}
