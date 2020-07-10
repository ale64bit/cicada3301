package vigenere

import (
	"fmt"

	"cicada/decoder"
	"cicada/decoder/mono"
	"cicada/gematria"
)

type vigenere struct {
	keyword []rune
	skips   []int
}

func New(keyword string, skips []int) decoder.Decoder {
	name := fmt.Sprintf("vigenere(%q)", keyword)
	ks := []rune(keyword)
	updateFn := func(_ int, codedRuneIndex int, r rune) string {
		m := gematria.Len()
		idx := codedRuneIndex % len(ks)
		return gematria.LetterOfIndex((gematria.IndexOfRune(r) + m - gematria.IndexOfRune(ks[idx])) % m)
	}
	return mono.New(name, updateFn, skips)
}
