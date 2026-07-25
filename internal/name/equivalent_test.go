package name_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"golang.socialmaps.org/api/internal/name"
)

func TestBasic(t *testing.T) {
	tests := []struct {
		a, b string
		exp  bool
	}{
		// returns true if the strings are exactly the same
		{"bora's cafe", "bora's cafe", true},
		// trims leading whitespace
		{"bora's cafe", " bora's cafe", true},
		// trims trailing whitespace
		{"bora's cafe", "bora's cafe ", true},
		// removes accents and diacritics
		{"bora's cafe", "bora's café", true},
		// removes non-alphanumeric characters
		{"bora's cafe", "bora’s-cafe!", true},
		// normalises ligatures
		{"fin cafe", "ﬁn cafe", true},
		// ignores case
		{"bora's cafe", "BORA'S CAFE", true},
		// returns false if the strings are different
		{"bora's cafe", "cafe", false},
	}

	for _, d := range tests {
		act := name.Equivalent(d.a, d.b)
		require.Equal(t, d.exp, act)

		act = name.Equivalent(d.b, d.a)
		require.Equal(t, d.exp, act)
	}
}
