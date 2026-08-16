package method

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"golang.socialmaps.org/api/internal/database"
	"golang.socialmaps.org/api/internal/j"
	"golang.socialmaps.org/api/internal/model"
	"golang.socialmaps.org/api/internal/must"
	"golang.socialmaps.org/api/internal/mytime"
)

func TestCreateReviewAuthorizationMissing(t *testing.T) {
	// Arrange
	ctx := t.Context()
	qs := model.New(database.OpenInTest(t))
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest(
		"POST", fmt.Sprintf("%s/v1/places/%d/reviews", srv.URL, plc.ID),
		strings.NewReader(`{"rating": 4, "comments": "I liked it!"}`),
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
	qs := model.New(database.OpenInTest(t))

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest(
		"POST", srv.URL+"/v1/places/42/reviews",
		strings.NewReader(`{"rating": 4, "comment": "I liked it!"}`),
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

	qs := model.New(database.OpenInTest(t))
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest(
		"POST", fmt.Sprintf("%s/v1/places/%d/reviews", srv.URL, plc.ID),
		strings.NewReader(`{"rating": 4, "comment": "I liked it!"}`),
	)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	resB, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusAccepted, res.StatusCode, string(resB))

	var rvwR any
	err = json.Unmarshal(resB, &rvwR)
	require.NoError(t, err)
	require.Equal(t, plc.ID, j.Get[int64](rvwR, "place", "id"))
	require.Equal(t, 4, j.Get[int](rvwR, "rating"))
	require.Equal(t, "I liked it!", j.Get[string](rvwR, "comment"))

	rvwM := must.Get(qs.LoadReview(ctx, j.Get[int64](rvwR, "id")))
	require.NotNil(t, rvwM)
	require.Equal(t, int32(4), rvwM.Rating)
	require.Equal(t, "I liked it!", *rvwM.Comment)
}

func TestCreateReviewInFuture(t *testing.T) {
	// Arrange
	ctx := t.Context()

	qs := model.New(database.OpenInTest(t))
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest(
		"POST", fmt.Sprintf("%s/v1/places/%d/reviews", srv.URL, plc.ID),
		strings.NewReader(`{"rating": 4, "comment": "I liked it!", "reviewed_at": 2524608000}`),
	)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	var errR any
	err = json.NewDecoder(res.Body).Decode(&errR)
	require.NoError(t, err)
	require.Equal(t, int64(http.StatusBadRequest), j.Get[int64](errR, "status"))
	require.Equal(t, "reviewed_at cannot be in the future", j.Get[string](errR, "detail"))
	require.Equal(t, 1, len(j.Get[[]any](errR, "errors")))
	require.Equal(t, "body.reviewed_at", j.Get[string](errR, "errors", 0, "location"))
	require.Equal(t, "reviewed_at cannot be in the future", j.Get[string](errR, "errors", 0, "message"))
	require.Equal(t, int64(2524608000), j.Get[int64](errR, "errors", 0, "value"))
}
