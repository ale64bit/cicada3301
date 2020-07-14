package main

import (
	"flag"
	"fmt"
	"runtime"
	"strings"
	"sync"

	"cicada/cipher"
	"cicada/cipher/affine"
	"cicada/cipher/caesar"
	"cicada/cipher/cat"
	"cicada/cipher/experimental"
	"cicada/cipher/fskip"
	"cicada/cipher/identity"
	"cicada/cipher/prime"
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
	for k := 0; k < 1; k++ {
		var loc []int
		index := 0
		for _, r := range msg {
			if !gematria.IsValidRune(r) {
				continue
			}
			if id := gematria.IndexOfRune(r); id == k {
				loc = append(loc, index)
			}
			index++
		}
		n := len(loc)
		for mask := 0; mask < (1 << n); mask++ {
			var skips []int
			for i := 0; i < n; i++ {
				if (mask & (1 << i)) != 0 {
					skips = append(skips, loc[i])
				}
			}
			ret = append(ret, skips)
		}
	}
	return ret
}

func buildExperimentalCiphers(ciphers chan<- cipher.Cipher) {
	m := gematria.Len()
	z := mathutil.Pow(m, 6)
	ciphers <- identity.New()
	base := []cipher.Cipher{
		experimental.NewStateless("rev", func(i int, r rune) int { return 28 - gematria.IndexOfRune(r) }),
		experimental.NewStateless("IR", func(_ int, r rune) int { return (z - gematria.IndexOfRune(r)) % m }),
		experimental.NewStateless("VR", func(_ int, r rune) int { return (z - gematria.ValueOfRune(r)) % m }),
		prime.New(),
		experimental.NewStateless("muIdx", func(i int, r rune) int { return (z - mathutil.Mu(i)) % m }),
		experimental.NewStateless("muIR", func(i int, r rune) int { return (z - mathutil.Mu(gematria.IndexOfRune(r))) % m }),
		experimental.NewStateless("muVR", func(i int, r rune) int { return (z - mathutil.Mu(gematria.ValueOfRune(r))) % m }),
		experimental.NewStateless("phiIdx", func(i int, r rune) int { return (z - mathutil.Phi(i)) % m }),
		experimental.NewStateless("phiIR", func(i int, r rune) int { return (z - mathutil.Phi(gematria.IndexOfRune(r))) % m }),
		experimental.NewStateless("phiVR", func(i int, r rune) int { return (z - mathutil.Phi(gematria.ValueOfRune(r))) % m }),
		experimental.NewStateless("pinv", func(i int, r rune) int {
			return (gematria.IndexOfRune(r) * mathutil.ModPow(mathutil.PrimeAt(i), m-2, m)) % m
		}),
		experimental.New("st_ixor", func(i int, r rune, state int) (int, int) { return state % m, (state ^ gematria.IndexOfRune(r)) % m }),
		experimental.New("st_vxor", func(i int, r rune, state int) (int, int) { return state % m, (state ^ gematria.ValueOfRune(r)) % m }),
		experimental.New("st_b10", func(i int, r rune, state int) (int, int) {
			return state % m, (z + state*10 - gematria.IndexOfRune(r)) % m
		}),
		experimental.New("st_b29", func(i int, r rune, state int) (int, int) {
			return state % m, (z + state*29 - gematria.IndexOfRune(r)) % m
		}),
		experimental.New("st_b60", func(i int, r rune, state int) (int, int) {
			return state % m, (z + state*60 - gematria.IndexOfRune(r)) % m
		}),
		experimental.New("st_b1033", func(i int, r rune, state int) (int, int) {
			return state % m, (z + state*1033 - gematria.IndexOfRune(r)) % m
		}),
		experimental.New("st_b3301", func(i int, r rune, state int) (int, int) {
			return state % m, (z + state*3301 - gematria.IndexOfRune(r)) % m
		}),
	}
	for shift := 0; shift < m; shift++ {
		for mask := 1; mask < (1 << len(base)); mask++ {
			var cs []cipher.Cipher
			for i := 0; i < len(base); i++ {
				if (mask & (1 << i)) != 0 {
					cs = append(cs, base[i])
				}
			}
			cs = append(cs, caesar.New(shift))
			ciphers <- cat.New(cs...)
		}
	}
}

func buildAffineCiphers(ciphers chan<- cipher.Cipher) {
	m := gematria.Len()
	for a := 0; a < m; a++ {
		if !mathutil.Coprime(a, m) {
			continue
		}
		for b := 0; b < m; b++ {
			for shift := 0; shift < m; shift++ {
				ciphers <- cat.New(affine.New(a, b), caesar.New(shift))
			}
		}
	}
}

func buildVigenereCiphers(ciphers chan<- cipher.Cipher) {
	m := gematria.Len()
	keywords := []string{
		"ADHERENCE",
		"ADHERE",
		"AETHEREAL",
		"AFFINE",
		"AMASS",
		"ANALOG",
		"ANCIENT",
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
		"CROSSROAD",
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
		"WAY",
		"YOUR",
		"YOURSELF",
	}
	for _, keyword := range keywords {
		for _, r := range keyword {
			for shift := 0; shift < m; shift++ {
				ciphers <- cat.New(vigenere.New([]string{gematria.Encode(strings.ReplaceAll(keyword, string(r), "F"))}), caesar.New(shift))
			}
		}
	}
	for i := 0; i < len(keywords); i++ {
		for j := i + 1; j < len(keywords); j++ {
			for k := j + 1; k < len(keywords); k++ {
				k1, k2, k3 := keywords[i], keywords[j], keywords[k]
				for shift := 0; shift < m; shift++ {
					ciphers <- cat.New(vigenere.New([]string{gematria.Encode(k1)}), caesar.New(shift))
					ciphers <- cat.New(vigenere.New([]string{gematria.Encode(k1), gematria.Encode(k2)}), caesar.New(shift))
					ciphers <- cat.New(vigenere.New([]string{gematria.Encode(k1), gematria.Encode(k2), gematria.Encode(k3)}), caesar.New(shift))
				}
			}
		}
	}
}

func buildCandidateCiphers(ciphers chan<- cipher.Cipher) {
	buildExperimentalCiphers(ciphers)
	buildAffineCiphers(ciphers)
	buildVigenereCiphers(ciphers)
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
	bestScore := 0.35
	var skipSets [][][]int
	for _, section := range sections {
		msg := section.Text()[:*prefixLen]
		skipSets = append(skipSets, buildSkipSets(msg))
	}
	for dec := range ciphers {
		for i, section := range sections {
			msg := section.Text()[:*prefixLen]
			for _, ss := range skipSets[i] {
				c := cat.New(dec, fskip.New(ss...))
				plainText := gematria.Decode(c.Decode(msg))
				score := dict.Score(plainText)
				if score >= *matchScore || score > bestScore {
					fmt.Printf("\033[32mid=%s score=%.3f cipher=%s\n%s(...)\n\n\033[39m", section.ID, score, c.ID(), plainText)
					bestScore = score
				}
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
		data.Unsolved6Prefix,
		data.Unsolved6Body,
		data.Unsolved6Suffix,
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
