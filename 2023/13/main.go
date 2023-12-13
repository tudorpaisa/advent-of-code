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

func readLines(fname string) [][]string {
	var patterns []string = strings.Split(readFileRaw(fname), "\n\n")
	var out [][]string = *new([][]string)

	for _, pattern := range patterns {
		if pattern == "" {
			continue
		}
		res := []string{}
		split := strings.Split(pattern, "\n")
		for _, row := range split {
			if row == "" {
				continue
			}
			res = append(res, row)
		}
		out = append(out, res)
	}
	return out
}

func transpose(a []string) []string {
	out := []string{}

	// i == colIdx
	for i := 0; i < len(a[0]); i++ {
		transposedColumn := ""
		// j == rowIdx
		for j := 0; j < len(a); j++ {
			transposedColumn += string(a[j][i])
		}
		out = append(out, transposedColumn)
	}
	return out
}

func scanMirror2(data []string, tol int) []int {
	// this counts ALL mirror lengths
	// if tol == 0 we'll get just 1 result
	// if tol == 1 we'll get 2
	// etc...
	bestIndices := []int{}
	for i := 1; i < len(data); i++ {

		count := 0
		iterIndex := i
		mirrorIndex := i-1
		for mirrorIndex >=0 && iterIndex < len(data) {
			nDiff := countDiff(data[iterIndex], data[mirrorIndex])
			if nDiff <= tol {
				count++
			} else {
				count = -1
				break
			}
			iterIndex++
			mirrorIndex--
		}
		if count > 0 {
			bestIndices = append(bestIndices, i)
		}
	}
	return bestIndices
}

func scanMirror(data []string, tol int) (int, int) {
	max := 0
	bestIdx := 0
	for i := 1; i < len(data); i++ {

		count := 0
		iterIndex := i
		mirrorIndex := i-1
		for mirrorIndex >=0 && iterIndex < len(data) {
			nDiff := countDiff(data[iterIndex], data[mirrorIndex])
			if nDiff <= tol {
				count++
			} else {
				count = -1
				break
			}
			iterIndex++
			mirrorIndex--
		}
		if count > max {
			max = count
			bestIdx = i
		}
	}
	return bestIdx, max
}

func countDiff(a string, b string) int {
	count := 0
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] { count++ }
	}
	return count
}

func findSmudged(a []string) (int, int) {
	for i := 0; i < len(a)-1; i++ {
		for j := i+1; j < len(a); j++ {
			count := countDiff(a[i], a[j])

			if count == 1 {
				return j, i
			}
		}
	}

	return 0, 0
}

func solution1(data [][]string) int {
	out := 0
	for _, i := range data {
		columnsTransposed := transpose(i)
		rowsBest := scanMirror2(i, 0)
		columnsBest := scanMirror2(columnsTransposed, 0)
		for _, j := range columnsBest {
			out += j
		}
		for _, j := range rowsBest {
			out += j *100
		}
		// rowsBest, _ := scanMirror(i, 0)
		// columnsBest, _ := scanMirror(columnsTransposed, 0)

		// out += rowsBest * 100
		// out += columnsBest
	}
	return out
}

func solution2(data [][]string) int {
	out := 0
	for _, i := range data {
		columnsTransposed := transpose(i)
		rowsBest := scanMirror2(i, 1)
		columnsBest := scanMirror2(columnsTransposed, 1)
		for _, j := range columnsBest {
			out += j
		}
		for _, j := range rowsBest {
			out += j *100
		}
	}

	return out
}


func main() {
	data := readLines("input1.txt")
	s1 := solution1(data)
	fmt.Printf("Solution 1: %d\n", s1)
	data = readLines("input1.txt")
	fmt.Printf("Solution 2: %d\n", solution2(data)-s1)
}

