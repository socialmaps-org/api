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
	require.Equal(t, "xyz_1AAAAAADIIPZ2M", id1.String()[:18])
	require.Equal(t, "xyz_1AAAAAADIIPZ2M", id2.String()[:18])

	require.NotEqual(t, id1.String()[18:], id2.String()[18:])
}

func TestParseID(t *testing.T) {
	// Arrange
	const idS = "abc_1AAAAAADJCL4WUrmzGj2YeBpd60CDUtDq1tW"

	// Act
	id := ParseID(idS)

	// Assert
	require.Equal(t, "abc", id.kind)
	require.Equal(t, "1", id.version)
	require.Equal(t, time.Date(2025, 11, 11, 8, 52, 58, 0, time.UTC), id.time)
	require.Len(t, id.rand, 22)
	require.Equal(t, "rmzGj2YeBpd60CDUtDq1tW", id.rand)
}
