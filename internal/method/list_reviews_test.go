package method

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/require"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/must"
	"codeberg.org/socialmaps/api/internal/mytime"
	"codeberg.org/socialmaps/api/internal/resource"
)

func TestListReviews(t *testing.T) {
	// Arrange
	ctx := t.Context()

	mockClock := clock.NewMock()
	mytime.SetClock(mockClock)

	qs := model.New(database.Open(":memory:"))
	usrA := must.Get(qs.CreateUser(ctx, 1, "Alice"))
	usrB := must.Get(qs.CreateUser(ctx, 2, "Bob"))
	plc := must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096))
	rvwA := must.Get(qs.CreateReview(ctx, plc.ID, usrA.ID, true, sql.NullString{String: "I like it!", Valid: true}))
	must.Get(qs.CreateReviewDecision(ctx, rvwA.ID, "test-mod", true, ""))
	mockClock.Add(24 * time.Hour)
	rvwB := must.Get(qs.CreateReview(ctx, plc.ID, usrB.ID, false, sql.NullString{String: "I don't like it!", Valid: true}))
	must.Get(qs.CreateReviewDecision(ctx, rvwB.ID, "test-mod", true, ""))

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

	list := resource.List[resource.Review]{}
	err = json.NewDecoder(res.Body).Decode(&list)
	require.NoError(t, err)
	require.Len(t, list.Data, 1)
	require.Equal(t, rvwB.ID, list.Data[0].ID)

	// Act (#2) [Next]
	req, err = http.NewRequest("GET", fmt.Sprintf("%s/v1/places/%d/reviews", srv.URL, plc.ID), nil)
	require.NoError(t, err)
	req.URL.RawQuery = url.Values{
		"limit":        {"1"},
		"last_id":      {fmt.Sprint(list.Data[0].ID)},
		"last_created": {fmt.Sprint(list.Data[0].Created)},
	}.Encode()

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#2) [Next]
	require.Equal(t, http.StatusOK, res.StatusCode)

	list = resource.List[resource.Review]{}
	err = json.NewDecoder(res.Body).Decode(&list)
	require.NoError(t, err)
	require.Len(t, list.Data, 1)
	require.Equal(t, rvwA.ID, list.Data[0].ID)

	// Act (#3) [Previous]
	req, err = http.NewRequest("GET", fmt.Sprintf("%s/v1/places/%d/reviews", srv.URL, plc.ID), nil)
	require.NoError(t, err)
	req.URL.RawQuery = url.Values{
		"limit":         {"1"},
		"first_id":      {fmt.Sprint(list.Data[0].ID)},
		"first_created": {fmt.Sprint(list.Data[0].Created)},
	}.Encode()

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#3) [Previous]
	require.Equal(t, http.StatusOK, res.StatusCode)

	list = resource.List[resource.Review]{}
	err = json.NewDecoder(res.Body).Decode(&list)
	require.NoError(t, err)
	require.Len(t, list.Data, 1)
	require.Equal(t, rvwB.ID, list.Data[0].ID)
}
