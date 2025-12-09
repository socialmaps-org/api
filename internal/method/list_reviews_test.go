package method

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/require"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/mytime"
	"codeberg.org/socialmaps/api/internal/resource"
)

func TestListReviews(t *testing.T) {
	// Arrange
	ctx := t.Context()

	mockClock := clock.NewMock()
	mytime.SetClock(mockClock)

	db := database.Open(":memory:")
	usrA := model.UpsertUser(ctx, db, "1", "Alice")
	usrB := model.UpsertUser(ctx, db, "2", "Bob")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	rvwA := model.CreateReview(ctx, db, plc.ID, usrA.ID, true, "I like it!")
	model.CreateReviewDecision(ctx, db, rvwA.ID, "test-mod", true, "")
	mockClock.Add(24 * time.Hour)
	rvwB := model.CreateReview(ctx, db, plc.ID, usrB.ID, false, "I don't like it!")
	model.CreateReviewDecision(ctx, db, rvwB.ID, "test-mod", true, "")

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, db, "")

	// Act (#1)
	req, err := http.NewRequest("GET", srv.URL+"/v1/places/"+plc.ID+"/reviews", nil)
	require.NoError(t, err)
	req.URL.RawQuery = url.Values{
		"limit": {"1"},
	}.Encode()

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#1)
	require.Equal(t, http.StatusOK, res.StatusCode)

	list := resource.List[*model.Review]{}
	err = json.NewDecoder(res.Body).Decode(&list)
	require.NoError(t, err)
	require.Len(t, list.Data, 1)
	require.Equal(t, rvwB.ID, list.Data[0].ID)

	// Act (#2) [Next]
	req, err = http.NewRequest("GET", srv.URL+"/v1/places/"+plc.ID+"/reviews", nil)
	require.NoError(t, err)
	req.URL.RawQuery = url.Values{
		"limit":          {"1"},
		"starting_after": {*list.StartingAfter},
	}.Encode()

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#2) [Next]
	require.Equal(t, http.StatusOK, res.StatusCode)

	list = resource.List[*model.Review]{}
	err = json.NewDecoder(res.Body).Decode(&list)
	require.NoError(t, err)
	require.Len(t, list.Data, 1)
	require.Equal(t, rvwA.ID, list.Data[0].ID)

	// Act (#3) [Previous]
	req, err = http.NewRequest("GET", srv.URL+"/v1/places/"+plc.ID+"/reviews", nil)
	require.NoError(t, err)
	req.URL.RawQuery = url.Values{
		"limit":         {"1"},
		"ending_before": {*list.EndingBefore},
	}.Encode()

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert (#3) [Previous]
	require.Equal(t, http.StatusOK, res.StatusCode)

	list = resource.List[*model.Review]{}
	err = json.NewDecoder(res.Body).Decode(&list)
	require.NoError(t, err)
	require.Len(t, list.Data, 1)
	require.Equal(t, rvwB.ID, list.Data[0].ID)
}
