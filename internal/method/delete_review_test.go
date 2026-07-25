package method

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"golang.socialmaps.org/api/internal/database"
	"golang.socialmaps.org/api/internal/model"
	"golang.socialmaps.org/api/internal/must"
	"golang.socialmaps.org/api/internal/mytime"
)

// TestDeleteReviewAuthorizationMissing tests that one cannot call the endpoint
// without authorization credentials.
func TestDeleteReviewAuthorizationMissing(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, new("I like it!"), mytime.Now()))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/reviews/%d", srv.URL, rvw.ID), nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Zero(t, authr.nIntrospectCalls)

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

// TestDeleteReviewAuthorizationInactive tests that one cannot call the endpoint
// with an authorization that's not "active".
func TestDeleteReviewAuthorizationInactive(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, new("I like it!"), mytime.Now()))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/reviews/%d", srv.URL, rvw.ID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

// TestDeleteReviewOthers tests that a user cannot delete someone else's review.
func TestDeleteReviewOthers(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usrA := must.Get(qs.CreateUser(ctx, now, 1, "Alice"))
	usrB := must.Get(qs.CreateUser(ctx, now, 2, "Bob"))
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usrA.ID, true, new("I like it!"), mytime.Now()))

	authr := NewTestAuthenticator(t, fmt.Sprint(usrB.ID))
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/reviews/%d", srv.URL, rvw.ID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusForbidden, res.StatusCode)
	require.NotNil(t, must.Get(qs.LoadReview(ctx, rvw.ID)))
}

// TestDeleteReview tests that a user (with an active authorization) can delete
// their own review.
func TestDeleteReview(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, new("I like it!"), mytime.Now()))

	authr := NewTestAuthenticator(t, fmt.Sprint(usr.ID))
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/reviews/%d", srv.URL, rvw.ID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)
	rvw, err = qs.LoadReview(ctx, rvw.ID)
	require.Equal(t, pgx.ErrNoRows, err)
}
