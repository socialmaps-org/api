package model

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/require"

	"codeberg.org/socialmaps/api/internal/mytime"
)

func TestNewRandomID(t *testing.T) {
	// Arrange
	mockClock := clock.NewMock()
	mytime.SetClock(mockClock)
	mockClock.Set(time.Date(2025, 6, 7, 8, 9, 10, 0, time.UTC))

	// Act
	id1 := NewRandomID("xyz")
	id2 := NewRandomID("xyz")

	// Assert
	require.Equal(t, "xyz_1000000388FPQC", id1.String()[:18])
	require.Equal(t, "xyz_1000000388FPQC", id2.String()[:18])

	require.NotEqual(t, id1.String()[18:], id2.String()[18:])
}

func TestParseID(t *testing.T) {
	// Arrange
	const idS = "abc_1000000392BSMK0mq7VVQGTSUH1O1xdF96iY"

	// Act
	id := ParseID(idS)

	// Assert
	require.Equal(t, "abc", id.kind)
	require.Equal(t, "1", id.version)
	require.Equal(t, time.Date(2025, 11, 11, 8, 52, 58, 0, time.UTC), id.time)
	require.Len(t, id.rand, 22)
	require.Equal(t, "0mq7VVQGTSUH1O1xdF96iY", id.rand)
}

func TestEarliestID(t *testing.T) {
	// Act (Generate)
	idG := EarliestID("xyz")

	// Assert (Generate)
	require.Equal(t, "xyz_100000000000000000000000000000000000", idG.String())

	// Act (Parse)
	idP := ParseID(idG.String())

	// Assert (Parse)
	require.Equal(t, "xyz", idP.kind)
	require.Equal(t, "1", idP.version)
	require.Equal(t, time.Date(1970, 1, 1, 00, 00, 00, 00, time.UTC), idP.time)
	require.Len(t, idP.rand, 22)
	require.Equal(t, "0000000000000000000000", idP.rand)
}

func TestLatestID(t *testing.T) {
	// Act (Generate)
	idG := LatestID("xyz")

	// Assert (Generate)
	require.Equal(t, "xyz_100000ENVUH0NUzzzzzzzzzzzzzzzzzzzzzz", idG.String())

	// Act (Parse)
	idP := ParseID(idG.String())

	// Assert (Parse)
	require.Equal(t, "xyz", idP.kind)
	require.Equal(t, "1", idP.version)
	require.Equal(t, time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC), idP.time)
	require.Len(t, idP.rand, 22)
	require.Equal(t, "zzzzzzzzzzzzzzzzzzzzzz", idP.rand)
}
