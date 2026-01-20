package method

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/must"
)

// TestDeleteReviewAuthorizationMissing tests that one cannot call the endpoint
// without authorization credentials.
func TestDeleteReviewAuthorizationMissing(t *testing.T) {
	// Arrange
	ctx := t.Context()

	qs := model.New(database.Open(":memory:"))
	usr := must.Get(qs.CreateUser(ctx, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, sql.NullString{String: "I like it!", Valid: true}))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs, "")

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

	qs := model.New(database.Open(":memory:"))
	usr := must.Get(qs.CreateUser(ctx, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, sql.NullString{String: "I like it!", Valid: true}))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs, "")

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

	qs := model.New(database.Open(":memory:"))
	usrA := must.Get(qs.CreateUser(ctx, 1, "Alice"))
	usrB := must.Get(qs.CreateUser(ctx, 2, "Bob"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usrA.ID, true, sql.NullString{String: "I like it!", Valid: true}))

	authr := NewTestAuthenticator(t, fmt.Sprint(usrB.ID))
	srv := NewTestServer(t, authr, qs, "")

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

	qs := model.New(database.Open(":memory:"))
	usr := must.Get(qs.CreateUser(ctx, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, sql.NullString{String: "I like it!", Valid: true}))

	authr := NewTestAuthenticator(t, fmt.Sprint(usr.ID))
	srv := NewTestServer(t, authr, qs, "")

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
	require.Equal(t, sql.ErrNoRows, err)
}
