package moderation

import (
	"testing"
	"time"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/must"
	"codeberg.org/socialmaps/api/internal/mytime"
	"github.com/stretchr/testify/require"
)

func TestProduce(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	db := database.OpenInTest(t)
	qs := model.New(db)
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, new("great little cafe!"), mytime.Now()))

	ch := make(chan model.Review, 1)

	// Act
	nextID, nextCreated := produce(ctx, qs, ch, -1, time.Unix(0, 0))

	// Assert
	require.Equal(t, rvw.ID, nextID)
	require.Equal(t, rvw.Created, nextCreated)
	act := <-ch
	require.Equal(t, rvw, act)
}
