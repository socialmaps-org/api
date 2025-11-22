package method

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"codeberg.org/socialmaps/api/internal/mytime"
	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/require"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/model"
)

func TestLikeReviewAuthorizationMissing(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	usr := model.UpsertUser(ctx, db, "1", "Steve")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	rvw := model.CreateReview(ctx, db, plc.ID, usr.ID, true, "I like it!")

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, db, "")

	// Act
	req, err := http.NewRequest("PUT", srv.URL+"/v1/reviews/"+rvw.ID+"/like", nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Zero(t, authr.nIntrospectCalls)

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	rvw = model.LoadReview(ctx, db, rvw.ID)
	require.Zero(t, rvw.NLikes)
}

func TestLikeReviewAuthorizationInactive(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	usr := model.UpsertUser(ctx, db, "1", "Steve")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	rvw := model.CreateReview(ctx, db, plc.ID, usr.ID, true, "I like it!")

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, db, "")

	// Act
	req, err := http.NewRequest("PUT", srv.URL+"/v1/reviews/"+rvw.ID+"/like", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	rvw = model.LoadReview(ctx, db, rvw.ID)
	require.Zero(t, rvw.NLikes)
}

func TestLikeReviewSelf(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	usr := model.UpsertUser(ctx, db, "1", "Alice")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	rvw := model.CreateReview(ctx, db, plc.ID, usr.ID, true, "I like it!")

	authr := NewTestAuthenticator(t, fmt.Sprint(usr.OSMID))
	srv := NewTestServer(t, authr, db, "")

	// Act
	req, err := http.NewRequest("PUT", srv.URL+"/v1/reviews/"+rvw.ID+"/like", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusForbidden, res.StatusCode)

	rvw = model.LoadReview(ctx, db, rvw.ID)
	require.Zero(t, rvw.NLikes)
}

func TestLikeReview(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	usrA := model.UpsertUser(ctx, db, "1", "Alice")
	usrB := model.UpsertUser(ctx, db, "2", "Bob")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	rvw := model.CreateReview(ctx, db, plc.ID, usrA.ID, true, "I like it!")

	authr := NewTestAuthenticator(t, fmt.Sprint(usrB.OSMID))
	srv := NewTestServer(t, authr, db, "")

	// Act
	req, err := http.NewRequest("PUT", srv.URL+"/v1/reviews/"+rvw.ID+"/like", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = model.LoadReview(ctx, db, rvw.ID)
	require.Equal(t, uint64(1), rvw.NLikes)
	require.Equal(t, float64(1), rvw.DecNLikes)
}

func TestLikeReviewIdempotent(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	usrA := model.UpsertUser(ctx, db, "1", "Alice")
	usrB := model.UpsertUser(ctx, db, "2", "Bob")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	rvw := model.CreateReview(ctx, db, plc.ID, usrA.ID, true, "I like it!")

	authr := NewTestAuthenticator(t, fmt.Sprint(usrB.OSMID), fmt.Sprint(usrB.OSMID))
	srv := NewTestServer(t, authr, db, "")

	// Act (#1)
	req, err := http.NewRequest("PUT", srv.URL+"/v1/reviews/"+rvw.ID+"/like", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#1)
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = model.LoadReview(ctx, db, rvw.ID)
	require.Equal(t, uint64(1), rvw.NLikes)
	require.Equal(t, float64(1), rvw.DecNLikes)

	// Act (#2)
	req, err = http.NewRequest("PUT", srv.URL+"/v1/reviews/"+rvw.ID+"/like", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#2)
	require.Equal(t, 2, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = model.LoadReview(ctx, db, rvw.ID)
	require.Equal(t, uint64(1), rvw.NLikes)
	require.Equal(t, float64(1), rvw.DecNLikes)
}

func TestLikeReviewDecay(t *testing.T) {
	// Arrange
	ctx := t.Context()

	mockClock := clock.NewMock()
	mytime.SetClock(mockClock)

	db := database.Open(":memory:")
	usrA := model.UpsertUser(ctx, db, "1", "Alice")
	usrB := model.UpsertUser(ctx, db, "2", "Bob")
	usrC := model.UpsertUser(ctx, db, "3", "Charlie")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	rvw := model.CreateReview(ctx, db, plc.ID, usrA.ID, true, "I like it!")

	authr := NewTestAuthenticator(t, fmt.Sprint(usrB.OSMID), fmt.Sprint(usrC.OSMID))
	srv := NewTestServer(t, authr, db, "")

	// Act (#1)
	req, err := http.NewRequest("PUT", srv.URL+"/v1/reviews/"+rvw.ID+"/like", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#1)
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = model.LoadReview(ctx, db, rvw.ID)
	require.Equal(t, uint64(1), rvw.NLikes)
	require.Equal(t, float64(1), rvw.DecNLikes)

	// Arrange (#2)
	mockClock.Add(180 * 24 * time.Hour)

	// Act (#2)
	req, err = http.NewRequest("PUT", srv.URL+"/v1/reviews/"+rvw.ID+"/like", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#2)
	require.Equal(t, 2, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = model.LoadReview(ctx, db, rvw.ID)
	require.Equal(t, uint64(2), rvw.NLikes)
	require.Equal(t, float64(1.5), rvw.DecNLikes)
}
