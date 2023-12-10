package main

import (
	"fmt"
	"log"
	"os"
	"strings"
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

const START_POSITION string = "S"
const VERTICAL_PIPE string = "|"
const HORIZONTAL_PIPE string = "-"
const NORTH_EAST_BEND string = "L"
const NORTH_WEST_BEND string = "J"
const SOUTH_WEST_BEND string = "7"
const SOUTH_EAST_BEND string = "F"
const GROUND string = "."

type Direction int

const (
	North Direction = iota
	South
	East
	West
	ErrorDirection
)

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

func findStart(data []string) (int, int, bool) {
	re := regexp.MustCompile(START_POSITION)

	for y, row := range data {
		matchIdx := re.FindIndex([]byte(row))
		if len(matchIdx) == 2 {
			return matchIdx[0], y, true
		}
	}
	return -1, -1, false
}

type Pipe struct {
	x int
	y int
	pipeType string
}

func withinBoundaries(x int, y int, maxX int, maxY int) bool {
	return x >= 0 && y >= 0 && x <= maxX && y <= maxY
}

func getDirection(curr Pipe, next Pipe) Direction {
	if curr.x == next.x && curr.y > next.y {
		return North
	} else if curr.x == next.x && curr.y < next.y {
		return South
	} else if curr.x > next.x && curr.y == next.y {
		return West
	} else if curr.x < next.x && curr.y == next.y {
		return East
	}

	panic("Pipes are going diagonally or not moving?!")
}

func pipesCompatible(curr Pipe, next Pipe) bool {
	if next.pipeType == GROUND { return false }
	if next.pipeType == START_POSITION { return true }

	direction := getDirection(curr, next)

	switch curr.pipeType {
	case VERTICAL_PIPE:  // |
		if next.pipeType == HORIZONTAL_PIPE {return false}

		if direction == North && next.pipeType == NORTH_EAST_BEND {return false}
		if direction == North && next.pipeType == NORTH_WEST_BEND {return false}

		if direction == South && next.pipeType == SOUTH_EAST_BEND {return false}
		if direction == South && next.pipeType == SOUTH_WEST_BEND {return false}

	case HORIZONTAL_PIPE: // -
		if next.pipeType == VERTICAL_PIPE {return false}

		if direction == East && next.pipeType == NORTH_EAST_BEND {return false}
		if direction == East && next.pipeType == SOUTH_EAST_BEND {return false}

		if direction == West && next.pipeType == NORTH_WEST_BEND {return false}
		if direction == West && next.pipeType == SOUTH_WEST_BEND {return false}

	case NORTH_EAST_BEND: // L
		if next.pipeType == NORTH_EAST_BEND {return false}

		if direction == North && next.pipeType == NORTH_WEST_BEND {return false}
		if direction == North && next.pipeType == HORIZONTAL_PIPE {return false}

		if direction == East && next.pipeType == SOUTH_EAST_BEND {return false}
		if direction == East && next.pipeType == VERTICAL_PIPE {return false}

	case NORTH_WEST_BEND:  // J
		if next.pipeType == NORTH_WEST_BEND {return false}

		if direction == North && next.pipeType == NORTH_EAST_BEND {return false}
		if direction == North && next.pipeType == HORIZONTAL_PIPE {return false}

		if direction == West && next.pipeType == SOUTH_WEST_BEND {return false}
		if direction == West && next.pipeType == VERTICAL_PIPE {return false}

	case SOUTH_EAST_BEND:  // F
		if next.pipeType == SOUTH_EAST_BEND {return false}

		if direction == South && next.pipeType == SOUTH_WEST_BEND {return false}
		if direction == South && next.pipeType == HORIZONTAL_PIPE {return false}

		if direction == East && next.pipeType == NORTH_EAST_BEND {return false}
		if direction == East && next.pipeType == VERTICAL_PIPE {return false}

	case SOUTH_WEST_BEND:  // 7
		if next.pipeType == SOUTH_WEST_BEND {return false}

		if direction == South && next.pipeType == SOUTH_EAST_BEND {return false}
		if direction == South && next.pipeType == HORIZONTAL_PIPE {return false}

		if direction == East && next.pipeType == NORTH_WEST_BEND {return false}
		if direction == East && next.pipeType == VERTICAL_PIPE {return false}
	}

	return true
}

func getNextDirection(pipe Pipe, direction Direction) Direction {
	if pipe.pipeType == START_POSITION { return direction }

	if pipe.pipeType == VERTICAL_PIPE && direction == North { return North }
	if pipe.pipeType == VERTICAL_PIPE && direction == South { return South }

	if pipe.pipeType == HORIZONTAL_PIPE && direction == East { return East }
	if pipe.pipeType == HORIZONTAL_PIPE && direction == West { return West }

	if pipe.pipeType == NORTH_EAST_BEND && direction == South { return East }
	if pipe.pipeType == NORTH_EAST_BEND && direction == West { return North }

	if pipe.pipeType == NORTH_WEST_BEND && direction == South { return West }
	if pipe.pipeType == NORTH_WEST_BEND && direction == East { return North }

	if pipe.pipeType == SOUTH_EAST_BEND && direction == North { return East }
	if pipe.pipeType == SOUTH_EAST_BEND && direction == West { return South }

	if pipe.pipeType == SOUTH_WEST_BEND && direction == North { return West }
	if pipe.pipeType == SOUTH_WEST_BEND && direction == East { return South }

	panic("Wrong directions!")
}

func getNextXY(pipe Pipe, nextDirection Direction) (int, int) {
	if nextDirection == North {
		return pipe.x, pipe.y - 1
	} else if nextDirection == South {
		return pipe.x, pipe.y + 1
	} else if nextDirection == East {
		return pipe.x + 1, pipe.y
	} else if nextDirection == West {
		return pipe.x - 1, pipe.y
	}

	panic("Wrong XY!")
}

func traversePipes(data []string, initX int, initY int, initDirection Direction) ( []Pipe, bool ) {
	pipesPath := *new([]Pipe)

	maxX := len(data[0]) - 1
	maxY := len(data) - 1

	currPipe := Pipe{initX, initY, string(data[initY][initX])}
	currDirection := initDirection

	for {
		nextDirection := getNextDirection(currPipe, currDirection)

		nextX, nextY := getNextXY(currPipe, nextDirection)
		if !withinBoundaries(nextX, nextY, maxX, maxY) {break}

		nextPipe := Pipe{nextX, nextY, string(data[nextY][nextX])}
		if !pipesCompatible(currPipe, nextPipe) { break }

		if nextPipe.pipeType == START_POSITION {
			return pipesPath, true
		}

		pipesPath = append(pipesPath, nextPipe)

		currPipe = nextPipe
		currDirection = nextDirection
	}

	return []Pipe{}, false
}

func formatPipes(pipes []Pipe) []string {
	out := *new([]string)
	for _, i := range pipes { out = append(out, i.pipeType)}
	return out
}

func getPipesPath(data []string, initX int, initY int) []Pipe {
	var pipes []Pipe
	var okPath bool

	pipes, okPath = traversePipes(data, initX, initY, North)
	if okPath { return pipes }

	pipes, okPath = traversePipes(data, initX, initY, South)
	if okPath { return pipes }

	pipes, okPath = traversePipes(data, initX, initY, East)
	if okPath { return pipes }

	pipes, okPath = traversePipes(data, initX, initY, West)
	if okPath { return pipes }

	panic("No path found!")
}

func prod(x1 int, y1 int, x2 int, y2 int) int {
	// fmt.Printf("(%d * %d) - (%d * %d)\n", x1, y2, x2, y1)
	return (x1 * y2) - (x2 * y1)
}

func shoelace(path []Pipe) float64 {
	out := 0

	for i := 0; i < len(path) - 1; i++ {
		out += prod(path[i].x, path[i].y, path[i+1].x, path[i+1].y)
	}
	out += prod(path[len(path)-1].x, path[len(path)-1].y, path[0].x, path[0].y)

	return float64(out) / 2.0
}

func solution1(data []string) int {
	startX, startY, _ := findStart(data)
	path := getPipesPath(data, startX, startY)

	if len(path) % 2 == 0 {
		return len(path) / 2
	}
	return (len(path) / 2) + 1
}

func solution2(data []string) int {
	startX, startY, _ := findStart(data)
	path := getPipesPath(data, startX, startY)
	area := shoelace(path)
	return int(area - (float64(len(path))/2) + 1)
}

func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

