package moderation

import (
	"database/sql"
	"testing"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/must"
	"github.com/stretchr/testify/require"
)

func TestProduce(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	qs := model.New(db)
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096))
	usr := must.Get(qs.CreateUser(ctx, 1, "Steve"))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, sql.NullString{String: "great little cafe!", Valid: true}))

	ch := make(chan model.Review, 1)

	// Act
	nextID, nextCreated := produce(ctx, qs, ch, -1, -1)

	// Assert
	require.Equal(t, rvw.ID, nextID)
	require.Equal(t, rvw.Created, nextCreated)
	act := <-ch
	require.Equal(t, rvw, act)
}
