package main

import (
	"fmt"
	"log"
	"os"
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

type Pos struct {
	x int
	y int
}

func findStartPoint(data []string) Pos {
	for y, row := range data {
		for x, c := range row {
			if c == 'S' {
				return Pos{x, y}
			}
		}
	}

	return Pos{-1,  -1}
}

func printMap(data []string) {
	for _, row := range data {
		fmt.Printf("%s\n", row)
	}
}

func markMap(data []string, positions []Pos) []string {
	for _, p := range positions {
		data[p.y] = data[p.y][:p.x] + "O" + data[p.y][p.x+1:]
	}
	return data
}

type Direction int
const (
	North = iota
	South
	West
	East
)

func searchInDirection(data []string, pos Pos, dir Direction) (Pos, bool) {
	var newPos Pos

	switch dir {
	case North:
		newPos = Pos{pos.x, pos.y - 1}
	case South:
		newPos = Pos{pos.x, pos.y + 1}
	case West:
		newPos = Pos{pos.x - 1, pos.y}
	case East:
		newPos = Pos{pos.x + 1, pos.y}
	default:
		panic("Wrong Direction!")
	}

	if newPos.y >= 0 && newPos.y < len(data) && newPos.x >= 0 && newPos.x < len(data[0]) {
		if data[newPos.y][newPos.x] != '#' {
			return newPos, true
		}
	}
	return Pos{-1, -1}, false
}

func getNextSteps(data []string, pos Pos) []Pos {
	out := []Pos{}

	nPos, ok := searchInDirection(data, pos, North)
	if ok { out = append(out, nPos) }
	sPos, ok := searchInDirection(data, pos, South)
	if ok { out = append(out, sPos) }
	wPos, ok := searchInDirection(data, pos, West)
	if ok { out = append(out, wPos) }
	ePos, ok := searchInDirection(data, pos, East)
	if ok { out = append(out, ePos) }

	return out
}

func stepIter(data []string, startPos Pos, n int) []Pos {
	it := 0

	evolutionMap := make(map[int]map[Pos]bool)
	evolutionMap[it] = map[Pos]bool{startPos: true}

	for it <= n {

		queue := []Pos{}
		for k := range evolutionMap[it] {
			queue = append(queue, k)
		}

		it++
		itPositions := make(map[Pos]bool)

		for len(queue) > 0 {

			pos := queue[0]
			queue = queue[1:]


			newPos := getNextSteps(data, pos)

			for _, i := range newPos {
				itPositions[i] = true
			}
		}

		evolutionMap[it] = itPositions
	}

	out := []Pos{}
	for k := range evolutionMap[n] {
		out = append(out, k)
	}
	return out

}

func solution1(data []string) int {
	// printMap(data)
	startPos := findStartPoint(data)
	// fmt.Printf("StartPos=%d\n", startPos)

	pos := stepIter(data, startPos, 64)
	// fmt.Printf("%d\n", pos)
	// printMap(markMap(data, pos))
	return len(pos)
}

type RelPos struct {
	x int
	y int
	relOffsetX int
	relOffsetY int
}

func searchInRelativeDirection(data []string, pos RelPos, dir Direction) (RelPos, bool) {
	var newPos RelPos

	switch dir {
	case North:
		newPos = RelPos{pos.x, pos.y - 1, pos.relOffsetX, pos.relOffsetY}
	case South:
		newPos = RelPos{pos.x, pos.y + 1, pos.relOffsetX, pos.relOffsetY}
	case West:
		newPos = RelPos{pos.x - 1, pos.y, pos.relOffsetX, pos.relOffsetY}
	case East:
		newPos = RelPos{pos.x + 1, pos.y, pos.relOffsetX, pos.relOffsetY}
	default:
		panic("Wrong Direction!")
	}

	if newPos.y < 0 {
		newPos.y = len(data) - 1
		newPos.relOffsetY = newPos.relOffsetY - 1
	}
	if newPos.y >= len(data) {
		newPos.y = 0
		newPos.relOffsetY = newPos.relOffsetY + 1
	}

	if newPos.x < 0 {
		newPos.x = len(data[0]) - 1
		newPos.relOffsetX = newPos.relOffsetX - 1
	}
	if newPos.x >= len(data[0]) {
		newPos.x = 0
		newPos.relOffsetX = newPos.relOffsetX + 1
	}

	if data[newPos.y][newPos.x] != '#' {
		return newPos, true
	}
	return RelPos{-1, -1, 0, 0}, false
}

func getNextRelativeSteps(data []string, pos RelPos) []RelPos {
	out := []RelPos{}

	nPos, ok := searchInRelativeDirection(data, pos, North)
	if ok { out = append(out, nPos) }
	sPos, ok := searchInRelativeDirection(data, pos, South)
	if ok { out = append(out, sPos) }
	wPos, ok := searchInRelativeDirection(data, pos, West)
	if ok { out = append(out, wPos) }
	ePos, ok := searchInRelativeDirection(data, pos, East)
	if ok { out = append(out, ePos) }

	return out
}

func stepRelIter(data []string, startPos RelPos, n int) map[int][]RelPos {
	it := 0

	evolutionMap := make(map[int]map[RelPos]bool)
	evolutionMap[it] = map[RelPos]bool{startPos: true}

	out := make(map[int][]RelPos)

	for it <= n {
		fmt.Printf("It=%d\r", it)

		queue := []RelPos{}
		for k := range evolutionMap[it] {
			queue = append(queue, k)
		}

		it++
		itPositions := make(map[RelPos]bool)

		for len(queue) > 0 {

			pos := queue[0]
			queue = queue[1:]


			newPos := getNextRelativeSteps(data, pos)

			for _, i := range newPos {
				itPositions[i] = true
			}
		}

		evolutionMap[it] = itPositions
	}

	for it, m := range evolutionMap {
		p := []RelPos{}
		for k := range m {
			p = append(p, k)
		}
		out[it] = p
	}

	return out
}


func solution2(data []string) string {
	startPos := findStartPoint(data)
	startPosRel := RelPos{startPos.x, startPos.y, 0, 0}

	nSteps := 26501365
	halfSize := startPos.x
	size := len(data)
	target := (nSteps - halfSize) / size

	pos := stepRelIter(data, startPosRel, halfSize + (size * 2))
	x := []int{0, 1, 2}
	y := []int{}
	for it := 1; it <= len(pos); it++ {
		p := pos[it]
		if it % size == halfSize {
			y = append(y, len(p))
		}
	}

	out := `Go to Wolfram Alpha, and search for 'quadratic fit calculator'.
For x and y, plug the following values:
`
	out = out + fmt.Sprintf("x=%d\ny=%d\n", x, y)
	out = out + `Once you have you coefficients, plug the a, b, c coefficients in the following formula:
a * (x^2) + b * x + c
`
	out = out + fmt.Sprintf("and solve for x=%d\n", target)

	return out
}


func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %s\n", solution2(data))
}

