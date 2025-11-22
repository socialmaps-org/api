package method

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/model"
)

func TestUnlikeReviewAuthorizationMissing(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	usrA := model.UpsertUser(ctx, db, "1", "Alice")
	usrB := model.UpsertUser(ctx, db, "2", "Bob")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	rvw := model.CreateReview(ctx, db, plc.ID, usrA.ID, true, "I like it!")
	model.LikeReview(ctx, db, rvw.ID, usrB.ID)

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, db, "")

	// Act
	req, err := http.NewRequest("PUT", srv.URL+"/v1/reviews/"+rvw.ID+"/unlike", nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Zero(t, authr.nIntrospectCalls)

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	rvw = model.LoadReview(ctx, db, rvw.ID)
	require.Equal(t, uint64(1), rvw.NLikes)
}

func TestUnlikeReviewAuthorizationInactive(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	usrA := model.UpsertUser(ctx, db, "1", "Alice")
	usrB := model.UpsertUser(ctx, db, "2", "Bob")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	rvw := model.CreateReview(ctx, db, plc.ID, usrA.ID, true, "I like it!")
	model.LikeReview(ctx, db, rvw.ID, usrB.ID)

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, db, "")

	// Act
	req, err := http.NewRequest("PUT", srv.URL+"/v1/reviews/"+rvw.ID+"/unlike", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	rvw = model.LoadReview(ctx, db, rvw.ID)
	require.Equal(t, uint64(1), rvw.NLikes)
}

func TestUnlikeReviewOthers(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	usrA := model.UpsertUser(ctx, db, "1", "Alice")
	usrB := model.UpsertUser(ctx, db, "2", "Bob")
	usrC := model.UpsertUser(ctx, db, "3", "Charlie")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	rvw := model.CreateReview(ctx, db, plc.ID, usrA.ID, true, "I like it!")
	model.LikeReview(ctx, db, rvw.ID, usrB.ID)

	authr := NewTestAuthenticator(t, fmt.Sprint(usrC.OSMID))
	srv := NewTestServer(t, authr, db, "")

	// Act
	req, err := http.NewRequest("PUT", srv.URL+"/v1/reviews/"+rvw.ID+"/unlike", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = model.LoadReview(ctx, db, rvw.ID)
	require.Equal(t, uint64(1), rvw.NLikes)
}

func TestUnlikeReview(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	usrA := model.UpsertUser(ctx, db, "1", "Alice")
	usrB := model.UpsertUser(ctx, db, "2", "Bob")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	rvw := model.CreateReview(ctx, db, plc.ID, usrA.ID, true, "I like it!")
	model.LikeReview(ctx, db, rvw.ID, usrB.ID)

	authr := NewTestAuthenticator(t, fmt.Sprint(usrB.OSMID))
	srv := NewTestServer(t, authr, db, "")

	// Act
	req, err := http.NewRequest("PUT", srv.URL+"/v1/reviews/"+rvw.ID+"/unlike", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = model.LoadReview(ctx, db, rvw.ID)
	require.Zero(t, rvw.NLikes)
}

func TestUnlikeReviewIdempotent(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	usrA := model.UpsertUser(ctx, db, "1", "Alice")
	usrB := model.UpsertUser(ctx, db, "2", "Bob")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	rvw := model.CreateReview(ctx, db, plc.ID, usrA.ID, true, "I like it!")
	model.LikeReview(ctx, db, rvw.ID, usrB.ID)

	authr := NewTestAuthenticator(t, fmt.Sprint(usrB.OSMID), fmt.Sprint(usrB.OSMID))
	srv := NewTestServer(t, authr, db, "")

	// Act (#1)
	req, err := http.NewRequest("PUT", srv.URL+"/v1/reviews/"+rvw.ID+"/unlike", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#1)
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = model.LoadReview(ctx, db, rvw.ID)
	require.Zero(t, rvw.NLikes)

	// Act (#2)
	req, err = http.NewRequest("PUT", srv.URL+"/v1/reviews/"+rvw.ID+"/unlike", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#2)
	require.Equal(t, 2, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = model.LoadReview(ctx, db, rvw.ID)
	require.Zero(t, rvw.NLikes)
}
