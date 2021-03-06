package totient

import (
	"cicada/cipher"
	"cicada/gematria"
	"cicada/mathutil"
)

type totient struct{}

func New(skips []int) cipher.Cipher {
	return &totient{}
}

func (c *totient) ID() string { return "totient()" }

func (c *totient) Encode(string) string {
	panic("TODO: not implemented")
}

func (c *totient) Decode(s string) string {
	m := gematria.Len()
	index := 0
	f := func(r rune) string {
		defer func() { index++ }()
		y := gematria.IndexOfRune(r)
		x := (y - (mathutil.PrimeAt(index)%m - 1) + m) % m
		return string(gematria.RuneOfIndex(x))
	}
	return gematria.MapRunes(s, f)
}
