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

type Direction int

const (
	North Direction = iota
	South
	West
	East
)

func initEnergyMap(data []string) [][]int {
	out := [][]int{}
	for i := 0; i < len(data); i++ {
		row := []int{}
		for j := 0; j < len(data[i]); j++ {
			row = append(row, 0)
		}
		out = append(out, row)
	}
	return out
}

type Beam struct {
	x int
	y int
	direction Direction
}

func incrementDirection(beam *Beam) {
	if beam.direction == North {
		beam.y -= 1
	} else if beam.direction == South {
		beam.y += 1
	} else if beam.direction == West {
		beam.x -= 1
	} else {
		beam.x += 1
	}
}

func withinBounds(beam *Beam, layout *[]string) bool {
	if beam.x < 0 || beam.x >=len((*layout)[0]) { return false }
	if beam.y < 0 || beam.y >=len(*layout) { return false }
	return true
}

func isSeen(a Beam, b []Beam) bool {
	for _, i := range b {
		if a.x == i.x && a.y == i.y && a.direction == i.direction {
			return true
		}
	}
	return false
}

func computeEnergized(data []string, startBeam Beam) [][]int {
	energyMap := initEnergyMap(data)
	energyMap[0][0] =  1
	beamQueue := []Beam{ startBeam }
	seen := []Beam{}
	for len(beamQueue) > 0 {
		currentBeam := beamQueue[0]
		// pop
		beamQueue = beamQueue[1:]
		if !withinBounds(&currentBeam, &data) {
			continue
		}
		if isSeen(currentBeam, seen) {
			continue
		}

		seen = append(seen, currentBeam)

		energyMap[currentBeam.y][currentBeam.x] = energyMap[currentBeam.y][currentBeam.x] + 1
		currentTile := data[currentBeam.y][currentBeam.x]

		if currentTile == '/' {
			if currentBeam.direction == North {
				beamQueue = append(beamQueue, Beam{currentBeam.x + 1, currentBeam.y, East})

			} else if currentBeam.direction == South {
				beamQueue = append(beamQueue, Beam{currentBeam.x - 1, currentBeam.y, West})

			} else if currentBeam.direction == West {
				beamQueue = append(beamQueue, Beam{currentBeam.x, currentBeam.y + 1, South})

			} else {
				beamQueue = append(beamQueue, Beam{currentBeam.x, currentBeam.y - 1, North})

			}
		} else if currentTile == '\\' {
			if currentBeam.direction == North {
				beamQueue = append(beamQueue, Beam{currentBeam.x - 1, currentBeam.y, West})

			} else if currentBeam.direction == South {
				beamQueue = append(beamQueue, Beam{currentBeam.x + 1, currentBeam.y, East})

			} else if currentBeam.direction == West {
				beamQueue = append(beamQueue, Beam{currentBeam.x, currentBeam.y - 1, North})

			} else {
				beamQueue = append(beamQueue, Beam{currentBeam.x, currentBeam.y + 1, South})

			}
		} else if currentTile == '|' {
			if currentBeam.direction == West || currentBeam.direction == East {
				beam1 := Beam{ currentBeam.x, currentBeam.y - 1, North }
				beam2 := Beam{ currentBeam.x, currentBeam.y + 1, South }
				beamQueue = append(beamQueue, beam1)
				beamQueue = append(beamQueue, beam2)
			} else {
				incrementDirection(&currentBeam)
				beamQueue = append(beamQueue, currentBeam)
			}

		} else if currentTile == '-' {
			if currentBeam.direction == North || currentBeam.direction == South {
				beam1 := Beam{ currentBeam.x + 1, currentBeam.y, East }
				beam2 := Beam{ currentBeam.x - 1, currentBeam.y, West }
				beamQueue = append(beamQueue, beam1)
				beamQueue = append(beamQueue, beam2)
			} else {
				incrementDirection(&currentBeam)
				beamQueue = append(beamQueue, currentBeam)
			}
		} else {
			incrementDirection(&currentBeam)
			beamQueue = append(beamQueue, currentBeam)
		}
	}

	return energyMap
}

func countEnergised(m [][]int) int {
	out := 0
	for _, i := range m {
		for _, j := range i {
			if j > 0 {
				out++
			}
		}
	}
	return out
}

func printEnergyMap(m [][]int) {
	fmt.Print("\n")
	for _, i := range m {
		for _, j := range i {
			if j == 0 {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Print("\n")
	}
}

func sumMaps(a [][]int, b [][]int) [][]int {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			a[i][j] = a[i][j] + b[i][j]
		}
	}
	return a
}

func solution1(data []string) int {
	start := Beam{0, 0, East}
	energyMap := computeEnergized(data, start)
	return countEnergised(energyMap)
}

func solution2(data []string) int {
	startPositions := []Beam{}

	for i := 0; i<len(data); i++ {
		startPositions = append(startPositions, Beam{0, i, East})
		startPositions = append(startPositions, Beam{len(data[0])-1, i, West})
	}
	for i := 1; i<len(data[0])-1; i++ {
		startPositions = append(startPositions, Beam{i, 0, South})
		startPositions = append(startPositions, Beam{i, len(data)-1, North})
	}

	max := 0
	for idx, i := range startPositions {
		fmt.Printf("%d/%d\r", idx, len(startPositions))
		res := computeEnergized(data, i)
		count := countEnergised(res)
		if count > max {
			max = count
		}
	}

	return max
}


func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

