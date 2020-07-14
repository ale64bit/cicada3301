package main

import (
	"flag"
	"fmt"

	"cicada/cipher"
	"cicada/cipher/atbash"
	"cicada/cipher/caesar"
	"cicada/cipher/cat"
	"cicada/cipher/fskip"
	"cicada/cipher/identity"
	"cicada/cipher/prime"
	"cicada/cipher/vigenere"
	"cicada/data"
	"cicada/dict"
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
		cipher  cipher.Cipher
	}{
		/*  0 */ {data.Warning, atbash.New()},
		/*  1 */ {data.Welcome, cat.New(vigenere.New([]string{gematria.Encode("DIUINITY")}), fskip.New(48, 74, 84, 132, 159, 160, 250, 421, 443, 465, 514))},
		/*  2 */ {data.SomeWisdom, identity.New()},
		/*  3 */ {data.Koan1, cat.New(atbash.New(), caesar.New(3))},
		/*  4 */ {data.LossOfDivinity, identity.New()},
		/*  5 */ {data.Koan2, cat.New(vigenere.New([]string{gematria.Encode("FIRFUMFERENFE")}), fskip.New(49, 58))},
		/*  6 */ {data.AnInstruction, identity.New()},
		/*  7 */ // TODO {data.Unsolved0, identity.New()},
		/*  8 */ // TODO {data.Unsolved1, identity.New()},
		/*  9 */ // TODO {data.Unsolved2, identity.New()},
		/* 10 */ // TODO {data.Unsolved3, identity.New()},
		/* 11 */ // TODO {data.Unsolved4, identity.New()},
		/* 12 */ // TODO {data.Unsolved5, identity.New()},
		/* 13 */ // TODO {data.Unsolved6, identity.New()},
		/* 14 */ // TODO {data.Unsolved7, identity.New()},
		/* 15 */ {data.AnEnd, cat.New(prime.New(), caesar.New(28), fskip.New(56))},
		/* 16 */ {data.Parable, identity.New()},
	}
	for _, e := range entries {
		if len(selected) == 0 || selected[e.section.ID] {
			cipherText := e.section.Text()
			plainText := gematria.Decode(e.cipher.Decode(cipherText))
			fmt.Printf("\033[33m# id=%s len=%d markings=%v score=%.3f method=%s\n\033[39m%s\n\n",
				e.section.ID,
				gematria.RuneCount(cipherText),
				e.section.Markings,
				dict.Score(plainText),
				e.cipher.ID(),
				plainText)
		}
	}
}

/*

   f(0)=29
   f(1)=59
   f(2)=2
   f(3)=3
   f(4)=149
   f(5)=5
   f(6)=151
   f(7)=7
   f(8)=37
   f(9)=67
   f(10)=97
   f(11)=11
   f(12)=41
   f(13)=13
   f(14)=43
   f(15)=73
   f(16)=103
   f(17)=17
   f(18)=47
   f(19)=19
   f(20)=107
   f(21)=79
   f(22)=109
   f(23)=23
   f(24)=53
   f(25)=83
   f(26)=113
   f(27)=317
   f(28)=173

   [2 3 5 7 11 13 17 19 23 29 31 37 41 43 47 53 59 61 67 71 73 79 83 89 97 101 103 107]
   [2 3 5 7 11 13 17 19 23 29    37 41 43 47 53 59    67    73 79 83    97     103 107 109 113 149 151 173 317]

   SKIPPED: 31 61 71 89 101

*/
