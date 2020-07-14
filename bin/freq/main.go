package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"cicada/data"
	"cicada/gematria"
	"cicada/mathutil"
)

var (
	maxCard = flag.Int("max_card", 3, "Max graph cardinality")
	full    = flag.Bool("full", false, "Whether to display the frequency for each rune/digraph")
)

func perc(x, n int) float64 {
	return float64(x) / float64(n) * 100.0
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func graphFreq(section data.Section) {
	type entry struct {
		key   string
		count int
	}

	n := gematria.RuneCount(section.Text())
	text := []rune(gematria.FilterRunes(section.Text()))
	for card := 1; card <= *maxCard; card++ {
		count := map[string]int{}
		for i := 0; i+card-1 < len(text); i += card {
			g := ""
			for j := 0; j < card; j++ {
				g += string(text[i+j])
			}
			count[g]++
		}
		var entries []entry
		for k, v := range count {
			entries = append(entries, entry{string(k), v})
		}
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].count == entries[j].count {
				return entries[i].key < entries[j].key
			}
			return entries[i].count > entries[j].count
		})

		fmt.Printf("\t# %d-graph freq: len_match=%d range=[%5.2f%% %5.2f%%] unique=%4d unique_freq=%6.2f%% unique_ratio=%6.2f%%\n",
			card,
			b2i(n%card == 0),
			perc(entries[len(entries)-1].count, n/card),
			perc(entries[0].count, n/card),
			len(count),
			perc(len(count), mathutil.Pow(n, card)),
			perc(len(count), n/card))
		if *full {
			fmt.Printf("\t\t")
			for _, e := range entries {
				fmt.Printf("%s=%3.1f%% ", e.key, perc(e.count, n/card))
			}
			fmt.Println()
		}
	}
}

func IOC(section data.Section) {
	rs := []rune(gematria.RuneStream(section.Text()))
	fmt.Printf("\t# IOC: ")
	const maxShift = 15
	for shift := 1; shift <= maxShift; shift++ {
		total := 0
		count := 0
		for i := 0; i+maxShift < len(rs); i++ {
			if rs[i] == rs[i+shift] {
				count++
			}
			total++
		}
		fmt.Printf("c(%d)=%.2f%% ", shift, perc(count, total))
	}
	fmt.Println()
}

func lineByLine(section data.Section) {
	fmt.Printf("\t# line-by=line:\n")
	for _, page := range section.Pages {
		lines := strings.Split(page, "\n")
		for _, line := range lines {
			n := 0
			sum := 0
			var values []int
			for _, r := range line {
				if gematria.IsValidRune(r) {
					x := gematria.IndexOfRune(r)
					n++
					sum += x
					values = append(values, x)
				}
			}
			fmt.Printf("\t\tn=%2d sum=%3d [ ", n, sum)
			for _, v := range values {
				fmt.Printf("%2d ", v)
			}
			fmt.Println("]")
		}
		fmt.Println()
	}
	fmt.Println()
}

func analyze(section data.Section) {

	n := gematria.RuneCount(section.Text())
	fmt.Printf("# section=%s len=%d\n", section.ID, n)

	graphFreq(section)
	IOC(section)
	lineByLine(section)

	fmt.Println()
}

func main() {
	flag.Parse()
	selected := map[string]bool{}
	for _, id := range flag.Args() {
		if len(id) == 0 {
			continue
		}
		selected[id] = true
	}

	sections := []data.Section{
		data.Warning,
		data.Welcome,
		data.SomeWisdom,
		data.Koan1,
		data.LossOfDivinity,
		data.Koan2,
		data.AnInstruction,
		data.Unsolved0,
		data.Unsolved1,
		data.Unsolved2,
		data.Unsolved3,
		data.Unsolved4,
		data.Unsolved5,
		data.Unsolved6,
		data.Unsolved7,
		data.AnEnd,
		data.Parable,
	}

	for _, section := range sections {
		if len(selected) == 0 || selected[section.ID] {
			analyze(section)
		}
	}
}
