package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"regexp"
	"strconv"
)


const digitsRegex = "1|2|3|4|5|6|7|8|9"
const spelledNumbersRegex = "one|two|three|four|five|six|seven|eight|nine"

func isDigit(char rune) bool {
	match, _ := regexp.Match(digitsRegex, []byte(string(char)))
	return match
}

func getSpelledNumber(s string) (string, bool) {
	re, _ := regexp.CompilePOSIX(spelledNumbersRegex)
	match := re.Find([]byte(s))
	if match != nil {
		return string(match), true
	}
	return "", false
}

func spelledNumberToDigit(s string) string {
	switch s {
	case "one":
		return "1"
	case "two":
		return "2"
	case "three":
		return "3"
	case "four":
		return "4"
	case "five":
		return "5"
	case "six":
		return "6"
	case "seven":
		return "7"
	case "eight":
		return "8"
	default:
		return "9"
	}
}

func readFileRaw(fname string) string {
	file, err := os.ReadFile(fname)

	if err != nil {
		msg := fmt.Sprintf("Encountered an error while reading '%s': %s", fname, err)
		log.Fatal(msg)
	}

	return string(file)
}

func readLines(fname string) []string {
	var data []string = strings.Split(readFileRaw("input1.txt"), "\n")
	var out []string = *new([]string)

	for _, row := range data {
		if row == "" {
			continue
		}
		out = append(out, row)
	}
	return out
}


func sum(a []int) int {
	var out int
	for _, i := range a {
		out += i
	}
	return out
}

func solution1(data []string) int {
	filtered_digits := *new([][]string)

	for _, row := range data {
		digits := *new([]string)

		for _, char := range row {
			if isDigit(char) {
				digits = append(digits, string(char))
			}
		}

		filtered_digits = append(filtered_digits, digits)
	}

	numbers := *new([]int)
	for _, i := range filtered_digits {
		var row_number string = i[0] + i[len(i)-1]
		n, _ := strconv.Atoi(row_number)
		numbers = append(numbers, n)
	}

	return sum(numbers)
}

func solution2(data []string) int {
	filtered_digits := *new([][]string)

	for _, row := range data {
		digits := *new([]string)
		collected_chars := ""

		for _, char := range row {
			collected_chars = collected_chars + string(char)

			if isDigit(char) {
				digits = append(digits, string(char))
				collected_chars = ""
			} else if number, ok := getSpelledNumber(collected_chars); ok {
				digits = append(digits, spelledNumberToDigit(number))
				collected_chars = string(collected_chars[len(collected_chars)-1])
			}
		}

		filtered_digits = append(filtered_digits, digits)
	}

	numbers := *new([]int)
	for _, i := range filtered_digits {
		var row_number string = i[0] + i[len(i)-1]
		n, _ := strconv.Atoi(row_number)
		numbers = append(numbers, n)
	}

	return sum(numbers)
}

func main() {
	data := readLines("input1.txt")

	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

