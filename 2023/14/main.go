package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"os/exec"
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

const ROUNDED_ROCK string = "O"
const EMPTY_SPACE string = "."
const CUBE_ROCK string = "#"

type RoundedRock struct {
	x int
	y int
}

func (r RoundedRock) canMoveNorth(a []string) bool {
	if r.y == 0 {return false}
	if string(a[r.y-1][r.x]) == ROUNDED_ROCK {return false}
	if string(a[r.y-1][r.x]) == CUBE_ROCK {return false}

	return true
}

func (r RoundedRock) rollNorth(a []string) []string {
	if !r.canMoveNorth(a) {
		return a
	}

	a[r.y] = a[r.y][:r.x] + EMPTY_SPACE + a[r.y][r.x+1:]
	a[r.y-1] = a[r.y-1][:r.x] + ROUNDED_ROCK + a[r.y-1][r.x+1:]

	r.y--

	return a
}

func getAllRoundedRock(a []string) []RoundedRock {
	out := []RoundedRock{}

	for y, i := range a {
		for x, j := range i {
			if string(j) == ROUNDED_ROCK {
				out = append(out, RoundedRock{x, y})
			}
		}
	}
	return out
}

func printMap(data []string) {
	for i := 0; i<len(data);i++ {fmt.Printf("%s\n", data[i])}
}

func moveAllNorth(data []string, rocks []RoundedRock) []string {
	rocksMap := make(map[int][]*RoundedRock )
	for i := 0; i<len(data); i++ {rocksMap[i] = []*RoundedRock{}}
	for _, i := range rocks { rocksMap[i.y] = append(rocksMap[i.y], &i) }

	for i := 0; i<len(data); i++ {
		nImmovable := 0
		for nImmovable < len(rocksMap[i]) {
			for _, r := range rocksMap[i] {
				data = r.rollNorth(data)
				if !r.canMoveNorth(data) { nImmovable++ }
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()
				fmt.Printf("N Immovable: %d\n", nImmovable)
				printMap(data)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}

	return data
}


func solution1(data []string) int {
	rocks := getAllRoundedRock(data)
	for _, i := range rocks { fmt.Printf("Rock(%d, %d)\n", i.x, i.y) }
	moveAllNorth(data, rocks)
	return 0
}

func solution2(data []string) int {
	return 0
}


func main() {
	data := readLines("input2.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

