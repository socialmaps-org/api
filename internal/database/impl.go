package database

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.socialmaps.org/api/internal/env"
	"golang.socialmaps.org/api/internal/must"
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
	ctx := context.Background()
	pool := Open(env.Var.DatabaseDSN)
	tx := must.Get(pool.Begin(ctx))
	t.Cleanup(func() {
		must.Do(tx.Rollback(ctx))
	})
	return tx
}
