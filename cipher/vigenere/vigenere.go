package vigenere

import (
	"fmt"

	"cicada/cipher"
	"cicada/gematria"
)

type vigenere struct {
	keys  []string
	skips []int
}

func New(keywords []string, skips []int) cipher.Cipher {
	return &vigenere{keywords, skips}
}

func (c *vigenere) ID() string {
	return fmt.Sprintf("vigenere(keys=%v, skips=%v)", c.keys, c.skips)
}

func (c *vigenere) Encode(string) string {
	panic("TODO: not implemented")
}

func (c *vigenere) Decode(s string) string {
	m := gematria.Len()
	index := 0
	f := func(r rune) string {
		defer func() { index++ }()
		x := gematria.IndexOfRune(r)
		for _, key := range c.keys {
			runes := []rune(key)
			pos := index % len(runes)
			x -= gematria.IndexOfRune(runes[pos])
			for x < 0 {
				x += m
			}
			x %= m
		}
		return gematria.LetterOfIndex(x)
	}
	return cipher.MapRunesWithSkips(s, f, c.skips)
}
