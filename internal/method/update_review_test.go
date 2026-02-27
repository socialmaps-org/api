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
	"codeberg.org/socialmaps/api/internal/mytime"
)

func TestUpdateReviewAuthorizationMissing(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, new("I liked it!"), mytime.Now()))

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
	qs := model.New(database.OpenInTest(t))

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
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, new("I liked it!"), mytime.Now()))
	must.Get(qs.CreateReviewDecision(ctx, now, rvw.ID, "test-mod", true, ""))

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

	var rvwR any
	err = json.NewDecoder(res.Body).Decode(&rvwR)
	require.NoError(t, err)
	require.Equal(t, plc.ID, j.Get[int64](rvwR, "place", "id"))
	require.False(t, j.Get[bool](rvwR, "liked"))
	require.Equal(t, "I didn't like it!", j.Get[string](rvwR, "comment"))

	rvwM := must.Get(qs.LoadReview(ctx, j.Get[int64](rvwR, "id")))
	require.NotNil(t, rvwM)
	require.False(t, rvwM.Liked)
	require.Equal(t, "I didn't like it!", *rvwM.Comment)
	require.Nil(t, rvwM.LastDecisionApproved)
	require.Nil(t, rvwM.LastDecisionBy)
}

func TestUpdateReviewCommentAndLike(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, false, new("I didn't like it!"), mytime.Now()))
	must.Get(qs.CreateReviewDecision(ctx, now, rvw.ID, "test-mod", true, ""))

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

	var rvwR any
	err = json.NewDecoder(res.Body).Decode(&rvwR)
	require.NoError(t, err)
	require.Equal(t, plc.ID, j.Get[int64](rvwR, "place", "id"))
	require.True(t, j.Get[bool](rvwR, "liked"))
	require.Equal(t, "I liked it!", j.Get[string](rvwR, "comment"))

	rvwM := must.Get(qs.LoadReview(ctx, j.Get[int64](rvwR, "id")))
	require.NotNil(t, rvwM)
	require.True(t, rvwM.Liked)
	require.Equal(t, "I liked it!", *rvwM.Comment)
	require.Nil(t, rvwM.LastDecisionApproved)
	require.Nil(t, rvwM.LastDecisionBy)
}

func TestUpdateReviewCommentOnly(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, false, new("I liked it!"), mytime.Now()))
	must.Get(qs.CreateReviewDecision(ctx, now, rvw.ID, "test-mod", true, ""))

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

	var rvwR any
	err = json.NewDecoder(res.Body).Decode(&rvwR)
	require.NoError(t, err)
	require.Equal(t, plc.ID, j.Get[int64](rvwR, "place", "id"))
	require.False(t, j.Get[bool](rvwR, "liked"))
	require.Equal(t, "I didn't like it!", j.Get[string](rvwR, "comment"))

	rvwM := must.Get(qs.LoadReview(ctx, j.Get[int64](rvwR, "id")))
	require.NotNil(t, rvwM)
	require.False(t, rvwM.Liked)
	require.Equal(t, "I didn't like it!", *rvwM.Comment)
	require.Nil(t, rvwM.LastDecisionApproved)
	require.Nil(t, rvwM.LastDecisionBy)
}

func TestUpdateReviewLikeOnly(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, false, new("I liked it!"), mytime.Now()))

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

	var rvwR any
	err = json.NewDecoder(res.Body).Decode(&rvwR)
	require.NoError(t, err)
	require.Equal(t, plc.ID, j.Get[int64](rvwR, "place", "id"))
	require.True(t, j.Get[bool](rvwR, "liked"))
	require.Equal(t, "I liked it!", j.Get[string](rvwR, "comment"))

	rvwM := must.Get(qs.LoadReview(ctx, j.Get[int64](rvwR, "id")))
	require.NotNil(t, rvwM)
	require.True(t, rvwM.Liked)
	require.Equal(t, "I liked it!", *rvwM.Comment)
}

func TestUpdateReviewDislikeOnly(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	qs := model.New(database.OpenInTest(t))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096, mytime.Now()))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, true, new("I didn't like it!"), mytime.Now()))

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

	var rvwR any
	err = json.NewDecoder(res.Body).Decode(&rvwR)
	require.NoError(t, err)
	require.Equal(t, plc.ID, j.Get[int64](rvwR, "place", "id"))
	require.False(t, j.Get[bool](rvwR, "liked"))
	require.Equal(t, "I didn't like it!", j.Get[string](rvwR, "comment"))

	rvwM := must.Get(qs.LoadReview(ctx, j.Get[int64](rvwR, "id")))
	require.NotNil(t, rvwM)
	require.False(t, rvwM.Liked)
	require.Equal(t, "I didn't like it!", *rvwM.Comment)
}
