package method

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/j"
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/must"
	"codeberg.org/socialmaps/api/internal/mytime"
)

func TestLookupNew(t *testing.T) {
	// Arrange
	ctx := t.Context()

	qs := model.New(database.OpenInTest(t))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest("GET", srv.URL+"/v1/places/lookup", nil)
	require.NoError(t, err)
	req.URL.RawQuery = url.Values{
		"name": {"Woo"},
		"lat":  {"43.733047"},
		"lon":  {"7.419294"},
	}.Encode()

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, http.StatusOK, res.StatusCode)

	var plcR any
	err = json.NewDecoder(res.Body).Decode(&plcR)
	require.NoError(t, err)
	require.Equal(t, "Woo", j.Get[string](plcR, "name"))
	require.Equal(t, 43.7330475, j.Get[float64](plcR, "location", "lat"))
	require.Equal(t, 7.4192941, j.Get[float64](plcR, "location", "lon"))

	plcM := must.Get(qs.LoadPlace(ctx, j.Get[int64](plcR, "id")))
	require.NotNil(t, plcM)
	require.Equal(t, "Woo", plcM.Name)
	require.Equal(t, 43.7330475, plcM.Lat)
	require.Equal(t, 7.4192941, plcM.Lon)
}

func TestLookupExisting(t *testing.T) {
	// Arrange
	ctx := t.Context()

	qs := model.New(database.OpenInTest(t))

	must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs)

	// Act
	req, err := http.NewRequest("GET", srv.URL+"/v1/places/lookup", nil)
	require.NoError(t, err)
	req.URL.RawQuery = url.Values{
		"name": {"Woo"},
		"lat":  {"43.733047"},
		"lon":  {"7.419294"},
	}.Encode()

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, http.StatusOK, res.StatusCode)

	var plcR any
	err = json.NewDecoder(res.Body).Decode(&plcR)
	require.NoError(t, err)
	require.Equal(t, "Woo", j.Get[string](plcR, "name"))
	require.Equal(t, 43.7330475, j.Get[float64](plcR, "location", "lat"))
	require.Equal(t, 7.4192941, j.Get[float64](plcR, "location", "lon"))
}
