package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
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

func parseHistory(data []string) [][]int {
	out := *new([][]int)

	for _, i := range data {
		rowSplit := strings.Split(i, " ")
		rowInt := *new([]int)
		for _, j := range rowSplit {
			k, err := strconv.Atoi(j)
			if err != nil { panic(err) }
			rowInt = append(rowInt, k)
		}
		out = append(out, rowInt)
	}

	return out
}

func allZeroes(a []int) bool {
	for _, i := range a {
		if i != 0 { return false }
	}
	return true
}

func computeSequenceDiffs(a []int) [][]int {
	out := *new([][]int)
	out = append(out, a)

	for {
		interimArr := *new([]int)

		for i := 1; i < len(out[len(out) - 1]); i++ {
			interimArr = append(interimArr, out[len(out) - 1][i] - out[len(out) - 1][i - 1])
		}

		out = append(out, interimArr)
		if allZeroes(interimArr) { break }
	}

	return out
}

func padSeqDiff(a [][]int, diff bool) [][]int {
	lenArr := len(a)

	for i := (lenArr-1); i >= 0; i-- {
		if i == len(a)-1 {
			a[i] = append(a[i], 0)
			continue
		}
		var val int
		if diff {
			val = a[i][len(a[i])-1] - a[i+1][len(a[i+1])-1]
		} else {
			val = a[i][len(a[i])-1] + a[i+1][len(a[i+1])-1]
		}
		a[i] = append(a[i], val)
	}

	return a
}

func sum(a []int) int {
	out := 0
	for _, i := range a {
		out += i
	}
	return out
}

func reverseArrArr(a [][]int) [][]int {
	out := *new([][]int)

	for _, i := range a {
		iRev := reverseArr(i)
		out = append(out, iRev)
	}
	return out
}

func reverseArr(a []int) []int {
	out := *new([]int)
	for i := len(a)-1; i >= 0; i-- {
		out = append(out, a[i])
	}
	return out
}

func solution1(data []string) int {
	history := parseHistory(data)

	toSum := *new([]int)

	for _, i := range history {
		seqDiff := computeSequenceDiffs(i)
		seqDiff = padSeqDiff(seqDiff, false)
		x := seqDiff[0][len(seqDiff[0]) - 1]
		toSum = append(toSum, x)
	}

	return sum(toSum)
}

func solution2(data []string) int {
	history := parseHistory(data)

	toSum := *new([]int)

	for _, i := range history {
		seqDiff := computeSequenceDiffs(i)
		seqDiff = reverseArrArr(seqDiff)
		seqDiff = padSeqDiff(seqDiff, true)
		x := seqDiff[0][len(seqDiff[0]) - 1]
		toSum = append(toSum, x)
	}

	return sum(toSum)
}

func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

