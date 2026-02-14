package multiline

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDedent(t *testing.T) {
	// Arrange
	ms := `
		This is a multiline string.

		    preprocessed code block

		* Unordered list \
		  with line breaks
	`

	// Act
	act := Dedent(ms)

	// Assert
	exp := `This is a multiline string.

    preprocessed code block

* Unordered list \
  with line breaks
`
	require.Equal(t, exp, act)
}
