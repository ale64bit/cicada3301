package experimental

import (
	"fmt"

	"cicada/cipher"
	"cicada/gematria"
)

type experimental struct {
	id string
	f  func(int, rune, int) (int, int)
}

func New(id string, f func(int, rune, int) (int, int)) cipher.Cipher {
	return &experimental{id, f}
}

func NewStateless(id string, f func(int, rune) int) cipher.Cipher {
	ff := func(i int, r rune, _ int) (int, int) {
		return f(i, r), 0
	}
	return New(id, ff)
}

func (c *experimental) ID() string {
	return fmt.Sprintf("experimental(%s)", c.id)
}

func (c *experimental) Encode(msg string) string {
	panic("TODO: not implemented")
}

func (c *experimental) Decode(msg string) string {
	index, state := 0, 0
	f := func(r rune) string {
		defer func() { index++ }()
		x, nextState := c.f(index, r, state)
		state = nextState
		return string(gematria.RuneOfIndex(x))
	}
	return gematria.MapRunes(msg, f)
}
