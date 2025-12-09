package moderation

import (
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func CleanUp(txt string) string {
	t := transform.Chain(
		norm.NFKC,
		runes.Remove(runes.Predicate(func(r rune) bool {
			isValid := unicode.IsLetter(r) ||
				unicode.IsNumber(r) ||
				unicode.IsSpace(r) ||
				unicode.IsPunct(r) ||
				unicode.IsGraphic(r)
			return !isValid
		})),
	)

	// Execute
	clean, _, _ := transform.String(t, txt)

	return clean
}
