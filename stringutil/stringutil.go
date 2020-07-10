package stringutil

func Reverse(s string) string {
	ret := ""
	for _, r := range s {
		ret = string(r) + ret
	}
	return ret
}
