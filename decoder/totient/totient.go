package totient

import (
	"fmt"

	"cicada/decoder"
	"cicada/decoder/mono"
	"cicada/gematria"
	"cicada/mathutil"
)

func New(offset int, skips []int) decoder.Decoder {
	name := fmt.Sprintf("totient(%d)", offset)
	updateFn := func(_ int, codedRuneIndex int, r rune) string {
		m := gematria.Len()
		x := (m + gematria.IndexOfRune(r) - (mathutil.PrimeAt(codedRuneIndex)%m - 1) + offset) % m
		return gematria.LetterOfIndex(x)
	}
	return mono.New(name, updateFn, skips)
}
