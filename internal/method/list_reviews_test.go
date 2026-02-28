package method

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/require"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/j"
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/must"
	"codeberg.org/socialmaps/api/internal/mytime"
)

func TestListReviews(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

	mockClock := clock.NewMock()
	mytime.SetClockInTest(t, mockClock)

	qs := model.New(database.OpenInTest(t))
	usrA := must.Get(qs.CreateUser(ctx, now, 1, "Alice"))
	usrB := must.Get(qs.CreateUser(ctx, now, 2, "Bob"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", -8.4715779, 51.8952597, "node", 7095470096, mytime.Now()))
	rvwA := must.Get(qs.CreateReview(ctx, plc.ID, usrA.ID, true, new("I like it!"), mytime.Now()))
	must.Get(qs.CreateReviewDecision(ctx, now, rvwA.ID, "test-mod", true, ""))
	mockClock.Add(24 * time.Hour)
	rvwB := must.Get(qs.CreateReview(ctx, plc.ID, usrB.ID, false, new("I don't like it!"), mytime.Now()))
	must.Get(qs.CreateReviewDecision(ctx, now, rvwB.ID, "test-mod", true, ""))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs, "")

	// Act (#1)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/places/%d/reviews", srv.URL, plc.ID), nil)
	require.NoError(t, err)
	req.URL.RawQuery = url.Values{
		"limit": {"1"},
	}.Encode()

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#1)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var listR any
	err = json.NewDecoder(res.Body).Decode(&listR)
	require.NoError(t, err)
	require.Len(t, j.Get[[]any](listR, "data"), 1)
	require.Equal(t, rvwB.ID, j.Get[int64](listR, "data", 0, "id"))

	// Act (#2) [Next]
	req, err = http.NewRequest("GET", fmt.Sprintf("%s/v1/places/%d/reviews", srv.URL, plc.ID), nil)
	require.NoError(t, err)
	req.URL.RawQuery = url.Values{
		"limit":        {"1"},
		"last_id":      {fmt.Sprint(j.Get[int64](listR, "data", 0, "id"))},
		"last_created": {fmt.Sprint(j.Get[int64](listR, "data", 0, "created"))},
	}.Encode()

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#2) [Next]
	require.Equal(t, http.StatusOK, res.StatusCode)

	err = json.NewDecoder(res.Body).Decode(&listR)
	require.NoError(t, err)
	require.Len(t, j.Get[[]any](listR, "data"), 1)
	require.Equal(t, rvwA.ID, j.Get[int64](listR, "data", 0, "id"))

	// Act (#3) [Previous]
	req, err = http.NewRequest("GET", fmt.Sprintf("%s/v1/places/%d/reviews", srv.URL, plc.ID), nil)
	require.NoError(t, err)
	req.URL.RawQuery = url.Values{
		"limit":         {"1"},
		"first_id":      {fmt.Sprint(j.Get[int64](listR, "data", 0, "id"))},
		"first_created": {fmt.Sprint(j.Get[int64](listR, "data", 0, "created"))},
	}.Encode()

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#3) [Previous]
	require.Equal(t, http.StatusOK, res.StatusCode)

	err = json.NewDecoder(res.Body).Decode(&listR)
	require.NoError(t, err)
	require.Len(t, j.Get[[]any](listR, "data"), 1)
	require.Equal(t, rvwB.ID, j.Get[int64](listR, "data", 0, "id"))
}
