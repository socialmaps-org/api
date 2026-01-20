package moderation

import (
	"context"
	"log/slog"

	"codeberg.org/socialmaps/api/internal/model"
)

func Consumer(ctx context.Context, qs *model.Queries, mod Moderator, ch <-chan model.Review) {
	for {
		consume(ctx, qs, mod, ch)
	}
}

func consume(ctx context.Context, qs *model.Queries, mod Moderator, ch <-chan model.Review) {
	rvw := <-ch

	if !rvw.Comment.Valid {
		return
	}

	dec, err := mod.Moderate(rvw.Comment.String)
	if err != nil {
		slog.ErrorContext(ctx, "CANONICAL-MODERATION-CONSUMER-LINE",
			"review", rvw.ID,
			"status", "error",
			"error", err.Error(),
		)
	}

	decM, err := qs.CreateReviewDecision(ctx, rvw.ID, mod.ID(), dec.Approved, dec.Details)
	if err != nil {
		slog.ErrorContext(ctx, "CANONICAL-MODERATION-CONSUMER-LINE",
			"review", rvw.ID,
			"status", "error",
			"error", err.Error(),
		)
		return
	}

	slog.InfoContext(ctx, "CANONICAL-MODERATION-CONSUMER-LINE",
		"review", rvw.ID,
		"status", "success",
		"decision_id", decM.ID,
		"approved", dec.Approved,
	)
}
