package direct

import (
	"fmt"

	"cicada/decoder"
	"cicada/decoder/mono"
	"cicada/gematria"
)

func New(offset int) decoder.Decoder {
	name := fmt.Sprintf("direct(%d)", offset)
	updateFn := func(_ int, _ int, r rune) string {
		return gematria.LetterOfIndex((gematria.IndexOfRune(r) + offset) % gematria.Len())
	}
	return mono.New(name, updateFn, nil)
}
