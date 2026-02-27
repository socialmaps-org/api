package database

import (
	"context"
	"testing"

	"codeberg.org/socialmaps/api/internal/env"
	"codeberg.org/socialmaps/api/internal/must"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Open(dataSourceName string) *pgxpool.Pool {
	ctx := context.Background()
	db, err := pgxpool.New(ctx, dataSourceName)
	if err != nil {
		panic(err)
	}

	err = db.Ping(ctx)
	if err != nil {
		panic(err)
	}

	return db
}

func OpenInTest(t *testing.T) pgx.Tx {
	ctx := t.Context()
	pool := Open(env.Var.DatabaseDSN)
	tx := must.Get(pool.Begin(ctx))
	t.Cleanup(func() {
		tx.Rollback(ctx)
	})
	return tx
}
