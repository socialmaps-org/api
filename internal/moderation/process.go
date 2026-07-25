package moderation

import (
	"context"

	"golang.socialmaps.org/api/internal/model"
)

func Process(ctx context.Context, qs *model.Queries, mod Moderator) {
	ch := make(chan model.Review)

	go Producer(ctx, qs, ch)
	go Consumer(ctx, qs, mod, ch)
}
