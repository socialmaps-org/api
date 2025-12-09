package moderation

import (
	"context"
	"database/sql"

	"codeberg.org/socialmaps/api/internal/model"
)

func Process(ctx context.Context, db *sql.DB, mod Moderator) {
	ch := make(chan *model.Review)

	go Producer(ctx, db, ch)
	go Consumer(ctx, db, mod, ch)
}
