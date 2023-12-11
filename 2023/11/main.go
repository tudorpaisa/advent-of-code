package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"regexp"
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

func getAllEmptyRowIndex(data []string) []int {
	out := []int{}
	re := regexp.MustCompile("\\.")
	for i, row := range data {
		result := re.FindAllIndex([]byte(row), len(row)+1)
		if len(result) == len(row) {
			out = append(out, i)
		}
	}
	return out
}

func isEmptyColumn(data []string, colIdx int) bool {
	for i := range data {
		if string(data[i][colIdx]) != "." {
			return false
		}
	}
	return true
}

func getAllEmptyColumnIndex(data []string) []int {
	out := []int{}

	for i := 0; i < len(data[0]); i++ {
		if isEmptyColumn(data, i) {
			out = append(out, i)
		}
	}

	return out
}

func expandRows(data []string, emptyRows []int, multiplier int) []string {
	increment := 0
	for x, rowIdx := range emptyRows {
		fmt.Printf("Expanding rows: %d/%d\r", x+1, len(emptyRows))
		for i := 0; i < multiplier - 1; i++ {

			data = append(data[:rowIdx + 1 + increment], data[rowIdx + increment:]...)

			increment = increment + 1

		}
	}
	fmt.Print("\n")
	return data
}

func expandColumns(data []string, emptyColumns []int, multiplier int) []string {
	increment := 0
	// yikes... i don't like this, but it's clear what's happening
	for x, colIdx := range emptyColumns {
		fmt.Printf("Expanding columns: %d/%d\r", x+1, len(emptyColumns))
		for i, row := range data {
			multipliedString := strings.Repeat(".", multiplier-1)
			data[i] = row[:colIdx+increment] + multipliedString + row[colIdx+increment:]
		}
		increment = increment + multiplier - 1
	}
	fmt.Print("\n")
	return data
}

func expandUniverse(data []string, multiplier int) []string {
	emptyRows := getAllEmptyRowIndex(data)
	emptyColumns := getAllEmptyColumnIndex(data)

	data = expandColumns(data, emptyColumns, multiplier)
	data = expandRows(data, emptyRows, multiplier)
	return data
}

func findUniverses(data []string) map[int][]int {
	out := make(map[int][]int)
	universeNumber := 0
	re := regexp.MustCompile("#")
	for y, row := range data {
		result := re.FindAllIndex([]byte(row), len(row)+1)
		if len(result) > 0 {
			for _, universeIdx := range result {
				out[universeNumber] = []int{universeIdx[0], y}
				universeNumber++
			}
		}
	}
	return out
}

func namePair(a int, b int) string {
	return strconv.Itoa(a) + " <-> " + strconv.Itoa(b)
}

func manhattanDist(a []int, b []int) int {
	return int(
		math.Abs(
			float64(
				a[0] - b[0])) +
		math.Abs(
			float64(
				a[1] - b[1])))
}

func calculateDistances(locations map[int][]int) map[string]int {
	dist := make(map[string]int)

	for a, coordsA := range locations {
		for b, coordsB := range locations {
			aBPairName := namePair(a, b)
			bAPairName := namePair(b, a)

			if _, ok := dist[aBPairName]; !ok {
				x := manhattanDist(coordsA, coordsB)
				dist[aBPairName] = x
				dist[bAPairName] = x
			}
		}
	}

	return dist
}

func findValuesInRange(a []int, low int, high int) []int {
	out := []int{}
	if low > high {
		low, high = high, low
	}

	for _, i := range a {
		if i > low && i < high {
			out = append(out, i)
		}
	}
	return out
}

func calculateDistancesWithExpansion(locations map[int][]int, emptyRows []int, emptyColumns []int, expansionMultiplier int) map[string]int {
	dist := make(map[string]int)

	for a, coordsA := range locations {
		for b, coordsB := range locations {

			aBPairName := namePair(a, b)
			bAPairName := namePair(b, a)

			if a == b {
				continue
			}

			spacesBetweenX := findValuesInRange(emptyColumns, coordsA[0], coordsB[0])
			spacesBetweenY := findValuesInRange(emptyRows, coordsA[1], coordsB[1])

			var xA int
			var yA int
			if coordsA[0] < coordsB[0] {
				xA = coordsA[0] - ( ( len(spacesBetweenX) * expansionMultiplier ) - len(spacesBetweenX) )
			} else {
				xA = coordsA[0] + ( ( len(spacesBetweenX) * expansionMultiplier ) - len(spacesBetweenX) )

			}
			if coordsA[1] < coordsB[1] {
				yA = coordsA[1] - ( ( len(spacesBetweenY) * expansionMultiplier ) - len(spacesBetweenY) )
			} else {
				yA = coordsA[1] + ( ( len(spacesBetweenY) * expansionMultiplier ) - len(spacesBetweenY) )

			}

			x := manhattanDist([]int{xA, yA}, coordsB)

			// sanity check because go gave me a few scares here
			dist[aBPairName] = x
			if val, ok := dist[bAPairName]; ok {
				if val != x {
					panic(fmt.Sprintf("VALUES DIFFER: %d - %d\n", val, x))
				}
			}
		}
	}

	return dist
}

func printUniverse(data []string) {
	for _, i := range data {
		fmt.Printf("%s\n", i)
	}
}

func solution1(data []string) int {
	data = expandUniverse(data, 2)
	universeLocations := findUniverses(data)
	distances := calculateDistances(universeLocations)

	out := 0
	for _, v := range distances {
		out += v
	}
	return out / 2
}

func solution2(data []string) int {
	emptyRows := getAllEmptyRowIndex(data)
	emptyColumns := getAllEmptyColumnIndex(data)
	universeLocations := findUniverses(data)
	distances := calculateDistancesWithExpansion(universeLocations, emptyRows, emptyColumns, 1000000)

	out := 0
	for _, v := range distances {
		out += v
	}
	return out / 2
}

func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	data = readLines("input1.txt")
	fmt.Printf("Solution 2: %d\n", solution2(data))
}
