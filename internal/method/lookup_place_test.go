package method

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/must"
	"codeberg.org/socialmaps/api/internal/resource"
)

const nominatimDoc = `
[
  {
    "place_id": 416921008,
    "licence": "Data © OpenStreetMap contributors, ODbL 1.0. http://osm.org/copyright",
    "osm_type": "node",
    "osm_id": 7095470096,
    "lat": "51.8952597",
    "lon": "-8.4715779",
    "category": "amenity",
    "type": "restaurant",
    "place_rank": 30,
    "importance": 6.924620431825769e-5,
    "addresstype": "amenity",
    "name": "Izz Cafe",
    "display_name": "Izz Cafe, 14, George's Quay, South Parish, South Gate A, Cork, County Cork, Munster, T12 EY24, Ireland",
    "namedetails": { "name": "Izz Cafe" },
    "boundingbox": ["51.8952097", "51.8953097", "-8.4716279", "-8.4715279"]
  }
]
`

func TestLookupNew(t *testing.T) {
	// Arrange
	ctx := t.Context()

	qs := model.New(database.Open(":memory:"))

	nominatimSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/search", r.URL.Path)
		require.Equal(t, "izz cafe", r.FormValue("amenity"))
		require.Equal(t, "jsonv2", r.FormValue("format"))
		require.Equal(t, "1", r.FormValue("namedetails"))
		require.Equal(t, "-8.4722987, 51.8948003, -8.4708413, 51.8956997", r.FormValue("viewbox"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(nominatimDoc))
	}))
	defer nominatimSrv.Close()

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs, nominatimSrv.URL)

	// Act
	req, err := http.NewRequest("GET", srv.URL+"/v1/places/lookup", nil)
	require.NoError(t, err)
	req.URL.RawQuery = url.Values{
		"name": {"izz cafe"},
		"lat":  {"51.89525"},
		"lon":  {"-8.47157"},
	}.Encode()

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, http.StatusOK, res.StatusCode)

	var plcR resource.Place
	err = json.NewDecoder(res.Body).Decode(&plcR)
	require.NoError(t, err)
	require.Equal(t, "Izz Cafe", plcR.Name)
	require.Equal(t, 51.8952597, plcR.Location.Lat)
	require.Equal(t, -8.4715779, plcR.Location.Lon)

	plcM := must.Get(qs.LoadPlace(ctx, plcR.ID))
	require.NotNil(t, plcM)
	require.Equal(t, "Izz Cafe", plcM.Name)
	require.Equal(t, 51.8952597, plcM.Lat)
	require.Equal(t, -8.4715779, plcM.Lon)
}

func TestLookupExisting(t *testing.T) {
	// Arrange
	ctx := t.Context()

	qs := model.New(database.Open(":memory:"))

	must.Get(qs.CreatePlace(ctx, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs, "")

	// Act
	req, err := http.NewRequest("GET", srv.URL+"/v1/places/lookup", nil)
	require.NoError(t, err)
	req.URL.RawQuery = url.Values{
		"name": {"izz cafe"},
		"lat":  {"51.89525"},
		"lon":  {"-8.47157"},
	}.Encode()

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// Assert
	require.Equal(t, http.StatusOK, res.StatusCode)

	var plcR resource.Place
	err = json.NewDecoder(res.Body).Decode(&plcR)
	require.NoError(t, err)
	require.Equal(t, "Izz Cafe", plcR.Name)
	require.Equal(t, 51.8952597, plcR.Location.Lat)
	require.Equal(t, -8.4715779, plcR.Location.Lon)
}
