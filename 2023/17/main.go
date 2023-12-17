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

func preprocess(data []string) [][]int {
	out := [][]int{}
	for _, i := range data {
		row := []int{}
		for _, c := range i {
			n, err := strconv.Atoi(string(c))
			if err != nil { panic(err) }
			row = append(row, n)
		}
		out = append(out, row)
	}
	return out
}

func printMap(m [][]int) {
	for _, i := range m {
		for _, j := range i {
			fmt.Printf("%d ", j)
		}
		fmt.Print("\b\n")
	}
}

type Direction int
const (
	North Direction= iota
	South
	West
	East
)

type Node struct {
	x int
	y int
	direction Direction
	dirCount int
}

type WeightedNode struct {
	node Node
	cost int
}

func directionToCoords(d Direction) (int, int) {
	switch d {
	case North:
		return 0, -1
	case South:
		return 0, 1
	case West:
		return -1, 0
	case East:
		return 1, 0
	default:
		panic("Wrong direction")
	}
}

func getCheapest(a []WeightedNode) ( WeightedNode, int ) {
	minCost := math.MaxInt
	var out WeightedNode
	var idx int
	for i, j := range a {
		if j.cost < minCost {
			idx = i
			out = j
			minCost = j.cost
		}
	}
	return out, idx
}

func getOppositeDirection(d Direction) Direction {
	switch d {
	case North:
		return South
	case South:
		return North
	case West:
		return East
	case East:
		return West
	default:
		panic("Wrong direction")
	}
}

func isOutOfBounds(x, y int, m [][]int) bool {
	if y < 0 || x < 0 || y >= len(m) || x >= len(m[0]) {return true}
	return false
}

func findPath(dataMap [][]int, startX, startY, destX, destY, minC, maxC int) int {
	visited := make(map[Node]bool)
	// too lazy to implement a priority queue
	queue := []WeightedNode{{Node{startX, startY, East, 1}, 0} , {Node{startX, startY, South, 1}, 0}}

	for len(queue) > 0 {
		wnode, idx := getCheapest(queue)
		if idx == -1 { panic("-1")}
		queue = append(queue[:idx], queue[idx+1:]...)

		if _, ok := visited[wnode.node]; ok {
			continue
		}
		visited[wnode.node] = true

		xInc, yInc := directionToCoords(wnode.node.direction)
		newX := wnode.node.x + xInc
		newY := wnode.node.y + yInc
		if isOutOfBounds(newX, newY, dataMap) { continue }

		newCost := wnode.cost + dataMap[newY][newX]
		if wnode.node.dirCount >= minC && wnode.node.dirCount <= maxC {
			if newX == destX && newY == destY {
				return newCost
			}
		}

		for _, dir := range []Direction{North, South, West, East} {
			// check rev
			if getOppositeDirection(wnode.node.direction) == dir { continue }

			newDirCount := 1
			if wnode.node.direction == dir {newDirCount += wnode.node.dirCount}
			if (wnode.node.direction != dir && wnode.node.dirCount < minC) || newDirCount > maxC  {
				continue
			}
			newWNode := WeightedNode{ Node{newX, newY, dir, newDirCount}, newCost }
			queue = append(queue, newWNode)
		}
	}

	return -1
}

func solution1(data []string) int {
	m := preprocess(data)
	return findPath(m,0, 0, len(m[0])-1, len(m)-1, 1, 3)
}

func solution2(data []string) int {
	m := preprocess(data)
	return findPath(m,0, 0, len(m[0])-1, len(m)-1, 4, 10)
}


func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

