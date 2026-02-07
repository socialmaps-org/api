package method

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/must"
	"codeberg.org/socialmaps/api/internal/resource"
)

func TestUpdateReviewAuthorizationMissing(t *testing.T) {
	// Arrange
	ctx := t.Context()

	qs := model.New(database.Open(":memory:"))
	usr := must.Get(qs.CreateUser(ctx, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, sql.NullString{String: "I liked it!", Valid: true}))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs, "")

	// Act
	req, err := http.NewRequest(
		"PUT", fmt.Sprintf("%s/v1/reviews/%d", srv.URL, rvw.ID),
		strings.NewReader(`{"liked": false, "comments": "I didn't like it!"}`),
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
	qs := model.New(database.Open(":memory:"))

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, qs, "")

	// Act
	req, err := http.NewRequest(
		"PUT", fmt.Sprintf("%s/v1/reviews/%d", srv.URL, 123),
		strings.NewReader(`{"liked": false, "comment": "I didn't like it!"}`),
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

func TestUpdateReviewCommentAndDislike(t *testing.T) {
	// Arrange
	ctx := t.Context()

	qs := model.New(database.Open(":memory:"))
	usr := must.Get(qs.CreateUser(ctx, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, sql.NullString{String: "I liked it!", Valid: true}))
	must.Get(qs.CreateReviewDecision(ctx, rvw.ID, "test-mod", true, ""))

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, qs, "")

	// Act
	req, err := http.NewRequest(
		"PUT", fmt.Sprintf("%s/v1/reviews/%d", srv.URL, rvw.ID),
		strings.NewReader(`{"liked": false, "comment": "I didn't like it!"}`),
	)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusAccepted, res.StatusCode)

	var rvwR resource.Review
	err = json.NewDecoder(res.Body).Decode(&rvwR)
	require.NoError(t, err)
	require.Equal(t, plc.ID, rvwR.Place.ID)
	require.False(t, rvwR.Liked)
	require.Equal(t, "I didn't like it!", rvwR.Comment)

	rvwM := must.Get(qs.LoadReview(ctx, rvwR.ID))
	require.NotNil(t, rvwM)
	require.False(t, rvwM.Liked)
	require.Equal(t, "I didn't like it!", rvwM.Comment.String)
	require.False(t, rvwM.LastDecisionApproved.Valid)
	require.False(t, rvwM.LastDecisionBy.Valid)
}

func TestUpdateReviewCommentAndLike(t *testing.T) {
	// Arrange
	ctx := t.Context()

	qs := model.New(database.Open(":memory:"))
	usr := must.Get(qs.CreateUser(ctx, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, false, sql.NullString{String: "I didn't like it!", Valid: true}))
	must.Get(qs.CreateReviewDecision(ctx, rvw.ID, "test-mod", true, ""))

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, qs, "")

	// Act
	req, err := http.NewRequest(
		"PUT", fmt.Sprintf("%s/v1/reviews/%d", srv.URL, rvw.ID),
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

	var rvwR resource.Review
	err = json.NewDecoder(res.Body).Decode(&rvwR)
	require.NoError(t, err)
	require.Equal(t, plc.ID, rvwR.Place.ID)
	require.True(t, rvwR.Liked)
	require.Equal(t, "I liked it!", rvwR.Comment)

	rvwM := must.Get(qs.LoadReview(ctx, rvwR.ID))
	require.NotNil(t, rvwM)
	require.True(t, rvwM.Liked)
	require.Equal(t, "I liked it!", rvwM.Comment.String)
	require.False(t, rvwM.LastDecisionApproved.Valid)
	require.False(t, rvwM.LastDecisionBy.Valid)
}

func TestUpdateReviewCommentOnly(t *testing.T) {
	// Arrange
	ctx := t.Context()

	qs := model.New(database.Open(":memory:"))
	usr := must.Get(qs.CreateUser(ctx, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, false, sql.NullString{String: "I liked it!", Valid: true}))
	must.Get(qs.CreateReviewDecision(ctx, rvw.ID, "test-mod", true, ""))

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, qs, "")

	// Act
	req, err := http.NewRequest(
		"PUT", fmt.Sprintf("%s/v1/reviews/%d", srv.URL, rvw.ID),
		strings.NewReader(`{"liked": false, "comment": "I didn't like it!"}`),
	)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer my-auth-token")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, 1, authr.nIntrospectCalls)

	require.Equal(t, http.StatusAccepted, res.StatusCode)

	var rvwR resource.Review
	err = json.NewDecoder(res.Body).Decode(&rvwR)
	require.NoError(t, err)
	require.Equal(t, plc.ID, rvwR.Place.ID)
	require.False(t, rvwR.Liked)
	require.Equal(t, "I didn't like it!", rvwR.Comment)

	rvwM := must.Get(qs.LoadReview(ctx, rvwR.ID))
	require.NotNil(t, rvwM)
	require.False(t, rvwM.Liked)
	require.Equal(t, "I didn't like it!", rvwM.Comment.String)
	require.False(t, rvwM.LastDecisionApproved.Valid)
	require.False(t, rvwM.LastDecisionBy.Valid)
}

func TestUpdateReviewLikeOnly(t *testing.T) {
	// Arrange
	ctx := t.Context()

	qs := model.New(database.Open(":memory:"))
	usr := must.Get(qs.CreateUser(ctx, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, false, sql.NullString{String: "I liked it!", Valid: true}))

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, qs, "")

	// Act
	req, err := http.NewRequest(
		"PUT", fmt.Sprintf("%s/v1/reviews/%d", srv.URL, rvw.ID),
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

	rvwM := must.Get(qs.LoadReview(ctx, rvwR.ID))
	require.NotNil(t, rvwM)
	require.True(t, rvwM.Liked)
	require.Equal(t, "I liked it!", rvwM.Comment.String)
}

func TestUpdateReviewDislikeOnly(t *testing.T) {
	// Arrange
	ctx := t.Context()

	qs := model.New(database.Open(":memory:"))
	usr := must.Get(qs.CreateUser(ctx, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, sql.NullString{String: "I didn't like it!", Valid: true}))

	authr := NewTestAuthenticator(t, "1")
	srv := NewTestServer(t, authr, qs, "")

	// Act
	req, err := http.NewRequest(
		"PUT", fmt.Sprintf("%s/v1/reviews/%d", srv.URL, rvw.ID),
		strings.NewReader(`{"liked": false, "comment": "I didn't like it!"}`),
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
	require.False(t, rvwR.Liked)
	require.Equal(t, "I didn't like it!", rvwR.Comment)

	rvwM := must.Get(qs.LoadReview(ctx, rvwR.ID))
	require.NotNil(t, rvwM)
	require.False(t, rvwM.Liked)
	require.Equal(t, "I didn't like it!", rvwM.Comment.String)
}
