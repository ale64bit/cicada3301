package mono

import (
	"strings"

	"cicada/decoder"
	"cicada/gematria"
)

type mono struct {
	name     string
	updateFn func(int, int, rune) string
	skips    []int
}

func New(name string, updateFn func(int, int, rune) string, skips []int) decoder.Decoder {
	return &mono{name, updateFn, skips}
}

func (d *mono) Name() string { return d.name }

func (d *mono) DecodeStream(msg string) (string, []string) {
	shouldSkip := map[int]bool{}
	for _, x := range d.skips {
		shouldSkip[x] = true
	}
	wordStream := ""
	ret := ""
	runeIndex := 0
	codedRuneIndex := 0
	for _, r := range msg {
		s, decoded, wordBreak := decoder.MapNonrune(r)
		if !decoded {
			if gematria.IsValidRune(r) {
				if !shouldSkip[runeIndex] {
					s = d.updateFn(runeIndex, codedRuneIndex, r)
					codedRuneIndex++
				} else {
					s = gematria.LetterOfRune(r)
				}
				runeIndex++
			}
			wordStream += s
		} else if wordBreak {
			wordStream += " "
		}
		ret += s
	}
	return ret, strings.Fields(wordStream)
}

func (d *mono) Decode(msg string) string {
	s, _ := d.DecodeStream(msg)
	return s
}
