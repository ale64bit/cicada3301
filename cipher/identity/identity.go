package identity

import (
	"cicada/cipher"
	"cicada/gematria"
)

type identity struct{}

func New() cipher.Cipher {
	return identity{}
}

func (identity) ID() string { return "identity" }

func (identity) Encode(s string) string { return s }

func (identity) Decode(s string) string {
	return gematria.MapRunes(s, func(r rune) string { return string(r) })
}
