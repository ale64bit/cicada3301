package atbash

import (
	"cicada/cipher"
	"cicada/cipher/affine"
	"cicada/gematria"
)

func New() cipher.Cipher {
	m := gematria.Len()
	return affine.New(m-1, m-1)
}
