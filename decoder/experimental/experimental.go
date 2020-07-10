package experimental

import (
	"fmt"

	"cicada/decoder"
	"cicada/decoder/mono"
	"cicada/gematria"
)

func New(f func(i int, r rune) int, skips []int) decoder.Decoder {
	name := fmt.Sprintf("experimental")
	updateFn := func(_ int, codedRuneIndex int, r rune) string {
		x := f(codedRuneIndex, r)
		return gematria.LetterOfIndex(x)
	}
	return mono.New(name, updateFn, skips)
}
