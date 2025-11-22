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
	"codeberg.org/socialmaps/api/internal/resource"
)

const overpassDoc = `
{
  "version": 0.6,
  "generator": "Overpass API 0.7.62.8 e802775f",
  "osm3s": {
    "timestamp_osm_base": "2025-10-18T16:51:30Z",
    "copyright": "The data included in this document is from www.openstreetmap.org. The data is made available under ODbL."
  },
  "elements": [
    {
      "type": "node",
      "id": 7095470096,
      "lat": 51.8952597,
      "lon": -8.4715779,
      "tags": {
        "addr:city": "Cork",
        "addr:housenumber": "14",
        "addr:postcode": "T12 EY24",
        "addr:street": "George's Quay",
        "amenity": "restaurant",
        "cuisine": "palestinian",
        "entrance": "main",
        "name": "Izz Cafe",
        "note": "Called \"cafe\" but self-identifies and functions as a restaurant",
        "opening_hours": "We 12:00-20:00; Th-Sa 12:00 21:00; Su 12:00-18:00",
        "phone": "+353 21 229 0689",
        "takeaway": "yes",
        "website": "https://izz.ie/",
        "wheelchair": "yes"
      }
    }
  ]
}
`

func TestLookupNew(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")

	overpassSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/interpreter", r.URL.Path)
		require.Equal(
			t,
			`[out:json];nwr(51.8951601, -8.4717157, 51.8953399, -8.4714243)[name];out center tags;`,
			r.FormValue("data"),
		)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(overpassDoc))
	}))
	defer overpassSrv.Close()

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, db, overpassSrv.URL+"/api/interpreter")

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

	plcM := model.LoadPlaceByID(ctx, db, plcR.ID)
	require.NotNil(t, plcM)
	require.Equal(t, "Izz Cafe", plcM.Name)
	require.Equal(t, 51.8952597, plcM.Lat)
	require.Equal(t, -8.4715779, plcM.Lon)
}

func TestLookupExisting(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")

	model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, db, "")

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
