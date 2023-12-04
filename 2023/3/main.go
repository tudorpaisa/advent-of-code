package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"regexp"
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

func findAllIndexByRegex(data []string, expression string) map[int][][]int {
	re, _ := regexp.Compile(expression)

	out := make(map[int][][]int)

	for i, row := range data {
		match := re.FindAllIndex([]byte(row), 100)
		if match != nil {
			out[i] = match
		}
	}

	return out
}

func regexMatches(str string, query string) bool {
	re, _ := regexp.Compile(query)
	return re.Match([]byte(str))
}

func extractParts(data []string, dataIndices map[int][][]int) []Part {
	var out []Part = *new([]Part)

	for i, rowIndices := range dataIndices {
		fmt.Printf("%s\n", data[i])
		for _, partIndices := range rowIndices {
			pi := []int{ partIndices[0], partIndices[1] }
			numberString := data[i][pi[0]:pi[1]]
			number, _ := strconv.Atoi(numberString)

			partNumber := initPart(i, number, pi)
			out = append(out, partNumber)
		}
	}

	return out
}

type Part struct {
	rowNumber int
	partNumber int
	index []int
}

func (p Part) hasSymbols(data []string) bool {
	return p.hasSymbolLeft(data) ||
		p.hasSymbolRight(data) ||
		p.hasSymbolTop(data) ||
		p.hasSymbolBelow(data) ||
		p.hasSymbolTopLeft(data) ||
		p.hasSymbolTopRight(data) ||
		p.hasSymbolBottomLeft(data) ||
		p.hasSymbolBottomRight(data)
}

func (p Part) hasSymbolTop(data []string) bool {
	// sanity check
	if p.rowNumber - 1 < 0 {
		return false
	}

	rowAbove := data[p.rowNumber - 1]

	targetString := rowAbove[p.index[0]:p.index[1]]

	if regexMatches(targetString, "[^0-9.]") {
		// fmt.Printf("Found match (%s) on row %d, part %d\n", targetString, p.rowNumber, p.partNumber)
		return true
	}

	return false
}

func (p Part) hasSymbolBelow(data []string) bool {
	// sanity check
	if p.rowNumber + 1 >= len(data) {
		return false
	}

	rowBelow := data[p.rowNumber + 1]

	targetString := rowBelow[p.index[0]:p.index[1]]

	if regexMatches(targetString, "[^0-9.]") {
		// fmt.Printf("Found match (%s) on row %d, part %d\n", targetString, p.rowNumber, p.partNumber)
		return true
	}

	return false
}

func (p Part) hasSymbolLeft(data []string) bool {
	x := p.index[0] - 1
	y := p.rowNumber

	// sanity check
	if x < 0 {
		return false
	}

	targetChar := string(data[y][x])
	if regexMatches(targetChar, "[^0-9.]") {
		// fmt.Printf("Found match (%s) on row %d, part %d\n", targetChar, p.rowNumber, p.partNumber)
		return true
	}
	return false
}

func (p Part) hasSymbolRight(data []string) bool {
	x := p.index[1]
	y := p.rowNumber

	// sanity check
	if x >= len(data[y]) {
		return false
	}

	targetChar := string(data[y][x])
	// if regexMatches(targetChar, "\\w|\\.") {
	// 	return false
	// }
	if regexMatches(targetChar, "[^0-9.]") {
		// fmt.Printf("Found match (%s) on row %d, part %d\n", targetChar, p.rowNumber, p.partNumber)
		return true
	}
	return false
}

func (p Part) hasSymbolTopLeft(data []string) bool {
	y := p.rowNumber - 1
	x := p.index[0] - 1

	// sanity check
	if y < 0 {
		return false
	}
	if x < 0 {
		return false
	}

	targetChar := string(data[y][x])
	if regexMatches(targetChar, "[^0-9.]") {
		// fmt.Printf("Found match (%s) on row %d, part %d\n", targetChar, p.rowNumber, p.partNumber)
		return true
	}
	return false

}

func (p Part) hasSymbolTopRight(data []string) bool {
	y := p.rowNumber - 1
	x := p.index[1]

	// sanity check
	if y < 0 {
		return false
	}
	if x >= len(data[y]) {
		return false
	}

	targetChar := string(data[y][x])
	if regexMatches(targetChar, "[^0-9.]") {
		// fmt.Printf("Found match (%s) on row %d, part %d\n", targetChar, p.rowNumber, p.partNumber)
		return true
	}
	return false

}

func (p Part) hasSymbolBottomLeft(data []string) bool {
	y := p.rowNumber + 1
	x := p.index[0] - 1

	// sanity check
	if y >= len(data) {
		return false
	}
	if x < 0 {
		return false
	}

	targetChar := string(data[y][x])
	if regexMatches(targetChar, "[^0-9.]") {
		// fmt.Printf("Found match (%s) on row %d, part %d\n", targetChar, p.rowNumber, p.partNumber)
		return true
	}
	return false

}

func (p Part) isEndOfRow(data []string) bool {
	x := p.index[1]
	return x == len(data[p.rowNumber])
}

func (p Part) hasSymbolBottomRight(data []string) bool {
	y := p.rowNumber + 1
	x := p.index[1]

	// sanity check
	if y >= len(data) {
		return false
	}
	if x >= len(data[y]) {
		return false
	}

	targetChar := string(data[y][x])
	if regexMatches(targetChar, "[^0-9.]") {
		// fmt.Printf("Found match (%s) on row %d, part %d\n", targetChar, p.rowNumber, p.partNumber)
		return true
	}
	return false

}

func initPart(rowNumber int, partNumber int, index []int) Part {
	out := Part{}
	out.rowNumber = rowNumber
	out.partNumber = partNumber
	out.index = index
	return out
}

func printParts(parts []Part) {
	for _, part := range parts {
		fmt.Printf("Line: %d\tNumber: %d\n", part.rowNumber, part.partNumber)
	}
}

func sum(a []int) int {
	out := 0
	for _, i := range a {
		out += i
	}

	return out
}

func solution1(data []string) int {
	matchIndices := findAllIndexByRegex(data, "(\\d+)")
	parts := extractParts(data, matchIndices)

	// printParts(parts)

	verifiedParts := *new([]int)
	for _, i := range parts {
		if i.hasSymbols(data) {
			verifiedParts = append(verifiedParts, i.partNumber)
		}
	}

	return sum(verifiedParts)
}

func numberWithinRange(number int, low int, high int) bool {
	return number >=low && number <= high
}

func indexStartsAt(index int, low int) bool {
	return index == low
}

func indexEndsAt(index int, high int) bool {
	return index == high
}

func getNumberSameRow(index int, partIndices [][]int, row string, indexBuffer int) ( int, bool ) {
	for _, partIdx := range partIndices {
		if numberWithinRange(index, partIdx[0] + indexBuffer, partIdx[1] + indexBuffer) {
			out, _ := strconv.Atoi(row[partIdx[0]:partIdx[1]])
			return out, true
		}
	}
	return 0, false
}

func getNumberStartsAt(index int, partIndices [][]int, row string) (int, bool) {
	for _, partIdx := range partIndices {
		if indexStartsAt(index, partIdx[0]) {
			out, _ := strconv.Atoi(row[partIdx[0]:partIdx[1]])
			return out, true
		}
	}
	return 0, false
}

func getNumberEndsAt(index int, partIndices [][]int, row string) (int, bool) {
	for _, partIdx := range partIndices {
		if indexEndsAt(index, partIdx[1]) {
			out, _ := strconv.Atoi(row[partIdx[0]:partIdx[1]])
			return out, true
		}
	}
	return 0, false
}

func solution2(data []string) int {
	starIndices := findAllIndexByRegex(data, "\\*")
	partIndices := findAllIndexByRegex(data, "(\\d+)")
	fmt.Println(partIndices)

	products := *new([]int)

	for rowNumber, starsRow := range starIndices {
		for _, starIndex := range starsRow {
			foundNumbers := *new([]int)
			var number int;
			var foundNumber bool;
			// check left & right
			if partIndices, ok := partIndices[rowNumber]; ok {
				// number left
				number, foundNumber = getNumberSameRow(starIndex[0], partIndices, data[rowNumber], 0)
				if foundNumber {
					foundNumbers = append(foundNumbers, number)
				}
				// number right
				number, foundNumber = getNumberSameRow(starIndex[1], partIndices, data[rowNumber], 0)
				if foundNumber {
					foundNumbers = append(foundNumbers, number)
				}
			}

			// check top
			if partIndices, ok := partIndices[rowNumber - 1]; ok {
				// right above
				number, foundNumber = getNumberSameRow(starIndex[0], partIndices, data[rowNumber - 1], 0)
				if foundNumber {
					foundNumbers = append(foundNumbers, number)
				}
				// top right
				number, foundNumber = getNumberStartsAt(starIndex[1], partIndices, data[rowNumber - 1])
				if foundNumber {
					foundNumbers = append(foundNumbers, number)
				}
			}

			// check bottom
			if partIndices, ok := partIndices[rowNumber + 1]; ok {
				number, foundNumber = getNumberSameRow(starIndex[0], partIndices, data[rowNumber + 1], 0)
				if foundNumber {
					foundNumbers = append(foundNumbers, number)
				}
				// top right
				number, foundNumber = getNumberStartsAt(starIndex[1], partIndices, data[rowNumber + 1])
				if foundNumber {
					foundNumbers = append(foundNumbers, number)
				}
			}

			fmt.Printf("%d: %d\n", rowNumber, foundNumbers)
			if len(foundNumbers) == 2 {
				products = append(products, foundNumbers[0]*foundNumbers[1])
			}
		}
	}

	return sum(products)
}

func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}
