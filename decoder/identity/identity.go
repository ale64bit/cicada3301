package identity

import (
	"cicada/decoder"
	"cicada/decoder/mono"
)

func New() decoder.Decoder {
	return mono.New("identity", update, nil)
}

func update(_ int, _ int, r rune) string {
	return string(r)
}
