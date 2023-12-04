package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
	"math"
	"regexp"
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

type Card struct {
	winningNumbers []string
	numbers []string
}

func initCard(winningNumbers []string, numbers []string) Card {
	out := Card{winningNumbers: winningNumbers, numbers: numbers}
	return out
}

func extractNumbers(data string) []string {
	re, _ := regexp.Compile("(\\d+)")
	matches := re.FindAll([]byte(data), 99)
	out := *new([]string)
	for _, i := range matches {
		out = append(out, string(i))
	}

	return out
}

func parseStrOfInts(data string) []int {
	split := strings.Split(data, " ")
	out := *new([]int)
	for _, i := range split {
		n, _ := strconv.Atoi(i)
		out = append(out, n)
	}
	return out
}

func parseData(data []string) map[string]Card {
	cards := make(map[string]Card)

	for _, row := range data {
		cardSplit := strings.Split(row, ": ")
		numbersSplit := strings.Split(cardSplit[1], " | ")

		cardN := strings.Replace(cardSplit[0], "Card ", "", 1)
		cardN = strings.Replace(cardN, " ", "", 9)
		// winningNumbers := parseStrOfInts(numbersSplit[0])
		// numbers := parseStrOfInts(numbersSplit[1])

		winningNumbers := extractNumbers(numbersSplit[0])
		numbers := extractNumbers(numbersSplit[1])

		cards[cardN] = initCard(winningNumbers, numbers)
	}

	return cards
}

func sum(a []float64) float64 {
	out := 0.0
	for _, i := range a {
		out += i
	}

	return out
}

func aInB(a string, b []string) bool {
	for _, i := range b {
		if a == i {
			return true
		}
	}
	return false
}

func calculateMatches(card Card) int {
	matches := 0
	for _, number := range card.numbers {
		if aInB(number, card.winningNumbers) {
			matches++
		}
	}
	return matches
}

func processCardWins(cards map[string]Card) map[string]int {
	out := make(map[string]int)
	for n, card := range cards {
		out[n] = calculateMatches(card)
	}
	return out
}

func collectCardNumbers(startCardNum string, increment int, maxCards int) []string  {
	// start, err := strconv.Atoi(strings.Replace(startCardNum, " ", "", 10))
	start, err := strconv.Atoi(startCardNum)
	if err != nil {
		fmt.Printf("ERRRRRORRRRRR!!!!!!!!!!!!!!\n%s\n%s\n", startCardNum,err)
	}

	cardsToGet := *new([]string)
	i := 1

	for i <= increment {
		newCardNum := start + i
		i++
		cardStr := strconv.Itoa(newCardNum)
		cardsToGet = append(cardsToGet, cardStr)
	}
	return cardsToGet
}

func solution1(data []string) float64 {
	cards := parseData(data)
	cardScores := *new([]float64)
	for _, card := range cards {
		matches := calculateMatches(card)

		if matches > 0 {
			var b float64 = float64(matches-1)
			score := math.Pow(2.0, b)
			cardScores = append(cardScores, score)
		}
	}
	return sum(cardScores)
}

func solution2(data []string) int {
	cards := parseData(data)
	cardWins := processCardWins(cards)

	queue := *new([]string)
	for n := range cardWins {
		queue = append(queue, n)
	}

	nScratchCards := 0

	for len(queue) > 0 {
		i := queue[0]
		queue = queue[1:]

		nWins := cardWins[i]
		if nWins > 0 {
			cardsToAdd := collectCardNumbers(i, nWins, len(cards)-1)
			queue = append(cardsToAdd, queue...)
		}

		nScratchCards++
	}

	return nScratchCards
}

func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %f\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}
