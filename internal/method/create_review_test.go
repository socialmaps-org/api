package method

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/j"
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/must"
)

func TestCreateReviewAuthorizationMissing(t *testing.T) {
	// Arrange
	ctx := t.Context()

	qs := model.New(database.Open(":memory:"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs, "")

	// Act
	req, err := http.NewRequest(
		"POST", fmt.Sprintf("%s/v1/places/%d/reviews", srv.URL, plc.ID),
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
	qs := model.New(database.Open(":memory:"))

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, qs, "")

	// Act
	req, err := http.NewRequest(
		"POST", srv.URL+"/v1/places/42/reviews",
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

	qs := model.New(database.Open(":memory:"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096))

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, qs, "")

	// Act
	req, err := http.NewRequest(
		"POST", fmt.Sprintf("%s/v1/places/%d/reviews", srv.URL, plc.ID),
		strings.NewReader(`{"liked": true, "comment": "I liked it!"}`),
	)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusAccepted, res.StatusCode)

	var rvwR any
	err = json.NewDecoder(res.Body).Decode(&rvwR)
	require.NoError(t, err)
	require.Equal(t, plc.ID, j.Get[int64](rvwR, "place", "id"))
	require.True(t, j.Get[bool](rvwR, "liked"))
	require.Equal(t, "I liked it!", j.Get[string](rvwR, "comment"))

	rvwM := must.Get(qs.LoadReview(ctx, j.Get[int64](rvwR, "id")))
	require.NotNil(t, rvwM)
	require.True(t, rvwM.Liked)
	require.Equal(t, "I liked it!", rvwM.Comment.String)
}
