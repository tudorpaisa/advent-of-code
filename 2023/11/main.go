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
	for _, rowIdx := range emptyRows {

		fmt.Printf("data:\n%s\n", data)
		var secondHalf []string
		var firstHalf []string
		copy(secondHalf, data[rowIdx + 1 + increment:])
		copy(firstHalf, data[:rowIdx + increment])

		for i := 0; i < multiplier; i++ {
			firstHalf = append(firstHalf, data[rowIdx])
		}

		fmt.Printf("first half:\n%s\n", firstHalf)
		fmt.Printf("second half:\n%s\n", secondHalf)
		data = append(firstHalf, secondHalf...)

		// data = append(data[:rowIdx + 1 + increment], data[rowIdx + increment:]...)

		increment = increment + multiplier - 1
	}
	printUniverse(data)
	return data
}

func expandColumns(data []string, emptyColumns []int) []string {
	increment := 0
	// ooof... i don't like this, but it's clear what's happening
	for _, colIdx := range emptyColumns {
		for i, row := range data {
			data[i] = row[:colIdx+increment] + string(row[colIdx+increment]) + row[colIdx+increment:]
		}
		increment++
	}
	// printUniverse(data)
	return data
}

func expandUniverse(data []string) []string {
	emptyRows := getAllEmptyRowIndex(data)
	emptyColumns := getAllEmptyColumnIndex(data)

	// fmt.Printf("Empty rows: %d\n", emptyRows)
	// fmt.Printf("Empty columns: %d\n", emptyColumns)
	data = expandRows(data, emptyRows, 2)
	data = expandColumns(data, emptyColumns)
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

func printUniverse(data []string) {
	for _, i := range data {
		fmt.Printf("%s\n", i)
	}
}

func solution1(data []string) int {
	data = expandUniverse(data)
	universeLocations := findUniverses(data)
	distances := calculateDistances(universeLocations)

	out := 0
	for _, v := range distances {
		out += v
	}
	return out / 2
}

func solution2(data []string) int {
	return 0
}


func main() {
	data := readLines("input2.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

