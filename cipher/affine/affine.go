package affine

import (
	"fmt"

	"cicada/cipher"
	"cicada/gematria"
	"cicada/mathutil"
)

type affine struct {
	a    int
	b    int
	aInv int
}

func New(a, b int) cipher.Cipher {
	m := gematria.Len()
	if !mathutil.Coprime(a, gematria.Len()) {
		panic("a must be coprime with the alphabet len")
	}
	aInv := mathutil.ModPow(a, m-2, m)
	return &affine{a, b, aInv}
}

func (c *affine) ID() string {
	return fmt.Sprintf("affine(a=%d, b=%d)", c.a, c.b)
}

func (c *affine) Encode(string) string {
	panic("TODO: not implemented")
}

func (c *affine) Decode(s string) string {
	m := gematria.Len()
	f := func(r rune) string {
		y := gematria.IndexOfRune(r)
		x := (c.aInv * (y - c.b + m)) % m
		return gematria.LetterOfIndex(x)
	}
	return gematria.MapRunes(s, f)
}
