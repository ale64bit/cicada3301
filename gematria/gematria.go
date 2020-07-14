package gematria

import "strings"

type graphInfo struct {
	letter string
	index  int
	value  int
}

type encoding struct {
	from string
	to   string
}

var (
	runeMap = map[rune]graphInfo{
		'ᚠ': {"F", 0, 2},
		'ᚢ': {"U", 1, 3},
		'ᚦ': {"TH", 2, 5},
		'ᚩ': {"O", 3, 7},
		'ᚱ': {"R", 4, 11},
		'ᚳ': {"C", 5, 13},
		'ᚷ': {"G", 6, 17},
		'ᚹ': {"W", 7, 19},
		'ᚻ': {"H", 8, 23},
		'ᚾ': {"N", 9, 29},
		'ᛁ': {"I", 10, 31},
		'ᛄ': {"J", 11, 37},
		'ᛇ': {"EO", 12, 41},
		'ᛈ': {"P", 13, 43},
		'ᛉ': {"X", 14, 47},
		'ᛋ': {"S", 15, 53},
		'ᛏ': {"T", 16, 59},
		'ᛒ': {"B", 17, 61},
		'ᛖ': {"E", 18, 67},
		'ᛗ': {"M", 19, 71},
		'ᛚ': {"L", 20, 73},
		'ᛝ': {"ING", 21, 79},
		'ᛟ': {"OE", 22, 83},
		'ᛞ': {"D", 23, 89},
		'ᚪ': {"A", 24, 97},
		'ᚫ': {"AE", 25, 101},
		'ᚣ': {"Y", 26, 103},
		'ᛡ': {"IO", 27, 107},
		'ᛠ': {"EA", 28, 109},
	}

	punctuationMap = map[rune]graphInfo{
		'-': {" ", 0, 1},
		',': {", ", 1, 2},
		'<': {"<", 2, 3},
		'>': {">", 3, 3},
		'.': {".", 4, 4},
		'~': {"~", 5, 5},
		'*': {"*", 6, 10},
		'!': {"!", 7, 11},
		'+': {"+", 8, 13},
		'#': {"#", 9, 23},
	}

	controlMap = map[rune]string{
		'[':  "\033[31m",
		']':  "\033[39m",
		'\'': "'",
		'"':  `"`,
		'0':  "0",
		'1':  "1",
		'2':  "2",
		'3':  "3",
		'4':  "4",
		'5':  "5",
		'6':  "6",
		'7':  "7",
		'8':  "8",
		'9':  "9",
	}

	encodings = []encoding{
		// trigraphs
		encoding{"ING", "ᛝ"},
		// digraphs
		encoding{"AE", "ᚫ"},
		encoding{"EA", "ᛠ"},
		encoding{"EO", "ᛇ"},
		encoding{"IO", "ᛡ"},
		encoding{"OE", "ᛟ"},
		encoding{"NG", "ᛝ"},
		encoding{"TH", "ᚦ"},
		// single
		encoding{"A", "ᚪ"},
		encoding{"B", "ᛒ"},
		encoding{"C", "ᚳ"},
		encoding{"D", "ᛞ"},
		encoding{"E", "ᛖ"},
		encoding{"F", "ᚠ"},
		encoding{"G", "ᚷ"},
		encoding{"H", "ᚻ"},
		encoding{"I", "ᛁ"},
		encoding{"J", "ᛄ"},
		encoding{"K", "ᚳ"},
		encoding{"L", "ᛚ"},
		encoding{"M", "ᛗ"},
		encoding{"N", "ᚾ"},
		encoding{"O", "ᚩ"},
		encoding{"P", "ᛈ"},
		encoding{"Q", "ᚳᚹ"},
		encoding{"R", "ᚱ"},
		encoding{"S", "ᛋ"},
		encoding{"T", "ᛏ"},
		encoding{"U", "ᚢ"},
		encoding{"V", "ᚢ"},
		encoding{"W", "ᚹ"},
		encoding{"X", "ᛉ"},
		encoding{"Y", "ᚣ"},
		// punctuation
		encoding{" ", "-"},
		encoding{",", ","},
		encoding{"<", "<"},
		encoding{">", ">"},
		encoding{".", "."},
		encoding{"~", "~"},
		encoding{"*", "*"},
		encoding{"!", "!"},
		encoding{"+", "+"},
		encoding{"#", "#"},
	}
)

func Len() int {
	return len(runeMap)
}

func IsValidRune(r rune) bool {
	_, ok := runeMap[r]
	return ok
}

func IsValidPunctuation(r rune) bool {
	_, ok := punctuationMap[r]
	return ok
}

func IsValid(r rune) bool {
	return IsValidRune(r) || IsValidPunctuation(r)
}

func LetterOfRune(r rune) string {
	info, ok := runeMap[r]
	if !ok {
		return "?"
	}
	return info.letter
}

func ValueOfRune(r rune) int {
	switch {
	case IsValidRune(r):
		return runeMap[r].value
	case IsValidPunctuation(r):
		return punctuationMap[r].value
	default:
		return -1
	}
}

func IndexOfRune(r rune) int {
	info, ok := runeMap[r]
	if !ok {
		return -1
	}
	return info.index
}

func LetterOfIndex(i int) string {
	for _, v := range runeMap {
		if v.index == i {
			return v.letter
		}
	}
	return "?"
}

func RuneOfLetter(s string) rune {
	for k, v := range runeMap {
		if v.letter == s {
			return k
		}
	}
	return '?'
}

func RuneOfIndex(i int) rune {
	for k, v := range runeMap {
		if v.index == i {
			return k
		}
	}
	return '?'
}

func RuneCount(s string) int {
	ret := 0
	for _, r := range s {
		if IsValidRune(r) {
			ret++
		}
	}
	return ret
}

func FilterRunes(s string) string {
	var b strings.Builder
	for _, r := range s {
		if IsValidRune(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func WordStream(s string) []string {
	var ret []string
	var cur strings.Builder
	for _, r := range s {
		if IsValidRune(r) {
			cur.WriteRune(r)
		} else if IsValidPunctuation(r) {
			if cur.Len() > 0 {
				ret = append(ret, cur.String())
				cur.Reset()
			}
		}
	}
	if cur.Len() > 0 {
		ret = append(ret, cur.String())
	}
	return ret
}

func RuneStream(s string) []rune {
	var ret []rune
	for _, r := range s {
		if IsValidRune(r) {
			ret = append(ret, r)
		}
	}
	return ret
}

func ValueStream(s string) []int {
	var ret []int
	for _, r := range s {
		if IsValidRune(r) {
			ret = append(ret, ValueOfRune(r))
		}
	}
	return ret
}

func IndexStream(s string) []int {
	var ret []int
	for _, r := range s {
		if IsValidRune(r) {
			ret = append(ret, IndexOfRune(r))
		}
	}
	return ret
}

func Decode(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case IsValidRune(r):
			b.WriteString(runeMap[r].letter)
		case IsValidPunctuation(r):
			p := punctuationMap[r]
			b.WriteString(p.letter)
			if p.value >= 4 {
				b.WriteString("\n")
			}
		case r == 'F':
			b.WriteRune('F')
		default:
			b.WriteString(controlMap[r])
		}
	}
	return b.String()
}

func EncodeLiteral(s string) string {
	var b strings.Builder
	for _, r := range s {
		b.WriteRune(RuneOfLetter(string(r)))
	}
	return b.String()
}

func Encode(s string) string {
	for _, e := range encodings {
		s = strings.ReplaceAll(s, e.from, e.to)
	}
	return s
}

type MapFunc = func(rune) string

func MapRunes(s string, f MapFunc) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case IsValidRune(r):
			b.WriteString(f(r))
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func MapRunesAndPunctuation(s string, f MapFunc) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case IsValid(r):
			b.WriteString(f(r))
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
