package method

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/resource"
)

func TestRetrieveExisting(t *testing.T) {
	// Arrange
	ctx := t.Context()

	db := database.Open(":memory:")
	plcM := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, db, "")

	// Act
	res, err := http.Get(srv.URL + "/v1/places/" + plcM.ID)
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

func TestRetrieveMissing(t *testing.T) {
	// Arrange
	db := database.Open(":memory:")

	authr := NewTestAuthenticator(t)
	srv := NewTestServer(t, authr, db, "")

	// Act
	res, err := http.Get(srv.URL + "/v1/places/plc_missing")
	require.NoError(t, err)

	// Assert
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}
