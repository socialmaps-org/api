package method

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/must"
	"codeberg.org/socialmaps/api/internal/mytime"
)

func TestUnlikeReviewAuthorizationMissing(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usrA := must.Get(qs.CreateUser(ctx, now, 1, "Alice"))
	usrB := must.Get(qs.CreateUser(ctx, now, 2, "Bob"))
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usrA.ID, true, new("I like it!"), mytime.Now()))
	must.Do(qs.LikeReview(ctx, rvw.ID, usrB.ID, mytime.Now()))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/reviews/%d/unlike", srv.URL, rvw.ID), nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Zero(t, authr.nIntrospectCalls)

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	rvw = must.Get(qs.LoadReview(ctx, rvw.ID))
	require.Equal(t, int64(1), rvw.NLikes)
}

func TestUnlikeReviewAuthorizationInactive(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usrA := must.Get(qs.CreateUser(ctx, now, 1, "Alice"))
	usrB := must.Get(qs.CreateUser(ctx, now, 2, "Bob"))
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usrA.ID, true, new("I like it!"), mytime.Now()))
	must.Do(qs.LikeReview(ctx, rvw.ID, usrB.ID, mytime.Now()))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/reviews/%d/unlike", srv.URL, rvw.ID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	rvw = must.Get(qs.LoadReview(ctx, rvw.ID))
	require.Equal(t, int64(1), rvw.NLikes)
}

func TestUnlikeReviewOthers(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usrA := must.Get(qs.CreateUser(ctx, now, 1, "Alice"))
	usrB := must.Get(qs.CreateUser(ctx, now, 2, "Bob"))
	usrC := must.Get(qs.CreateUser(ctx, now, 3, "Charlie"))
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usrA.ID, true, new("I like it!"), mytime.Now()))
	must.Do(qs.LikeReview(ctx, rvw.ID, usrB.ID, mytime.Now()))

	authr := NewTestAuthenticator(t, fmt.Sprint(usrC.ID))
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/reviews/%d/unlike", srv.URL, rvw.ID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = must.Get(qs.LoadReview(ctx, rvw.ID))
	require.Equal(t, int64(1), rvw.NLikes)
}

func TestUnlikeReview(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usrA := must.Get(qs.CreateUser(ctx, now, 1, "Alice"))
	usrB := must.Get(qs.CreateUser(ctx, now, 2, "Bob"))
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usrA.ID, true, new("I like it!"), mytime.Now()))
	must.Do(qs.LikeReview(ctx, rvw.ID, usrB.ID, mytime.Now()))

	authr := NewTestAuthenticator(t, fmt.Sprint(usrB.ID))
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/reviews/%d/unlike", srv.URL, rvw.ID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode, string(must.Get(io.ReadAll(res.Body))))

	rvw = must.Get(qs.LoadReview(ctx, rvw.ID))
	require.Zero(t, rvw.NLikes)
}

func TestUnlikeReviewIdempotent(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usrA := must.Get(qs.CreateUser(ctx, now, 1, "Alice"))
	usrB := must.Get(qs.CreateUser(ctx, now, 2, "Bob"))
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usrA.ID, true, new("I like it!"), mytime.Now()))
	must.Do(qs.LikeReview(ctx, rvw.ID, usrB.ID, mytime.Now()))

	authr := NewTestAuthenticator(t, fmt.Sprint(usrB.ID), fmt.Sprint(usrB.ID))
	srv := NewTestServer(t, authr, qs)

	// Act (#1)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/reviews/%d/unlike", srv.URL, rvw.ID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#1)
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = must.Get(qs.LoadReview(ctx, rvw.ID))
	require.Zero(t, rvw.NLikes)

	// Act (#2)
	req, err = http.NewRequest("PUT", fmt.Sprintf("%s/v1/reviews/%d/unlike", srv.URL, rvw.ID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#2)
	require.Equal(t, 2, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = must.Get(qs.LoadReview(ctx, rvw.ID))
	require.Zero(t, rvw.NLikes)
}
