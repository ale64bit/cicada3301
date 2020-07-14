package stringutil

import "strings"

func Reverse(s string) string {
	var b strings.Builder
	for i := len(s) - 1; i >= 0; i-- {
		b.WriteRune(s[i])
	}
	return b.String()
}
