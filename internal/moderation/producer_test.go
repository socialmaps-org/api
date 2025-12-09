package moderation

import (
	"testing"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/model"
	"github.com/stretchr/testify/require"
)

func TestProduce(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	usr := model.UpsertUser(ctx, db, "1", "Steve")
	rvw := model.CreateReview(ctx, db, plc.ID, usr.ID, true, "great little cafe!")

	ch := make(chan *model.Review, 1)

	// Act
	next := produce(ctx, db, ch, model.EarliestID("rvw").String())

	// Assert
	require.Equal(t, rvw.ID, next)
	act := <-ch
	require.Equal(t, rvw.ID, act.ID)
}
