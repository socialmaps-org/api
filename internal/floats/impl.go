package floats

import "math"

func AlmostEqual(a, b, t float64) bool {
	// shortcut, handles infinities
	if a == b {
		return true
	}

	return math.Abs(a-b) <= t
}
