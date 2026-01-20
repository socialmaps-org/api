package moderation

import (
	"context"
	"log/slog"
	"time"

	"codeberg.org/socialmaps/api/internal/model"
)

func Producer(ctx context.Context, qs *model.Queries, ch chan<- model.Review) {
	nextID, nextCreated := int64(-1), int64(-1)

	for {
		nextID, nextCreated = produce(ctx, qs, ch, nextID, nextCreated)
	}
}

func produce(ctx context.Context, qs *model.Queries, ch chan<- model.Review, nextID, nextCreated int64) (int64, int64) {
	reviews, err := qs.ListEarliestUnapprovedReviews(ctx, nextCreated, nextID, 100)
	if err != nil {
		slog.ErrorContext(ctx, "CANONICAL-MODERATION-PRODUCER-LINE",
			"status", "error",
			"error", err.Error(),
			"nextID", nextID,
			"nextCreated", nextCreated,
		)
		return nextID, nextCreated
	}

	for _, rvw := range reviews {
		ch <- rvw
	}

	slog.InfoContext(ctx, "CANONICAL-MODERATION-PRODUCER-LINE",
		"status", "success",
		"n_reviews", len(reviews),
		"nextID", nextID,
		"nextCreated", nextCreated,
	)

	if len(reviews) > 0 {
		lastReview := reviews[len(reviews)-1]
		return lastReview.ID, lastReview.Created
	} else {
		<-time.Tick(15 * time.Second)
		return nextID, nextCreated
	}
}
