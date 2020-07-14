package fskip

import (
	"fmt"
	"sort"
	"strings"

	"cicada/cipher"
	"cicada/gematria"
)

type fskip struct {
	skips []int
}

func New(skips ...int) cipher.Cipher {
	sort.Ints(skips)
	return &fskip{skips}
}

func (c *fskip) ID() string { return fmt.Sprintf("fskip(%v)", c.skips) }

func (c *fskip) Encode(string) string {
	panic("TODO: not implemented")
}

func (c *fskip) Decode(s string) string {
	var b strings.Builder
	index := 0
	cur := 0
	for _, r := range s {
		switch {
		case gematria.IsValidRune(r):
			if r == 'ᚠ' && cur < len(c.skips) && c.skips[cur] == index {
				b.WriteRune('F')
				cur++
			} else {
				b.WriteRune(r)
			}
			index++
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
