package gematria

type runeInfo struct {
	letter string
	index  int
	value  int
}

var (
	runeMap = map[rune]runeInfo{
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
)

func Len() int {
	return len(runeMap)
}

func IsValidRune(r rune) bool {
	_, ok := runeMap[r]
	return ok
}

func LetterOfRune(r rune) string {
	info, ok := runeMap[r]
	if !ok {
		return "?"
	}
	return info.letter
}

func ValueOfRune(r rune) int {
	info, ok := runeMap[r]
	if !ok {
		return -1
	}
	return info.value
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

func RuneCount(s string) int {
	ret := 0
	for _, r := range s {
		if _, ok := runeMap[r]; ok {
			ret++
		}
	}
	return ret
}

func Encode(s string) string {
	ret := ""
	for _, r := range s {
		ret += string(RuneOfLetter(string(r)))
	}
	return ret
}
