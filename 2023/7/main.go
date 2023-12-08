package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func readFileRaw(fname string) string {
	file, err := os.ReadFile(fname)

	if err != nil {
		msg := fmt.Sprintf("Encountered an error while reading '%s': %s", fname, err)
		log.Fatal(msg)
	}

	return string(file)
}

func readLines(fname string) []string {
	var data []string = strings.Split(readFileRaw(fname), "\n")
	var out []string = *new([]string)

	for _, row := range data {
		if row == "" {
			continue
		}
		out = append(out, row)
	}
	return out
}

type HandType int

const (
	FiveOfAKind HandType = iota
	FourOfAKind
	FullHouse
	ThreeOfAKind
	TwoPair
	OnePair
	HighCard
)

func (h HandType) toString() string {
	switch h {
	case FiveOfAKind:
		return "FiveOfAKind"
	case FourOfAKind:
		return "FourOfAKind"
	case FullHouse:
		return "FullHouse"
	case ThreeOfAKind:
		return "ThreeOfAKind"
	case TwoPair:
		return "TwoPair"
	case OnePair:
		return "OnePair"
	default:
		return "HighCard"
	}
}

type Hand struct {
	cards        string
	payout       int
	distribution [][]string
	handType     HandType
}

func printHands(hands []Hand) {
	firstLetter := string(hands[0].cards[0])
	for i := 0; i < len(hands); i++ {
		currFirstLetter := string(hands[i].cards[0])
		if firstLetter != currFirstLetter {
			firstLetter = currFirstLetter
			fmt.Print("\n")
		}
		// for _, i := range hands {
		fmt.Printf("%s, ", hands[i].cards)
		// fmt.Printf(
		// 	"Cards: %s | Payout: %d | Distribution: %s | Type: %d\n",
		// 	i.cards, i.payout, i.distribution, i.handType)
	}
	fmt.Print("\n")
}

func calculateDistribution(cards string) (map[string]int, [][]string) {
	counts := make(map[string]int)
	for _, i := range cards {
		j := string(i)
		if n, ok := counts[j]; ok {
			counts[j] = n + 1
		} else {
			counts[j] = 1
		}
	}

	// put cards in array where index == num card in hand
	var out [][]string = make([][]string, 6)

	for card, n := range counts {
		out[n] = append(out[n], card)
	}
	return counts, out
}

const cardRank string = "23456789TJQKA"
const cardRankJokerRule string = "J23456789TQKA"

func isLower(a Hand, b Hand, jokerRule bool) bool {
	if a.cards == b.cards {
		return false
	}
	// fmt.Printf("%s vs %s ", a.cards, b.cards)

	for i := range a.cards {
		var aIdx, bIdx int
		if jokerRule {
			aIdx = strings.Index(cardRankJokerRule, string(a.cards[i]))
			bIdx = strings.Index(cardRankJokerRule, string(b.cards[i]))
		} else {
			aIdx = strings.Index(cardRank, string(a.cards[i]))
			bIdx = strings.Index(cardRank, string(b.cards[i]))
		}

		if aIdx < 0 {
			panic(fmt.Sprintf("Could not find '%s' in charset='%s'", string(a.cards[i]), cardRank))
		}
		if bIdx < 0 {
			panic(fmt.Sprintf("Could not find '%s' in charset='%s'", string(b.cards[i]), cardRank))
		}

		if aIdx > bIdx {
			// fmt.Print("false\n")
			return false
		} else if aIdx < bIdx {
			// fmt.Print("true\n")
			return true
		}
	}
	// fmt.Print("true\n")
	return true
}

func sort(hands []Hand, jokerRule bool) {
	if len(hands) < 2 {
		return
	}

	low, high := 0, len(hands)-1
	mid := len(hands) / 2

	hands[mid], hands[high] = hands[high], hands[mid]

	for i := range hands {
		if isLower(hands[i], hands[high], jokerRule) {
			hands[i], hands[low] = hands[low], hands[i]
			low++
		}
	}

	hands[low], hands[high] = hands[high], hands[low]

	sort(hands[low+1:], jokerRule)
	sort(hands[:low], jokerRule)
}

func hasJoker(a string) bool {
	return strings.Index(a, "J") >= 0
}

func getNJokers(a string) int {
	return strings.Count(a, "J")
}

func promoteHandType(a HandType) HandType {
	switch a {
	case FiveOfAKind:
		return FiveOfAKind
	case FourOfAKind:
		return FiveOfAKind
	case FullHouse:
		return FourOfAKind
	case ThreeOfAKind:
		return FullHouse
	case TwoPair:
		return ThreeOfAKind
	default:
		return TwoPair
	}
}

func getType(dist [][]string) HandType {
	if len(dist) != 6 {
		panic("Invalid input. Expected array of size 6")
	}

	// maybe rewrite this as switch-case
	switch {
	case len(dist[5]) == 1:
		return FiveOfAKind
	case len(dist[4]) == 1:
		return FourOfAKind
	case len(dist[3]) == 1 && len(dist[2]) == 1:
		return FullHouse
	case len(dist[3]) == 1:
		return ThreeOfAKind
	case len(dist[2]) == 2:
		return TwoPair
	case len(dist[2]) == 1:
		return OnePair
	default:
		return HighCard
	}
}

func getMostCommon(counts map[string]int, includeJoker bool) string {
	var mostCommonLetter string = ""
	mostCounts := 0

	for letter, count := range counts {
		if letter == "J" && !includeJoker {
			continue
		}
		if count > mostCounts {
			mostCommonLetter = letter
			mostCounts = count
		}
	}

	return mostCommonLetter
}

func solution1(data []string) int {
	handBins := make(map[HandType][]Hand)
	for _, i := range data {
		s := strings.Split(i, " ")
		cards := s[0]
		bid, err := strconv.Atoi(s[1])
		if err != nil {
			panic(err)
		}

		_, dist := calculateDistribution(cards)
		handType := getType(dist)

		if bin, ok := handBins[handType]; ok {
			handBins[handType] = append(bin, Hand{cards, bid, dist, handType})
		} else {
			handBins[handType] = append(*new([]Hand), Hand{cards, bid, dist, handType})
		}
	}

	for handType, bin := range handBins {
		sort(bin, false)
		handBins[handType] = bin
	}
	sortedCards := *new([]Hand)
	sortedCards = append(sortedCards, handBins[HighCard]...)
	sortedCards = append(sortedCards, handBins[OnePair]...)
	sortedCards = append(sortedCards, handBins[TwoPair]...)
	sortedCards = append(sortedCards, handBins[ThreeOfAKind]...)
	sortedCards = append(sortedCards, handBins[FullHouse]...)
	sortedCards = append(sortedCards, handBins[FourOfAKind]...)
	sortedCards = append(sortedCards, handBins[FiveOfAKind]...)

	out := 0
	for i, hand := range sortedCards {
		rank := i + 1
		out += rank * hand.payout
	}

	return out
}

func solution2(data []string) int {
	handBins := make(map[HandType][]Hand)
	for _, i := range data {
		s := strings.Split(i, " ")
		cards := s[0]
		bid, err := strconv.Atoi(s[1])
		if err != nil {
			panic(err)
		}

		counts, dist := calculateDistribution(cards)

		var handType HandType

		if hasJoker(cards) {
			mostCommonLetter := getMostCommon(counts, false)
			if mostCommonLetter == "" {
				mostCommonLetter = "2"
			}
			inter := strings.ReplaceAll(cards, "J", mostCommonLetter)
			_, interDist := calculateDistribution(inter)
			handType = getType(interDist)
			// fmt.Printf("%s -> %s. handType=%s\n", cards, inter, handType.toString())
		} else {
			handType = getType(dist)
		}

		if bin, ok := handBins[handType]; ok {
			handBins[handType] = append(bin, Hand{cards, bid, dist, handType})
		} else {
			handBins[handType] = append(*new([]Hand), Hand{cards, bid, dist, handType})
		}
	}

	for handType, bin := range handBins {
		sort(bin, true)
		handBins[handType] = bin
		// fmt.Printf("\nHandType: %s\n", handType.toString())
		// printHands(bin)
	}
	sortedCards := *new([]Hand)
	sortedCards = append(sortedCards, handBins[HighCard]...)
	sortedCards = append(sortedCards, handBins[OnePair]...)
	sortedCards = append(sortedCards, handBins[TwoPair]...)
	sortedCards = append(sortedCards, handBins[ThreeOfAKind]...)
	sortedCards = append(sortedCards, handBins[FullHouse]...)
	sortedCards = append(sortedCards, handBins[FourOfAKind]...)
	sortedCards = append(sortedCards, handBins[FiveOfAKind]...)

	out := 0
	for i, hand := range sortedCards {
		rank := i + 1
		out += rank * hand.payout
		fmt.Printf("%s\t%d\t%s\n", hand.cards, hand.payout, hand.handType.toString())
	}

	return out
}

func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
	// fmt.Printf("%s\n", isLower(
	//     Hand{"55553", 0, *new( [][]string ), FiveOfAKind},
	//     Hand{"3J333", 0, *new([][]string), FiveOfAKind}))

	// a := []int{2,4,5,6,3,2,1,2,5421,10, 953, 2321, 2,4,1,0, 0}
	// sortIntArr(a[:])
	// fmt.Printf("%d\n", a)

}
