package geo_test

import (
	"testing"

	"codeberg.org/socialmaps/api/internal/floats"
	"codeberg.org/socialmaps/api/internal/geo"
)

const tolerance = 0.000001

func almostEqual(a, b geo.BBox) bool {
	return floats.AlmostEqual(a.South, b.South, tolerance) &&
		floats.AlmostEqual(a.West, b.West, tolerance) &&
		floats.AlmostEqual(a.North, b.North, tolerance) &&
		floats.AlmostEqual(a.East, b.East, tolerance)
}

func TestBasic(t *testing.T) {
	testCases := []struct {
		lat, lon, radius float64
		exp              geo.BBox
	}{
		{
			40.7580, -73.9855, 10.0,
			geo.BBox{South: 40.757910, West: -73.985619, North: 40.758090, East: -73.985381},
		},
		{
			0.0, -78.5000, 10.0,
			geo.BBox{South: -0.000090, West: -78.500090, North: 0.000090, East: -78.499910},
		},
	}

	for i, tc := range testCases {
		act := geo.NewBBox(tc.lat, tc.lon, tc.radius)
		if !almostEqual(tc.exp, act) {
			t.Errorf(
				"Test[%d]: NewBBox(%f, %f, %f)\nreturned: {S: %f, W: %f, N: %f, E: %f}\nexpected: {S: %f, W: %f, N: %f, E: %f}",
				i, tc.lat, tc.lon, tc.radius,
				act.South, act.West, act.North, act.East,
				tc.exp.South, tc.exp.West, tc.exp.North, tc.exp.East,
			)
		}
	}
}
