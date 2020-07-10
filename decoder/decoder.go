package decoder

import (
	"math"

	"cicada/dict"
	"cicada/gematria"
)

type Decoder interface {
	Name() string
	Decode(string) string
	DecodeStream(string) (string, []string)
}

func MapNonrune(r rune) (string, bool, bool) {
	switch r {
	case '-':
		return " ", true, true
	case '.':
		return "::", true, true
	case ' ':
		return " ", true, false
	case '\n':
		return "\n", true, false
	case 'R':
		return "\033[31m", true, false
	case 'B':
		return "\033[39m", true, false
	default:
		if gematria.IndexOfRune(r) == -1 {
			return string(r), true, false
		} else {
			return "", false, false
		}
	}
}

func Score(words []string) float64 {
	if len(words) == 0 {
		return 0.0
	}
	total := 0
	score := 0.0
	for _, word := range words {
		total += len(word)
		if dict.Exists(word) {
			score += math.Pow(float64(len(word)), 1.0)
		}
	}
	return score / float64(total)
}
