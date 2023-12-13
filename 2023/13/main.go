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

func scanMirror(data []string) (int, int) {
	max := -1
	bestIdx := -1
	for i := 1; i < len(data); i++ {

		count := 0
		iterIndex := i
		mirrorIndex := i-1
		for mirrorIndex >=0 && iterIndex < len(data) {
			if data[iterIndex] == data[mirrorIndex] {
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

func solution1(data [][]string) int {
	out := 0
	for idx, i := range data {
		columnsTransposed := transpose(i)
		rowsBest, rMax := scanMirror(i)
		columnsBest, cMax := scanMirror(columnsTransposed)
		if columnsBest == -1 {
			out += rowsBest * 100
		} else if rowsBest == -1 {

			out += columnsBest
		} else if (columnsBest - cMax == 0 && rowsBest - rMax == 0) || (columnsBest + cMax == len(columnsTransposed) && rowsBest + rMax == len(i)) {
			fmt.Printf("Found symmetrical mirroring at index: %d\nCols=%d (%d), Rows=%d (%d)\nProgram will panic.\n", idx, columnsBest, cMax, rowsBest, rMax)
			panic("SYMMETRICAL?!")
		} else if columnsBest - cMax == 0 || columnsBest + cMax == len(columnsTransposed) {
			out += columnsBest
		} else {
			out += rowsBest * 100
		}
	}
	return out
}

func solution2(data [][]string) int {
	return 0
}


func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

