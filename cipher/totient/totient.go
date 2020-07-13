package totient

import (
	"fmt"

	"cicada/cipher"
	"cicada/gematria"
	"cicada/mathutil"
)

type totient struct {
	skips []int
}

func New(skips []int) cipher.Cipher {
	return &totient{skips}
}

func (c *totient) ID() string {
	return fmt.Sprintf("totient(skips=%v)", c.skips)
}

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
		return gematria.LetterOfIndex(x)
	}
	return cipher.MapRunesWithSkips(s, f, c.skips)
}
