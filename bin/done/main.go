package main

import (
	"flag"
	"fmt"

	"cicada/data"
	"cicada/decoder"
	"cicada/decoder/affine"
	"cicada/decoder/direct"
	"cicada/decoder/totient"
	"cicada/decoder/vigenere"
	"cicada/gematria"
)

func main() {
	flag.Parse()
	selected := map[string]bool{}
	for _, s := range flag.Args() {
		selected[s] = true
	}

	entries := []struct {
		section data.Section
		decoder decoder.Decoder
	}{
		/*  0 */ {data.AWarning, affine.New(gematria.Len()-1, gematria.Len()-1, nil)},
		/*  1 */ {data.Welcome, vigenere.New(gematria.Encode("DIUINITY"), []int{48, 74, 84, 132, 159, 160, 250, 421, 443, 465, 514})},
		/*  2 */ {data.SomeWisdom, direct.New(0)},
		/*  3 */ {data.Koan1, affine.New(gematria.Len()-1, 2, nil)},
		/*  4 */ {data.LossOfDivinity, direct.New(0)},
		/*  5 */ {data.Koan2, vigenere.New(gematria.Encode("FIRFUMFERENFE"), []int{49, 58})},
		/*  6 */ {data.AnInstruction, direct.New(0)},
		/*  7 */ // TODO {data.Unsolved0, identity.New()},
		/*  8 */ // TODO {data.Unsolved1, identity.New()},
		/*  9 */ // TODO {data.Unsolved2, identity.New()},
		/* 10 */ // TODO {data.Unsolved3, identity.New()},
		/* 11 */ // TODO {data.Unsolved4, identity.New()},
		/* 12 */ // TODO {data.Unsolved5, identity.New()},
		/* 13 */ // TODO {data.Unsolved6, identity.New()},
		/* 14 */ // TODO {data.Unsolved7, identity.New()},
		/* 15 */ {data.AnEnd, totient.New(0, []int{56})},
		/* 16 */ {data.Parable, direct.New(0)},
	}
	for _, e := range entries {
		if len(selected) == 0 || selected[e.section.ID] {
			result, wordStream := e.decoder.DecodeStream(e.section.Text)
			fmt.Printf("\033[33m# id=%s len=%d markings=%v score=%.3f method=%s\n\033[39m%s\n\n",
				e.section.ID,
				gematria.RuneCount(e.section.Text),
				e.section.Markings,
				decoder.Score(wordStream),
				e.decoder.Name(),
				result)
		}
	}
}
