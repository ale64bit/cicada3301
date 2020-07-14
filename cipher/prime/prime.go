package prime

import (
	"cicada/cipher"
	"cicada/gematria"
	"cicada/mathutil"
)

type prime struct{}

func New() cipher.Cipher {
	return &prime{}
}

func (c *prime) ID() string { return "prime()" }

func (c *prime) Encode(string) string {
	panic("TODO: not implemented")
}

func (c *prime) Decode(s string) string {
	m := gematria.Len()
	index := 0
	f := func(r rune) string {
		defer func() { index++ }()
		y := gematria.IndexOfRune(r)
		x := (y - mathutil.PrimeAt(index)%m + m) % m
		return string(gematria.RuneOfIndex(x))
	}
	return gematria.MapRunes(s, f)
}
