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

type Direction int
const (
	Up = iota
	Down
	Left
	Right
)

type Action struct {
	direction Direction
	size int
	color string
}

type Coordinates struct {
	x int
	y int
}

func preprocess(data []string) []Action {
	out := []Action{}

	for _, i := range data {
		s := strings.Split(i, " ")
		if len(s) != 3 {panic("len != 3")}

		var dir Direction
		switch s[0] {
		case "U":
			dir = Up
		case "D":
			dir = Down
		case "L":
			dir = Left
		case "R":
			dir = Right
		default:
			panic("Wrong Direction")
		}

		size, err := strconv.Atoi(s[1])
		if err != nil {panic(err)}

		color := strings.Replace(strings.Replace(s[2], "(", "", 1), ")", "", 1)

		out = append(out, Action{dir, size, color})
	}

	return out
}

func traverse(actions []Action) []Coordinates {
	out := []Coordinates{}
	curr := Coordinates{1, 1}
	out = append(out, curr)

	for _, i := range actions {
		switch i.direction {
		case Up:
			curr = Coordinates{curr.x, curr.y - i.size}
		case Down:
			curr = Coordinates{curr.x, curr.y + i.size}
		case Left:
			curr = Coordinates{curr.x - i.size, curr.y}
		case Right:
			curr = Coordinates{curr.x + i.size, curr.y}
		default:
			panic("Wrong Direction")
		}
		out = append(out, curr)
	}

	return out
}

func prod(x1 int, y1 int, x2 int, y2 int) int {
	// fmt.Printf("(%d * %d) - (%d * %d)\n", x1, y2, x2, y1)
	return (x1 * y2) - (x2 * y1)
}

func shoelace(path []Coordinates) float64 {
	out := 0

	for i := 0; i < len(path) - 1; i++ {
		out += prod(path[i].x, path[i].y, path[i+1].x, path[i+1].y)
	}
	out += prod(path[len(path)-1].x, path[len(path)-1].y, path[0].x, path[0].y)

	return float64(out) / 2.0
}

func solution1(data []string) int {
	actions := preprocess(data)
	coords := traverse(actions)
	// for _, i := range coords {fmt.Println(i)}
	return int(shoelace(coords))
}

func solution2(data []string) int {
	return 0
}


func main() {
	data := readLines("input2.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}
