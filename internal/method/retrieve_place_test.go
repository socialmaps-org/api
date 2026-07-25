package method

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"golang.socialmaps.org/api/internal/database"
	"golang.socialmaps.org/api/internal/j"
	"golang.socialmaps.org/api/internal/model"
	"golang.socialmaps.org/api/internal/must"
	"golang.socialmaps.org/api/internal/mytime"
)

func TestRetrieveExisting(t *testing.T) {
	// Arrange
	ctx := t.Context()

	qs := model.New(database.OpenInTest(t))
	plcM := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs)

	// Act
	res, err := http.Get(fmt.Sprintf("%s/v1/places/%d", srv.URL, plcM.ID))
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

func TestRetrieveMissing(t *testing.T) {
	// Arrange
	qs := model.New(database.OpenInTest(t))

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, qs)

	// Act
	res, err := http.Get(srv.URL + "/v1/places/42")
	require.NoError(t, err)

	// Assert
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}
