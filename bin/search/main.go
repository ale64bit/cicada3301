package main

import (
	"flag"
	"fmt"
	"runtime"
	"strings"
	"sync"

	"cicada/cipher"
	"cicada/cipher/affine"
	"cicada/cipher/direct"
	"cicada/cipher/experimental"
	"cicada/cipher/vigenere"
	"cicada/data"
	"cicada/dict"
	"cicada/gematria"
	"cicada/mathutil"
)

var (
	prefixLen        = flag.Int("prefix_len", 180, "Length of the message prefix to try to decode (the longer prefix, the longer it takes)")
	selectedSections = flag.String("selected_sections", "", "Comma-separated list of the sections to try decoding (if empty, all sessions are tried)")
	matchScore       = flag.Float64("match_score", 0.5, "Minimum score to be reported as a match")
)

func buildSkipSets(msg string) [][]int {
	ret := [][]int{}
	// for k := 0; k < gematria.Len(); k++ {
	for k := 0; k < 1; k++ {
		var loc []int
		index := 0
		for _, r := range msg {
			id := gematria.IndexOfRune(r)
			if id == -1 {
				continue
			}
			if id == k {
				loc = append(loc, index)
			}
			index++
		}
		n := len(loc)
		for mask := 0; mask < (1 << n); mask++ {
			skips := []int{}
			for i := 0; i < n; i++ {
				if (mask & (1 << i)) != 0 {
					skips = append(skips, loc[i])
				}
			}
			if len(skips) > 0 {
				ret = append(ret, skips)
			}
		}
	}
	return ret
}

func buildExperimentalCiphers(skipSets [][]int, ciphers chan<- cipher.Cipher) {
	m := gematria.Len()
	ciphers <- direct.New()
	for offset := 0; offset < m; offset++ {
		for _, ss := range skipSets {
			// Simple
			ciphers <- experimental.NewStateless("D(x) = (VR(x) + off) mod 29", func(i int, r rune) int {
				return (gematria.ValueOfRune(r) + offset) % m
			}, ss)
			// ===========================================================================================
			// Mu-based
			ciphers <- experimental.NewStateless("D(X) = (IR(x) - mu(x) + off) mod 29", func(i int, r rune) int {
				return (m*m*m*m + gematria.IndexOfRune(r) - mathutil.Mu(i) + offset) % m
			}, ss)
			ciphers <- experimental.NewStateless("D(x) = (VR(x) - mu(x) + off) mod 29", func(i int, r rune) int {
				return (m*m*m*m + gematria.ValueOfRune(r) - mathutil.Mu(i) + offset) % m
			}, ss)
			// ===========================================================================================
			// Phi-based
			ciphers <- experimental.NewStateless("D(x) = (IR(x) - phi(x) + off) mod 29", func(i int, r rune) int {
				return (m*m*m*m + gematria.IndexOfRune(r) - mathutil.Phi(i) + offset) % m
			}, ss)
			ciphers <- experimental.NewStateless("D(x) = (VR(x) - phi(x) + off) mod 29", func(i int, r rune) int {
				return (m*m*m*m + gematria.ValueOfRune(r) - mathutil.Phi(i) + offset) % m
			}, ss)
			ciphers <- experimental.NewStateless("D(x) = (IR(x) - phi(primes(x)) + off) mod 29", func(i int, r rune) int {
				return (m*m*m*m + gematria.IndexOfRune(r) - mathutil.Phi(mathutil.PrimeAt(i)) + offset) % m
			}, ss)
			ciphers <- experimental.NewStateless("D(x) = (VR(x) - phi(primes(x)) + off) mod 29", func(i int, r rune) int {
				return (m*m*m*m + gematria.ValueOfRune(r) - mathutil.Phi(mathutil.PrimeAt(i)) + offset) % m
			}, ss)
			// ===========================================================================================
			// Stateful
			for _, base := range []int{2, 8, 10, 29, 59, 60} {
				base := base
				ciphers <- experimental.New("B???", func(i int, r rune, state int) (int, int) {
					newState := (m*m + state*base - gematria.IndexOfRune(r)) % m
					return (m*m*m*m + state + offset) % m, newState
				}, ss)
				ciphers <- experimental.New("B???", func(i int, r rune, state int) (int, int) {
					newState := (m*m + state*base - gematria.ValueOfRune(r)) % m
					return (m*m*m*m + state + offset) % m, newState
				}, ss)
			}
			ciphers <- experimental.New("P???", func(i int, r rune, state int) (int, int) {
				newState := (m*m + state - gematria.ValueOfRune(r)) % m
				return (m*m*m*m + state + offset) % m, newState
			}, ss)
			ciphers <- experimental.New("P???", func(i int, r rune, state int) (int, int) {
				newState := (m*m + state - gematria.IndexOfRune(r) - mathutil.PrimeAt(i)) % m
				return (m*m*m*m + state + offset) % m, newState
			}, ss)
			ciphers <- experimental.New("P???", func(i int, r rune, state int) (int, int) {
				newState := (m*m + state + gematria.ValueOfRune(r) - mathutil.PrimeAt(i)) % m
				return (m*m*m*m + state + offset) % m, newState
			}, ss)
			ciphers <- experimental.New("P???", func(i int, r rune, state int) (int, int) {
				newState := (m*m + state - gematria.ValueOfRune(r)*mathutil.PrimeAt(i)) % m
				return (m*m*m*m + state + offset) % m, newState
			}, ss)
			// ===========================================================================================
			// Others
			ciphers <- experimental.NewStateless("D(x) = (IR(x) - primes(x) + off) mod 29", func(i int, r rune) int {
				return (m*m*m*m + gematria.IndexOfRune(r) - mathutil.PrimeAt(i) + offset) % m
			}, ss)
			ciphers <- experimental.NewStateless("D(x) = (VR(x) - primes(x) + off) mod 20", func(i int, r rune) int {
				return (m*m*m*m + gematria.ValueOfRune(r) - mathutil.PrimeAt(i) + offset) % m
			}, ss)
			// ===========================================================================================
			// Others
			ciphers <- experimental.NewStateless("D(x) = ((29-IR(x)-1) * 60^x + off) mod 29", func(i int, r rune) int {
				return (m*m*m*m + (m-gematria.IndexOfRune(r)-1)*mathutil.ModPow(60, i, m) + offset) % m
			}, ss)
			ciphers <- experimental.NewStateless("D(x) = (VR(x) * 60^x + off) mod 29", func(i int, r rune) int {
				return (m*m*m*m + gematria.ValueOfRune(r)*mathutil.ModPow(60, i, m) + offset) % m
			}, ss)
		}
	}
}

func buildAffineCiphers(skipSets [][]int, ciphers chan<- cipher.Cipher) {
	m := gematria.Len()
	for a := 0; a < m; a++ {
		if !mathutil.Coprime(a, m) {
			continue
		}
		for b := 0; b < m; b++ {
			ciphers <- affine.New(a, b)
		}
	}
}

func buildVigenereCiphers(skipSets [][]int, ciphers chan<- cipher.Cipher) {
	// m := gematria.Len()
	keywords := []string{
		"ADHERENCE",
		"ADHERE",
		"AETHEREAL",
		"AFFINE",
		"AMASS",
		"ANALOG",
		"BEHAVIORS",
		"BUFFERS",
		"CABAL",
		"CARNAL",
		"CERTAINTY",
		"CICADA",
		"CIRCUMFERENCE",
		"CIRCUMFERENCES",
		"COMMAND",
		"CONSUME",
		"CONSUMPTION",
		"DECEPTION",
		"DIVINITY",
		"DOGMA",
		"EMERGE",
		"ENCRYPT",
		"ENCRYPTED",
		"ENCRYPTING",
		"ENLIGHTENED",
		"ENLIGHTEN",
		"FORM",
		"HOLY",
		"ILLUSIONS",
		"INHABITING",
		"INHABIT",
		"INNOCENCE",
		"INSTAR",
		"IRRITATE",
		"KOAN",
		"MIND",
		"MOBIUS",
		"MOURNFUL",
		"MOURN",
		"OBSCURA",
		"OWN",
		"PARABLE",
		"PILGRIMAGE",
		"PILGRIM",
		"PRESERVATION",
		"PRIMALITY",
		"PRIMAL",
		"PRIMES",
		"REALITY",
		"SACRED",
		"SELF",
		"SHADOWS",
		"TUNNELING",
		"TRUTH",
		"UNREASONABLE",
		"VOID",
		"YOUR",
		"YOURSELF",
	}
	// for mask := 1; mask < (1 << len(keywords)); mask++ {
	// 	for _, ss := range skipSets {
	// 		for offset := 0; offset < 29; offset++ {
	// 			ks := []string{}
	// 			for i := 0; i < len(keywords); i++ {
	// 				if (mask & (1 << i)) != 0 {
	// 					ks = append(ks, keywords[i])
	// 				}
	// 			}
	// 			ciphers <- vigenere.New(ks, offset, ss)
	// 		}
	// 	}
	// }

	for i := 0; i < len(keywords); i++ {
		for j := i + 1; j < len(keywords); j++ {
			for k := j + 1; k < len(keywords); k++ {
				k1, k2, k3 := keywords[i], keywords[j], keywords[k]
				for _, ss := range skipSets {
					ciphers <- vigenere.New([]string{gematria.Encode(k1)}, ss)
					ciphers <- vigenere.New([]string{gematria.Encode(k1), gematria.Encode(k2)}, ss)
					ciphers <- vigenere.New([]string{gematria.Encode(k1), gematria.Encode(k2), gematria.Encode(k3)}, ss)
				}
			}
		}
	}
	for _, keyword := range keywords {
		for _, r := range keyword {
			for _, ss := range skipSets {
				ciphers <- vigenere.New([]string{gematria.Encode(strings.ReplaceAll(keyword, string(r), "F"))}, ss)
			}
		}
	}
}

func buildCandidateCiphers(ciphers chan<- cipher.Cipher) {
	// skipSets := buildSkipSets(msg)
	skipSets := [][]int{[]int{}}
	buildExperimentalCiphers(skipSets, ciphers)
	buildAffineCiphers(skipSets, ciphers)
	buildVigenereCiphers(skipSets, ciphers)
	close(ciphers)
}

func revRunes(s string) string {
	rs := []rune(s)
	for l, r := 0, len(rs)-1; l < r; {
		if !gematria.IsValidRune(rs[l]) {
			l++
			continue
		}
		if !gematria.IsValidRune(rs[r]) {
			r--
			continue
		}
		rs[l], rs[r] = rs[r], rs[l]
		l++
		r--
	}
	return string(rs)
}

func eval(sections []data.Section, ciphers <-chan cipher.Cipher) {
	for dec := range ciphers {
		for _, section := range sections {
			msg := section.Text()[:*prefixLen]
			if section.ID == "unsolved3" {
				msg = section.Text()[100 : 100+*prefixLen]
			}
			output := dec.Decode(msg)
			score := dict.Score(output)
			if score >= *matchScore {
				fmt.Printf("\033[32mid=%s score=%.3f cipher=%s\n%s(...)\n\n\033[39m", section.ID, score, dec.ID(), output)
			}
		}
	}
}

func main() {
	flag.Parse()
	selected := map[string]bool{}
	for _, id := range strings.Split(*selectedSections, ",") {
		if len(id) == 0 {
			continue
		}
		selected[id] = true
	}

	sections := []data.Section{
		data.Unsolved0,
		data.Unsolved1,
		data.Unsolved2,
		data.Unsolved3,
		data.Unsolved4,
		data.Unsolved5,
		data.Unsolved6,
		data.Unsolved7,
	}

	fmt.Printf("Searching on %d CPUs...\n", runtime.NumCPU()-1)
	var wg sync.WaitGroup
	ciphers := make(chan cipher.Cipher)
	for i := 0; i < runtime.NumCPU()-1; i++ {
		wg.Add(1)
		go func() {
			eval(sections, ciphers)
			wg.Done()
		}()
	}
	buildCandidateCiphers(ciphers)
	wg.Wait()
}
