package main

import (
	"flag"
	"fmt"
	"strings"

	"cicada/data"
	"cicada/decoder"
	"cicada/decoder/affine"
	"cicada/decoder/direct"
	"cicada/decoder/experimental"
	"cicada/decoder/vigenere"
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
	for k := 0; k < gematria.Len(); k++ {
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

func buildExperimentalDecoders(skipSets [][]int) []decoder.Decoder {
	m := gematria.Len()
	var decoders []decoder.Decoder
	for offset := 0; offset < m; offset++ {
		decoders = append(decoders, direct.New(offset))
		for _, ss := range skipSets {
			decoders = append(decoders, experimental.New(func(i int, r rune) int {
				return (m*m*m*m + gematria.IndexOfRune(r) + mathutil.Mu(i) + offset) % m
			}, ss))
			decoders = append(decoders, experimental.New(func(i int, r rune) int {
				return (m*m*m*m + gematria.ValueOfRune(r) + mathutil.Mu(i) + offset) % m
			}, ss))
			decoders = append(decoders, experimental.New(func(i int, r rune) int {
				return (m*m*m*m + gematria.IndexOfRune(r) + mathutil.Phi(i) + offset) % m
			}, ss))
			decoders = append(decoders, experimental.New(func(i int, r rune) int {
				return (m*m*m*m + gematria.ValueOfRune(r) + mathutil.Phi(i) + offset) % m
			}, ss))
			decoders = append(decoders, experimental.New(func(i int, r rune) int {
				return (m*m*m*m + gematria.IndexOfRune(r) + mathutil.PrimeAt(i) + offset) % m
			}, ss))
			decoders = append(decoders, experimental.New(func(i int, r rune) int {
				return (m*m*m*m + gematria.ValueOfRune(r) + mathutil.PrimeAt(i) + offset) % m
			}, ss))
			decoders = append(decoders, experimental.New(func(i int, r rune) int {
				return (m*m*m*m + (m-gematria.IndexOfRune(r)-1)*mathutil.ModPow(60, i, m) + offset) % m
			}, ss))
			decoders = append(decoders, experimental.New(func(i int, r rune) int {
				return (m*m*m*m + gematria.ValueOfRune(r)*mathutil.ModPow(60, i, m) + offset) % m
			}, ss))
			decoders = append(decoders, experimental.New(func(i int, r rune) int {
				return (m*m*m*m + (m - gematria.IndexOfRune(r) - 1) - (mathutil.PrimeAt(i) - 1) + offset) % m
			}, ss))
			decoders = append(decoders, experimental.New(func(i int, r rune) int {
				return (m*m*m*m + (m - gematria.IndexOfRune(r) - 1) + (mathutil.PrimeAt(i) - 1) + offset) % m
			}, ss))
		}
	}
	return decoders
}

func buildAffineDecoders(skipSets [][]int) []decoder.Decoder {
	var decoders []decoder.Decoder
	for a := 0; a < gematria.Len(); a++ {
		if !mathutil.Coprime(a, gematria.Len()) {
			continue
		}
		for b := 0; b < gematria.Len(); b++ {
			for _, ss := range skipSets {
				decoders = append(decoders, affine.New(a, b, ss))
			}
		}
	}
	return decoders
}

func buildVigenereDecoders(skipSets [][]int) []decoder.Decoder {
	var decoders []decoder.Decoder
	keywords := []string{
		"ADHERENCE",
		"ADHERE",
		"AETHEREAL",
		"AFFINE",
		"AMASS",
		"BEHAVIORS",
		"BUFFERS",
		"CABAL",
		"CARNAL",
		"CERTAINTY",
		"CICADA",
		"CIRCUMFERENCE",
		"CIRCUMFERENCES",
		"CONSUME",
		"DECEPTION",
		"DIVINITY",
		"DOGMA",
		"EMERGE",
		"ENCRYPT",
		"ENCRYPTED",
		"ENCRYPTING",
		"ENLIGHTENED",
		"ENLIGHTEN",
		"HASHES",
		"ILLUSIONS",
		"INHABITING",
		"INHABIT",
		"INNOCENCE",
		"INSTAR",
		"IRRITATED",
		"IRRITATE",
		"KOAN",
		"MOBIUS",
		"MOURNFUL",
		"MOURN",
		"OBSCURA",
		"PARABLE",
		"PILGRIMAGE",
		"PILGRIM",
		"PRIMALITY",
		"PRIMAL",
		"PRIMES",
		"REALITIES",
		"TOTIENT",
		"TUNNELING",
		"UNREASONABLE",
	}
	for _, keyword := range keywords {
		for _, ss := range skipSets {
			decoders = append(decoders, vigenere.New(gematria.Encode(keyword), ss))
			decoders = append(decoders, vigenere.New("ᚫᚦᛖᚱᛠᛚ", ss))
			decoders = append(decoders, vigenere.New("ᛏᚱᚢᚦ", ss))
			decoders = append(decoders, vigenere.New("ᛏᚱᚢᛏᚻ", ss))
		}
		for _, r := range keyword {
			for _, ss := range skipSets {
				decoders = append(decoders, vigenere.New(gematria.Encode(strings.ReplaceAll(keyword, string(r), "F")), ss))
			}
		}
	}
	return decoders
}

func buildCandidateDecoders(msg string) []decoder.Decoder {
	skipSets := buildSkipSets(msg)
	var decoders []decoder.Decoder
	decoders = append(decoders, buildExperimentalDecoders(skipSets)...)
	decoders = append(decoders, buildAffineDecoders(skipSets)...)
	decoders = append(decoders, buildVigenereDecoders(skipSets)...)
	return decoders
}

func eval(section data.Section) {
	msg := section.Text[:*prefixLen]
	decoders := buildCandidateDecoders(msg)
	maxScore := -1.0
	output := ""
	var bestDecoder decoder.Decoder
	for _, d := range decoders {
		s, stream := d.DecodeStream(msg)
		score := decoder.Score(stream)
		if score > maxScore {
			maxScore = score
			bestDecoder = d
			output = s
		}
	}
	if maxScore >= *matchScore {
		fmt.Println("\033[32m!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! MATCH !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\033[39m")
	}
	fmt.Printf("\033[33mid=%s max_score=%.3f decoder=%s\n%s(...)\n\n\033[39m", section.ID, maxScore, bestDecoder.Name(), output)
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
	for _, s := range sections {
		if len(selected) > 0 && !selected[s.ID] {
			continue
		}
		eval(s)
	}
}
