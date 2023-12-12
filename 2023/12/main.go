package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
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

func generateArrangement(pattern string, index int) []string {
	if index == len(pattern) {
		return []string{pattern}
	}

	out := []string{}
	if string(pattern[index]) == "?" {
		out = append(out, generateArrangement(pattern[:index]+"."+pattern[index+1:], index+1)...)
		out = append(out, generateArrangement(pattern[:index]+"#"+pattern[index+1:], index+1)...)
	} else {
		out = append(out, generateArrangement(pattern, index+1)...)
	}
	return out
}

type Arrangement struct {
	pattern string
	damages []int
}

func preprocess(data []string) []Arrangement {
	out := []Arrangement{}

	for _, i := range data {
		s := strings.Split(i, " ")
		pattern := s[0]
		damages := []int{}
		for _, j := range strings.Split(s[1], ",") {
			res, err := strconv.Atoi(j)
			if err != nil {
				panic(err)
			}
			damages = append(damages, res)
		}

		out = append(out, Arrangement{pattern, damages})
	}
	return out
}

func printArrangements(arr []Arrangement) {
	for _, i := range arr {
		fmt.Printf("Pattern=%20s Damages=%d\n", i.pattern, i.damages)
	}
}

func getBrokenCountSingle(a string) []int {
	counts := []int{}
	count := 0
	for i := 0; i < len(a); i++ {
		if string(a[i]) == "#" {
			count++
		} else if string(a[i]) == "." && count > 0 {
			counts = append(counts, count)
			count = 0
		}
	}
	if count > 0 {
		counts = append(counts, count)
	}
	return counts
}

func countBroken(a []string) [][]int {
	out := [][]int{}

	for i := 0; i < len(a); i++ {
		counts := getBrokenCountSingle(a[i])
		out = append(out, counts)
	}

	return out
}

func findBrokenMatches(broken [][]int, damages []int) int {
	out := 0

	for _, set := range broken {
		if len(set) != len(damages) {
			continue
		}
		allOk := true

		for i, j := range set {
			if j != damages[i] {
				allOk = false
				break
			}
		}
		if allOk == true {
			out++
		}
	}

	return out
}

func countArrangements(patternIndex int, damagesIndex int, pattern string, damages []int, cache [][]int) int {
	if patternIndex >= len(pattern) {
		if damagesIndex < len(damages) {
			return 0
		}
		return 1
	}

	if cache[patternIndex][damagesIndex] != -1 {
		return cache[patternIndex][damagesIndex]
	}

	result := 0
	if pattern[patternIndex] == '.' {
		result = countArrangements(patternIndex + 1, damagesIndex, pattern, damages, cache)
	} else {
		if pattern[patternIndex] == '?' {
			result += countArrangements(patternIndex + 1, damagesIndex, pattern, damages, cache)
		}
		if damagesIndex < len(damages) {
			count := 0
			for i := patternIndex; i < len(pattern); i++ {
				if pattern[i] == '.' || count > damages[damagesIndex] || count == damages[damagesIndex] && pattern[i] == '?' {
					break
				}
				count++
			}

			if count == damages[damagesIndex] {
				if count + patternIndex < len(pattern) && pattern[count + patternIndex] != '#' {
					result += countArrangements(patternIndex + count + 1, damagesIndex + 1, pattern, damages, cache)
				} else {
					result += countArrangements(patternIndex + count, damagesIndex + 1, pattern, damages, cache)
				}
			}
		}
	}

	cache[patternIndex][damagesIndex] = result

	return result
}

func solution1(data []string) int {
	arr := preprocess(data)
	// printArrangements(arr)
	out := 0
	for _, i := range arr {
		gen := generateArrangement(i.pattern, 0)
		broken := countBroken(gen)
		out += findBrokenMatches(broken, i.damages)
	}
	return out
}


func unfoldPattern(a string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out = out + a
		if i < n-1 {
			out = out + "?"
		}
	}
	return out
}

func unfoldDamages(a []int, n int) []int {
	out := []int{}
	for i := 0; i<n; i++ {
		out = append(out, a...)
	}
	return out
}

func solution2(data []string) int {
	arr := preprocess(data)

	out := 0
	for _, i := range arr {
		i.pattern = unfoldPattern(i.pattern, 5)
		i.damages = unfoldDamages(i.damages, 5)
		cache := [][]int{}
		for x := 0; x < len(i.pattern); x++ {
			d := make([]int, len(i.damages)+1)
			for y := 0; y < len(i.damages) +1; y++ {
				d[y] = -1
			}
			cache = append(cache, d)
		}

		out += countArrangements(0, 0, i.pattern, i.damages, cache)

	}
	return out
}

func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	data = readLines("input1.txt")
	fmt.Printf("Solution 2: %d\n", solution2(data))
}
