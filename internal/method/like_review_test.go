package method

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"codeberg.org/socialmaps/api/internal/must"
	"codeberg.org/socialmaps/api/internal/mytime"
	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/require"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/model"
)

func TestLikeReviewAuthorizationMissing(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, new("I like it!"), mytime.Now()))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs, "")

	// Act
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/reviews/%d/like", srv.URL, rvw.ID), nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Zero(t, authr.nIntrospectCalls)

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	rvw = must.Get(qs.LoadReview(ctx, rvw.ID))
	require.Zero(t, rvw.NLikes)
}

func TestLikeReviewAuthorizationInactive(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, new("I like it!"), mytime.Now()))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs, "")

	// Act
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/reviews/%d/like", srv.URL, rvw.ID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	rvw = must.Get(qs.LoadReview(ctx, rvw.ID))
	require.Zero(t, rvw.NLikes)
}

func TestLikeReviewSelf(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Alice"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, new("I like it!"), mytime.Now()))

	authr := NewTestAuthenticator(t, fmt.Sprint(usr.ID))
	srv := NewTestServer(t, authr, qs, "")

	// Act
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/reviews/%d/like", srv.URL, rvw.ID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusForbidden, res.StatusCode)

	rvw = must.Get(qs.LoadReview(ctx, rvw.ID))
	require.Zero(t, rvw.NLikes)
}

func TestLikeReview(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usrA := must.Get(qs.CreateUser(ctx, now, 1, "Alice"))
	usrB := must.Get(qs.CreateUser(ctx, now, 2, "Bob"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usrA.ID, true, new("I like it!"), mytime.Now()))

	authr := NewTestAuthenticator(t, fmt.Sprint(usrB.ID))
	srv := NewTestServer(t, authr, qs, "")

	// Act
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/reviews/%d/like", srv.URL, rvw.ID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = must.Get(qs.LoadReview(ctx, rvw.ID))
	require.Equal(t, int64(1), rvw.NLikes)
	require.Equal(t, float64(1), rvw.DecNLikes)
}

func TestLikeReviewIdempotent(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usrA := must.Get(qs.CreateUser(ctx, now, 1, "Alice"))
	usrB := must.Get(qs.CreateUser(ctx, now, 2, "Bob"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usrA.ID, true, new("I like it!"), mytime.Now()))

	authr := NewTestAuthenticator(t, fmt.Sprint(usrB.ID), fmt.Sprint(usrB.ID))
	srv := NewTestServer(t, authr, qs, "")

	// Act (#1)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/reviews/%d/like", srv.URL, rvw.ID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#1)
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = must.Get(qs.LoadReview(ctx, rvw.ID))
	require.Equal(t, int64(1), rvw.NLikes)
	require.Equal(t, float64(1), rvw.DecNLikes)

	// Act (#2)
	req, err = http.NewRequest("PUT", fmt.Sprintf("%s/v1/reviews/%d/like", srv.URL, rvw.ID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#2)
	require.Equal(t, 2, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = must.Get(qs.LoadReview(ctx, rvw.ID))
	require.Equal(t, int64(1), rvw.NLikes)
	require.Equal(t, float64(1), rvw.DecNLikes)
}

func TestLikeReviewDecay(t *testing.T) {
	// Arrange
	ctx := t.Context()

	mockClock := clock.NewMock()
	mytime.SetClockInTest(t, mockClock)

	qs := model.New(database.OpenInTest(t))
	usrA := must.Get(qs.CreateUser(ctx, mytime.Now(), 1, "Alice"))
	usrB := must.Get(qs.CreateUser(ctx, mytime.Now(), 2, "Bob"))
	usrC := must.Get(qs.CreateUser(ctx, mytime.Now(), 3, "Charlie"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usrA.ID, true, new("I like it!"), mytime.Now()))

	authr := NewTestAuthenticator(t, fmt.Sprint(usrB.ID), fmt.Sprint(usrC.ID))
	srv := NewTestServer(t, authr, qs, "")

	// Act (#1)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/reviews/%d/like", srv.URL, rvw.ID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#1)
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = must.Get(qs.LoadReview(ctx, rvw.ID))
	require.Equal(t, int64(1), rvw.NLikes)
	require.Equal(t, float64(1), rvw.DecNLikes)

	// Arrange (#2)
	mockClock.Add(180 * 24 * time.Hour)

	// Act (#2)
	req, err = http.NewRequest("PUT", fmt.Sprintf("%s/v1/reviews/%d/like", srv.URL, rvw.ID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#2)
	require.Equal(t, 2, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)

	rvw = must.Get(qs.LoadReview(ctx, rvw.ID))
	require.Equal(t, int64(2), rvw.NLikes)
	require.Equal(t, float64(1.5), rvw.DecNLikes)
}
