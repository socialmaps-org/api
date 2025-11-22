package method

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/model"
)

// TestDeleteReviewAuthorizationMissing tests that one cannot call the endpoint
// without authorization credentials.
func TestDeleteReviewAuthorizationMissing(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	usr := model.UpsertUser(ctx, db, "1", "Steve")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	rvw := model.CreateReview(ctx, db, plc.ID, usr.ID, true, "I like it!")

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, db, "")

	// Act
	req, err := http.NewRequest("DELETE", srv.URL+"/v1/reviews/"+rvw.ID, nil)
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

	db := database.Open(":memory:")
	usr := model.UpsertUser(ctx, db, "1", "Steve")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	rvw := model.CreateReview(ctx, db, plc.ID, usr.ID, true, "I like it!")

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, db, "")

	// Act
	req, err := http.NewRequest("DELETE", srv.URL+"/v1/reviews/"+rvw.ID, nil)
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

	db := database.Open(":memory:")
	usrA := model.UpsertUser(ctx, db, "1", "Alice")
	usrB := model.UpsertUser(ctx, db, "2", "Bob")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	rvw := model.CreateReview(ctx, db, plc.ID, usrA.ID, true, "I like it!")

	authr := NewTestAuthenticator(t, fmt.Sprint(usrB.OSMID))
	srv := NewTestServer(t, authr, db, "")

	// Act
	req, err := http.NewRequest("DELETE", srv.URL+"/v1/reviews/"+rvw.ID, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusForbidden, res.StatusCode)
	require.NotNil(t, model.LoadReview(ctx, db, rvw.ID))
}

// TestDeleteReview tests that a user (with an active authorization) can delete
// their own review.
func TestDeleteReview(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	usr := model.UpsertUser(ctx, db, "1", "Steve")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	rvw := model.CreateReview(ctx, db, plc.ID, usr.ID, true, "I like it!")

	authr := NewTestAuthenticator(t, fmt.Sprint(usr.OSMID))
	srv := NewTestServer(t, authr, db, "")

	// Act
	req, err := http.NewRequest("DELETE", srv.URL+"/v1/reviews/"+rvw.ID, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNoContent, res.StatusCode)
	require.Nil(t, model.LoadReview(ctx, db, rvw.ID))
}
