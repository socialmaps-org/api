package method

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/api/interpreter"; got != want {
			t.Errorf("path: %#v, want: %#v", got, want)
		}

		if got, want := r.FormValue("data"), `[out:json];nwr(51.8951601, -8.4717157, 51.8953399, -8.4714243)[name];out center tags;`; got != want {
			t.Errorf("body: %#v, want: %#v", got, want)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(overpassDoc))
	}))
	defer server.Close()

	db := database.Open(":memory:")
	method := LookupPlace{
		Common: Common{
			DB: db,
		},
		OverpassEndpoint: server.URL + "/api/interpreter",
	}

	// Act
	res := method.Execute(t.Context(), &lookupPlaceArgs{
		Name: "izz cafe",
		Lat:  51.89525,
		Lon:  -8.47157,
	})

	// Assert
	if got, want := res.StatusCode, 200; got != want {
		t.Fatalf("statusCode: %d, want: %d", got, 200)
	}

	plc, ok := res.Body.(*resource.Place)
	if !ok {
		t.Fatalf("returned resource is not a Place")
	}

	if got, want := plc.Name, "Izz Cafe"; got != want {
		t.Errorf("name: %s, want: %s", got, want)
	}

	if got, want := plc.Location.Lat, 51.8952597; got != want {
		// OpenStreetMap requires 7 decimal places for geographic coordinates.
		t.Errorf("lat: %.7f, want: %.7f", got, want)
	}

	if got, want := plc.Location.Lon, -8.4715779; got != want {
		// OpenStreetMap requires 7 decimal places for geographic coordinates.
		t.Errorf("lon: %.7f, want: %.7f", got, want)
	}
}

func TestLookupExisting(t *testing.T) {
	// Arrange
	db := database.Open(":memory:")
	model.CreatePlace(t.Context(), db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	method := LookupPlace{
		Common: Common{
			DB: db,
		},
		OverpassEndpoint: "",
	}

	// Act
	res := method.Execute(t.Context(), &lookupPlaceArgs{
		Name: "izz cafe",
		Lat:  51.89525,
		Lon:  -8.47157,
	})

	// Assert
	if got, want := res.StatusCode, 200; got != want {
		t.Fatalf("statusCode: %d, want: %d", got, 200)
	}

	plc, ok := res.Body.(*resource.Place)
	if !ok {
		t.Fatalf("returned resource is not a Place")
	}

	if got, want := plc.Name, "Izz Cafe"; got != want {
		t.Errorf("name: %s, want: %s", got, want)
	}

	if got, want := plc.Location.Lat, 51.8952597; got != want {
		// OpenStreetMap requires 7 decimal places for geographic coordinates.
		t.Errorf("lat: %.7f, want: %.7f", got, want)
	}

	if got, want := plc.Location.Lon, -8.4715779; got != want {
		// OpenStreetMap requires 7 decimal places for geographic coordinates.
		t.Errorf("lon: %.7f, want: %.7f", got, want)
	}
}
