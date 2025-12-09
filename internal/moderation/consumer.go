package moderation

import (
	"context"
	"database/sql"
	"log/slog"

	"codeberg.org/socialmaps/api/internal/model"
)

func Consumer(ctx context.Context, db *sql.DB, mod Moderator, ch <-chan *model.Review) {
	for {
		consume(ctx, db, mod, ch)
	}
}

func consume(ctx context.Context, db *sql.DB, mod Moderator, ch <-chan *model.Review) {
	rvw := <-ch

	dec, err := mod.Moderate(rvw.Comment)
	if err != nil {
		slog.ErrorContext(ctx, "CANONICAL-MODERATION-CONSUMER-LINE",
			"review", rvw.ID,
			"status", "error",
			"error", err.Error(),
		)
	}

	decM := model.CreateReviewDecision(ctx, db, rvw.ID, mod.ID(), dec.Approved, dec.Details)

	slog.InfoContext(ctx, "CANONICAL-MODERATION-CONSUMER-LINE",
		"review", rvw.ID,
		"status", "success",
		"decision_id", decM.ID,
		"approved", dec.Approved,
	)
}
