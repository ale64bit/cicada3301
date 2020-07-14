package caesar

import (
	"cicada/cipher"
	"cicada/cipher/affine"
)

func New(b int) cipher.Cipher {
	return affine.New(1, b)
}
