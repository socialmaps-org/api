package name

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/unicode/rangetable"
)

var spaces = regexp.MustCompile(`\s+`)

func Equivalent(a, b string) bool {
	if a == b {
		return true
	}

	letnum := runes.In(rangetable.Merge(unicode.Letter, unicode.Number))

	t := transform.Chain(
		norm.NFKD,
		runes.Map(func(r rune) rune {
			if letnum.Contains(r) {
				return r
			} else {
				return ' '
			}
		}),
	)

	a, _, err := transform.String(t, a)
	if err != nil {
		panic(err)
	}
	b, _, err = transform.String(t, b)
	if err != nil {
		panic(err)
	}

	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)

	a = spaces.ReplaceAllLiteralString(a, " ")
	b = spaces.ReplaceAllLiteralString(b, " ")

	return strings.EqualFold(a, b)
}
