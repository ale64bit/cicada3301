package cat

import (
	"fmt"
	"strings"

	"cicada/cipher"
)

type cat struct {
	ciphers []cipher.Cipher
}

func New(ciphers ...cipher.Cipher) cipher.Cipher {
	return &cat{ciphers}
}

func (c *cat) ID() string {
	var ids []string
	for _, ci := range c.ciphers {
		ids = append(ids, ci.ID())
	}
	return fmt.Sprintf("cat(%s)", strings.Join(ids, ", "))
}

func (c *cat) Encode(s string) string {
	for _, ci := range c.ciphers {
		s = ci.Encode(s)
	}
	return s
}

func (c *cat) Decode(s string) string {
	for i := len(c.ciphers) - 1; i >= 0; i-- {
		s = c.ciphers[i].Decode(s)
	}
	return s
}
