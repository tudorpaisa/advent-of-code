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

func generateArrangement(pattern string, index int) []string {
    if index == len(pattern) {
        return []string{pattern}
    }

    out := []string{}
    if string(pattern[index]) == "?" {
        out = append(out, generateArrangement(pattern[:index] + "." + pattern[index+1:], index+1)...)
        out = append(out, generateArrangement(pattern[:index] + "#" + pattern[index+1:], index+1)...)
    } else {
        out = append(out, generateArrangement(pattern, index + 1)...)
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
            if err != nil { panic(err) }
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

func countBroken(a []string) [][]int {
    out := [][]int{}
    
    for i := 0; i < len(a); i++ {
        // fmt.Printf("%d, %s\n", i, a[i])
        counts := []int{}
        count := 0
        for j := 0; j < len(a[i]); j++ {
            if string(a[i][j]) == "#" {
                count++
            } else if string(a[i][j]) == "." && count > 0 {
                counts = append(counts, count)
                count = 0
            }
        }
        if count > 0 {
            counts = append(counts, count)
        }
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
        if allOk == true { out++ }
    }

    return out
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

func solution2(data []string) int {
	return 0
}


func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

