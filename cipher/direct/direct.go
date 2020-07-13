package direct

import (
	"cicada/cipher"
	"cicada/gematria"
)

type direct struct{}

func New() cipher.Cipher {
	return direct{}
}

func (direct) ID() string             { return "direct" }
func (direct) Encode(s string) string { return gematria.Encode(s) }
func (direct) Decode(s string) string { return gematria.MapRunes(s, gematria.LetterOfRune) }
