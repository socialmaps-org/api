package method

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/resource"
)

func TestCreateReviewAuthorizationMissing(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, db, "")

	// Act
	req, err := http.NewRequest(
		"POST", srv.URL+"/v1/places/"+plc.ID+"/reviews",
		strings.NewReader(`{"liked": true, "comments": "I liked it!"}`),
	)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Zero(t, authr.nIntrospectCalls)

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestCreateReviewMissingPlace(t *testing.T) {
	// Arrange
	db := database.Open(":memory:")

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, db, "")

	// Act
	req, err := http.NewRequest(
		"POST", srv.URL+"/v1/places/plc_missing/reviews",
		strings.NewReader(`{"liked": true, "comment": "I liked it!"}`),
	)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestCreateReview(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, db, "")

	// Act
	req, err := http.NewRequest(
		"POST", srv.URL+"/v1/places/"+plc.ID+"/reviews",
		strings.NewReader(`{"liked": true, "comment": "I liked it!"}`),
	)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusOK, res.StatusCode)

	var rvwR resource.Review
	err = json.NewDecoder(res.Body).Decode(&rvwR)
	require.NoError(t, err)
	require.Equal(t, plc.ID, rvwR.Place.ID)
	require.True(t, rvwR.Liked)
	require.Equal(t, "I liked it!", rvwR.Comment)

	rvwM := model.LoadReview(ctx, db, rvwR.ID)
	require.NotNil(t, rvwM)
	require.True(t, rvwM.Liked)
	require.Equal(t, "I liked it!", rvwM.Comment)
}
