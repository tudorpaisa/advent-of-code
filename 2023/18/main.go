package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
	"math"
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
	Right = iota
	Down
	Left
	Up
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
	curr := Coordinates{0, 0}
	out = append(out, curr)

	for _, i := range actions {
		var v Coordinates
		switch i.direction {
		case Up:
			// curr = Coordinates{curr.x, curr.y - i.size}
			for j := 1; j <= i.size; j++ {
				v = Coordinates{curr.x, curr.y - j}
				out = append(out, v)
			}
			curr = v
		case Down:
			// curr = Coordinates{curr.x, curr.y + i.size}
			for j := 1; j <= i.size; j++ {
				v = Coordinates{curr.x, curr.y + j}
				out = append(out, v)
			}
			curr = v
		case Left:
			// curr = Coordinates{curr.x - i.size, curr.y}
			for j := 1; j <= i.size; j++ {
				v = Coordinates{curr.x - j, curr.y}
				out = append(out, v)
			}
			curr = v
		case Right:
			// curr = Coordinates{curr.x + i.size, curr.y}
			for j := 1; j <= i.size; j++ {
				v = Coordinates{curr.x + j, curr.y}
				out = append(out, v)
			}
			curr = v
		default:
			panic("Wrong Direction")
		}
		// out = append(out, curr)
	}

	return out
}

type MapEdges struct {
	minX int
	maxX int
	minY int
	maxY int
}

func getMapEdges(coords []Coordinates) MapEdges {
	minX, minY, maxX, maxY := 0, 0, 0, 0

	for _, i := range coords {
		if i.x > maxX { maxX = i.x }
		if i.x < minX { minX = i.x }
		if i.y > maxY { maxY = i.y }
		if i.y < minY { minY = i.y }
	}
	return MapEdges{minX, maxX, minY, maxY}
}

func initMap(edges MapEdges) []string {
	m := []string{}
	height := int( math.Abs(float64(edges.minY)) + math.Abs(float64(edges.maxY)) ) + 1
	width := int( math.Abs(float64(edges.minX)) + math.Abs(float64(edges.maxX)) ) + 1
	for i:=0; i < height; i++ {
		m = append(m, strings.Repeat(".", width))
	}
	return m
}

func printMap(m []string) {
	for _, i := range m {
		fmt.Printf("%s\n", i)
	}
}

func drawOutline(m []string, coords []Coordinates, edges MapEdges) []string {
	for _, i := range coords {
		// add offset
		y := i.y + int(math.Abs(float64(edges.minY)))
		x := i.x + int(math.Abs(float64(edges.minX)))

		s := m[y]
		s = s[:x] + "#" + s[min(x+1, len(s)):]
		m[y] = s
	}
	return m
}

func fillMap(x, y int, m []string) []string {
	if x < 0 || x >= len(m[0]) || y < 0 || y >= len(m) { return m }

	if m[y][x] == '#' {
		return m
	}

	s := m[y]
	s = s[:x] + "#" + s[x+1:]
	m[y] = s
	m = fillMap(x - 1, y, m)
	m = fillMap(x + 1, y, m)
	m = fillMap(x, y - 1, m)
	m = fillMap(x, y + 1, m)

	return m
}

func countFilled(m []string) int {
	out := 0
	for _, i := range m {
		for _, j := range i { if j == '#' { out++ } }
	}
	return out
}

func hexToInt(a string) int {
	i, err := strconv.Atoi(string(a))
	if err != nil {
		switch a {
		case "a":
			return 10
		case "b":
			return 11
		case "c":
			return 12
		case "d":
			return 13
		case "e":
			return 14
		case "f":
			return 15
		default:
			panic(err)
		}
	}
	return i
}

func updateActions(actions []Action) []Action {
	out := []Action{}
	for _, i := range actions {
		color := strings.Replace(i.color, "#", "", 1)

		strDirection := string(color[len(color)-1])
		var dir Direction
		switch strDirection {
		case "0":
			dir = Right
		case "1":
			dir = Down
		case "2":
			dir = Left
		case "3":
			dir = Up
		default:
			panic("Wrong Direction")
		}

		nPow := len(color)-2
		num := 0
		for j := 0; j < len(color); j++ {
			x := hexToInt(string(color[j]))
			num += x * int(math.Pow(16, float64(nPow)))
			nPow--
		}

		i.size = num
		i.direction = dir

		out = append(out, i)
	}
	return out
}

func prod(x1 int, y1 int, x2 int, y2 int) int {
	return (x1 * y2) - (x2 * y1)
}

func shoelace(path []Coordinates) float64 {
	out := 0

	for i := 1; i < len(path) - 1; i++ {
		out += prod(path[i].x, path[i].y, path[i+1].x, path[i+1].y)
	}
	out += prod(path[len(path)-1].x, path[len(path)-1].y, path[0].x, path[0].y)

	return float64(out) / 2
}

func offsetCoords(coords []Coordinates, edges MapEdges) []Coordinates {
	out := []Coordinates{}
	for _, i := range coords {
		out = append(out, Coordinates{i.x+edges.minX, i.y+edges.minY})
	}
	return out
}

func solution1(data []string) int {
	actions := preprocess(data)
	coords := traverse(actions)
	return int(shoelace(coords)) + (len(coords) / 2) + 1

	// This is with polyfill. Uncomment for some fun

	// edges := getMapEdges(coords)
	// m := initMap(edges)
	// m = drawOutline(m, coords, edges)
	// m = fillMap(coords[1].x+  int(math.Abs( float64( edges.minX ) ) ) +1, coords[0].y+int(math.Abs( float64( edges.minY ) ) )+1, m)
	// // printMap(m)
	// return countFilled(m)
}

func solution2(data []string) int {
	actions := preprocess(data)
	actions = updateActions(actions)
	coords := traverse(actions)
	return int(shoelace(coords)) + (len(coords) / 2) + 1
}


func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}
