package method

import (
	"encoding/json"
	"fmt"
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

func TestUpdateReviewAuthorizationMissing(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, 4, new("I liked it!"), mytime.Now(), mytime.Now()))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest(
		"PUT", fmt.Sprintf("%s/v1/reviews/%d", srv.URL, rvw.ID),
		strings.NewReader(`{"rating": 2, "comments": "I didn't like it!"}`),
	)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Zero(t, authr.nIntrospectCalls)

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestUpdateReviewMissingReview(t *testing.T) {
	// Arrange
	qs := model.New(database.OpenInTest(t))

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest(
		"PUT", fmt.Sprintf("%s/v1/reviews/%d", srv.URL, 123),
		strings.NewReader(`{"rating": 2, "comment": "I didn't like it!"}`),
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

func TestUpdateReviewCommentAndRating(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, 4, new("I liked it!"), mytime.Now(), mytime.Now()))
	must.Get(qs.CreateReviewDecision(ctx, now, rvw.ID, "test-mod", true, ""))

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest(
		"PUT", fmt.Sprintf("%s/v1/reviews/%d", srv.URL, rvw.ID),
		strings.NewReader(`{"rating": 2, "comment": "I didn't like it!"}`),
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
	require.Equal(t, 2, j.Get[int](rvwR, "rating"))
	require.Equal(t, "I didn't like it!", j.Get[string](rvwR, "comment"))

	rvwM := must.Get(qs.LoadReview(ctx, j.Get[int64](rvwR, "id")))
	require.NotNil(t, rvwM)
	require.Equal(t, int32(2), rvwM.Rating)
	require.Equal(t, "I didn't like it!", *rvwM.Comment)
	require.Nil(t, rvwM.LastDecisionApproved)
	require.Nil(t, rvwM.LastDecisionBy)
}

func TestUpdateReviewCommentOnly(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, 2, new("I liked it!"), mytime.Now(), mytime.Now()))
	must.Get(qs.CreateReviewDecision(ctx, now, rvw.ID, "test-mod", true, ""))

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest(
		"PUT", fmt.Sprintf("%s/v1/reviews/%d", srv.URL, rvw.ID),
		strings.NewReader(`{"rating": 2, "comment": "I didn't like it!"}`),
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
	require.Equal(t, 2, j.Get[int](rvwR, "rating"))
	require.Equal(t, "I didn't like it!", j.Get[string](rvwR, "comment"))

	rvwM := must.Get(qs.LoadReview(ctx, j.Get[int64](rvwR, "id")))
	require.NotNil(t, rvwM)
	require.Equal(t, int32(2), rvwM.Rating)
	require.Equal(t, "I didn't like it!", *rvwM.Comment)
	require.Nil(t, rvwM.LastDecisionApproved)
	require.Nil(t, rvwM.LastDecisionBy)
}

func TestUpdateReviewRatingOnly(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, 2, new("I liked it!"), mytime.Now(), mytime.Now()))

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest(
		"PUT", fmt.Sprintf("%s/v1/reviews/%d", srv.URL, rvw.ID),
		strings.NewReader(`{"rating": 4, "comment": "I liked it!"}`),
	)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusOK, res.StatusCode)

	var rvwR any
	err = json.NewDecoder(res.Body).Decode(&rvwR)
	require.NoError(t, err)
	require.Equal(t, plc.ID, j.Get[int64](rvwR, "place", "id"))
	require.Equal(t, 4, j.Get[int](rvwR, "rating"))
	require.Equal(t, "I liked it!", j.Get[string](rvwR, "comment"))

	rvwM := must.Get(qs.LoadReview(ctx, j.Get[int64](rvwR, "id")))
	require.NotNil(t, rvwM)
	require.Equal(t, int32(4), rvwM.Rating)
	require.Equal(t, "I liked it!", *rvwM.Comment)
}
