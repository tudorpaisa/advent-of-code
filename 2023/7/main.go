package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
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


type Hand struct {
	cards string
	payout int
	distribution [][]string
	handType HandType
}

func printHands(hands []Hand) {
	for _, i := range hands {
		fmt.Printf(
			"Cards: %s | Payout: %d | Distribution: %s | Type: %d\n",
			i.cards, i.payout, i.distribution, i.handType)
	}
}

func calculateDistribution(cards string) [][]string {
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
	return out
}

func getType(dist [][]string) HandType {
	if len(dist) != 6 { panic("Invalid input. Expected array of size 6") }

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

func solution1(data []string) int {
	hands := *new([]Hand)
	for _, i := range data {
		s := strings.Split(i, " ")
		cards := s[0]
		bid, err := strconv.Atoi(s[1])
		if err != nil {
			panic(err)
		}

		dist := calculateDistribution(cards)
		handType := getType(dist)

		hands = append(hands, Hand{cards, bid, dist, handType})
	}

	printHands(hands)

	return 0
}

func main() {
	data := readLines("input2.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
}

