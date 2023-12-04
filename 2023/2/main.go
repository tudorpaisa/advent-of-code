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

type Balls struct {
	red int
	green int
	blue int
}

type Game struct {
	number int
	subsets []Balls
}

func initBalls() Balls {
	balls := Balls{}
	balls.red = 0
	balls.green = 0
	balls.blue = 0
	return balls
}

func initGame(number int, subsets []Balls) Game {
	game := Game{}
	game.number = number
	game.subsets = subsets
	return game
}

func getBalls(data []string) Balls {
	balls := initBalls()

	for _, i := range data {
		i_split := strings.Split(i, " ")
		n, _ := strconv.Atoi(i_split[0])
		color := i_split[1]

		switch color {
		case "red":
			balls.red = n
		case "green":
			balls.green = n
		case "blue":
			balls.blue = n
		}
	}

	return balls
}

func extractSubsets(data []string) []Balls {
	subsets := *new([]Balls)

	for _, subset := range data {
		subsets = append(subsets, getBalls(strings.Split(subset, ", ")))
	}

	return subsets
}

func extractGames(data []string) map[string]Game {
	var games map[string]Game = make(map[string]Game)

	for _, row := range data {
		rowSplit := strings.Split(row, ": ")
		gameName := rowSplit[0]
		subsetsArr := strings.Split(rowSplit[1], "; ")


		gameNumberStr := strings.Replace(gameName, "Game ", "", 1)
		gameNumber, _ := strconv.Atoi(gameNumberStr)

		subsets := extractSubsets(subsetsArr)
		game := initGame(gameNumber, subsets)

		games[gameNumberStr] = game
	}

	return games
}

func getMaxBalls(game Game) Balls {
	maxBalls := initBalls()
	for _, subset := range game.subsets {
		if maxBalls.red < subset.red { maxBalls.red = subset.red }
		if maxBalls.green < subset.green { maxBalls.green = subset.green }
		if maxBalls.blue < subset.blue { maxBalls.blue = subset.blue }
	}
	return maxBalls
}

func isWithinConstraint(game Game, constraint Balls) bool {
	maxBalls := getMaxBalls(game)

	redOk := maxBalls.red <= constraint.red
	greenOk := maxBalls.green <= constraint.green
	blueOk := maxBalls.blue <= constraint.blue

	return redOk && greenOk && blueOk
}

func printGames(games map[string]Game) {
	for _, v := range games {
		fmt.Printf("Game %d:\n", v.number)
		for _, game := range v.subsets {
			fmt.Printf("    %d red\t%d green\t%d blue\n", game.red, game.green, game.blue)
		}
	}
}

func sum(a []int) int {
	out := 0
	for _, i := range a {
		out += i
	}

	return out
}

func getMaxPerGame(games map[string]Game) map[string]Balls {
	out := make(map[string]Balls)

	for k, v := range games {
		out[k] = getMaxBalls(v)
	}

	return out
}

func product(a Balls) int {
	return a.red * a.green * a.blue
}

func solution1(data []string) int {
	var games map[string]Game = extractGames(data)

	constraint := initBalls()
	constraint.red = 12
	constraint.green = 13
	constraint.blue = 14

	legalGames := *new([]int)
	for _, game := range games {
		if isWithinConstraint(game, constraint) {
			legalGames = append(legalGames, game.number)
		}
	}

	return sum(legalGames)
}

func solution2(data []string) int {
	var games map[string]Game = extractGames(data)

	maxes := getMaxPerGame(games)
	products := *new([]int)
	for k, v := range maxes {
		p := product(v)
		if p == 0 {
			fmt.Printf("We have a zero product at game %s", k)
		}
		products = append(products, p)
	}

	return sum(products)

}

func main() {
	data := readLines("input1.txt")

	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

