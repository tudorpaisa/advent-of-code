package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	// "time"
	// "os/exec"
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

func (r RoundedRock) canMoveSouth(a []string) bool {
	if r.y == len(a)-1 {return false}
	if string(a[r.y+1][r.x]) == ROUNDED_ROCK {return false}
	if string(a[r.y+1][r.x]) == CUBE_ROCK {return false}
	return true
}

func (r RoundedRock) canMoveWest(a []string) bool {
	if r.x == 0 {return false}
	if string(a[r.y][r.x - 1]) == ROUNDED_ROCK {return false}
	if string(a[r.y][r.x - 1]) == CUBE_ROCK {return false}
	return true
}

func (r RoundedRock) canMoveEast(a []string) bool {
	if r.x == len(a)-1 {return false}
	if string(a[r.y][r.x + 1]) == ROUNDED_ROCK {return false}
	if string(a[r.y][r.x + 1]) == CUBE_ROCK {return false}
	return true
}

// func (r RoundedRock) rollNorth(a []string) []string {
// 	if !r.canMoveNorth(a) {
// 		return a
// 	}
// 	a[r.y] = a[r.y][:r.x] + EMPTY_SPACE + a[r.y][r.x+1:]
// 	a[r.y-1] = a[r.y-1][:r.x] + ROUNDED_ROCK + a[r.y-1][r.x+1:]
// 	r.y = r.y - 1
// 	return a
// }

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

func copyMap(data []string) []string {
    out := []string{}

    for i:=0; i <len(data); i++ {
        a := data[i]
        out = append(out, a)
    }
    return out
}

func printMap(data []string) {
	for i := 0; i<len(data);i++ {fmt.Printf("%s\n", data[i])}
}

func all(a []bool) bool {
    for _, i := range a {
        if i != true {
            return false
        }
    }
    return true
}

func any(a []bool) bool {
    for _, i := range a {
        if i == true {
            return true
        }
    }
    return false
}

func moveAllNorth(data []string, rocksMap map[int][]RoundedRock) ([]string, map[int][]RoundedRock) {
        d := copyMap(data)

	for i := 0; i<len(d); i++ {
                movableStatus := make([]bool, len(rocksMap[i]))
                for j := range movableStatus { movableStatus[j] = true }
		for any(movableStatus) {
			for j, r := range rocksMap[i] {
				if !r.canMoveNorth(d) {
                                    movableStatus[j] = false
                                } else {
                                    d[r.y] = d[r.y][:r.x] + EMPTY_SPACE + d[r.y][r.x+1:]
                                    d[r.y-1] = d[r.y-1][:r.x] + ROUNDED_ROCK + d[r.y-1][r.x+1:]
                                    r.y = r.y-1
                                    // update the map with the rock (and its new position)
                                    rocksMap[i][j] = r
                                }

				// cmd := exec.Command("clear")
				// cmd.Stdout = os.Stdout
				// cmd.Run()
                                // fmt.Printf("Row: %d | Movable: %t\n",i, movableStatus)
				// printMap(d)
				// time.Sleep(50 * time.Millisecond)
			}
		}
	}

	return d, rocksMap
}

func moveAllSouth(data []string, rocksMap map[int][]RoundedRock) ([]string, map[int][]RoundedRock) {
        d := copyMap(data)

	for i := len(d)-1; i>=0; i-- {
                movableStatus := make([]bool, len(rocksMap[i]))
                for j := range movableStatus { movableStatus[j] = true }
		for any(movableStatus) {
			for j, r := range rocksMap[i] {
				if !r.canMoveSouth(d) {
                                    movableStatus[j] = false
                                } else {
                                    d[r.y] = d[r.y][:r.x] + EMPTY_SPACE + d[r.y][r.x+1:]
                                    d[r.y+1] = d[r.y+1][:r.x] + ROUNDED_ROCK + d[r.y+1][r.x+1:]
                                    r.y = r.y+1
                                    // update the map with the rock (and its new position)
                                    rocksMap[i][j] = r
                                }
			}
		}
	}

	return d, rocksMap
}

func moveAllWest(data []string, rocksMap map[int][]RoundedRock) ([]string, map[int][]RoundedRock) {
        d := copyMap(data)
	for i := 0; i<len(d); i++ {
                movableStatus := make([]bool, len(rocksMap[i]))
                for j := range movableStatus { movableStatus[j] = true }
		for any(movableStatus) {
                    for j := 0; j < len(rocksMap[i]); j++ {
                                r := rocksMap[i][j]
				if !r.canMoveWest(d) {
                                    movableStatus[j] = false
                                } else {
                                    d[r.y] = d[r.y][:r.x-1] + ROUNDED_ROCK  + EMPTY_SPACE + d[r.y][r.x+1:]
                                    r.x = r.x-1
                                    // update the map with the rock (and its new position)
                                    rocksMap[i][j] = r
                                }
			}
		}
	}

	return d, rocksMap
}

func moveAllEast(data []string, rocksMap map[int][]RoundedRock) ([]string, map[int][]RoundedRock) {
        d := copyMap(data)
	for i := 0; i<len(d); i++ {
                movableStatus := make([]bool, len(rocksMap[i]))
                for j := range movableStatus { movableStatus[j] = true }
		for any(movableStatus) {
                    for j := len(rocksMap[i])-1; j >= 0; j-- {
                                r := rocksMap[i][j]
				if !r.canMoveEast(d) {
                                    movableStatus[j] = false
                                } else {
                                    d[r.y] = d[r.y][:r.x-1] +  EMPTY_SPACE + ROUNDED_ROCK  + d[r.y][r.x+1:]
                                    r.x = r.x+1
                                    // update the map with the rock (and its new position)
                                    rocksMap[i][j] = r
                                }
			}
		}
	}

	return d, rocksMap
}

func getRocksMap(data []string) map[int][]RoundedRock {
	rocks := getAllRoundedRock(data)
	rocksMap := make(map[int][]RoundedRock )
	for i := 0; i<len(data); i++ {rocksMap[i] = []RoundedRock{}}
	for _, i := range rocks { rocksMap[i.y] = append(rocksMap[i.y], i) }
        return rocksMap
}

func solution1(data []string) int {
	rocksMap := getRocksMap(data)
	data, rocksMap = moveAllNorth(data, rocksMap)

        out := 0
        mapSize := len(data)
        for _, rs := range rocksMap {
            for i := 0; i < len(rs); i++ {
                out += mapSize - rs[i].y
            }
        }
	return out
}

func solution2(data []string) int {
	rocksMap := getRocksMap(data)

        nCycles := 1000000000
        for i:=0; i<nCycles; i++ {
            fmt.Printf("Cycle: %d/%d\r", i, nCycles)
            data, rocksMap = moveAllNorth(data, rocksMap)
            rocksMap = getRocksMap(data)
            data, rocksMap = moveAllWest(data, rocksMap)
            rocksMap = getRocksMap(data)
            data, rocksMap = moveAllSouth(data, rocksMap)
            rocksMap = getRocksMap(data)
            data, rocksMap = moveAllEast(data, rocksMap)
            rocksMap = getRocksMap(data)
        }
        fmt.Print("\n")

        out := 0
        mapSize := len(data)
        for _, rs := range rocksMap {
            for i := 0; i < len(rs); i++ {
                out += mapSize - rs[i].y
            }
        }
	return out
}


func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

