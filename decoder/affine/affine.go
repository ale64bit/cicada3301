package affine

import (
	"fmt"

	"cicada/decoder"
	"cicada/decoder/mono"
	"cicada/gematria"
	"cicada/mathutil"
)

func New(a, b int, skips []int) decoder.Decoder {
	m := gematria.Len()
	if !mathutil.Coprime(a, gematria.Len()) {
		panic("a must be coprime with the alphabet len")
	}
	aInv := mathutil.ModPow(a, m-2, m)
	name := fmt.Sprintf("affine(%d, %d, %v)", a, b, skips)
	updateFn := func(_ int, _ int, r rune) string {
		x := (aInv * (m + gematria.IndexOfRune(r) - b)) % m
		return gematria.LetterOfIndex(x)
	}
	return mono.New(name, updateFn, skips)
}
