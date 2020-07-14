package cipher

import (
	"cicada/gematria"
)

type Cipher interface {
	ID() string
	Encode(string) string
	Decode(string) string
}

func MapRunesWithSkips(s string, f gematria.MapFunc, skips []int) string {
	shouldSkip := map[int]bool{}
	for _, x := range skips {
		shouldSkip[x] = true
	}
	index := 0
	ff := func(r rune) string {
		defer func() { index++ }()
		if shouldSkip[index] {
			return string(r)
		}
		return f(r)
	}
	return gematria.MapRunes(s, ff)
}
