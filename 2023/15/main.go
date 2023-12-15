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

func preprocess(data []string) []string {
	out := []string{}
	for _, i := range data {
		out = append(out, strings.Split(i, ",")...)
	}
	return out
}

func hashString(a string) int {
	out := 0
	for i := 0; i < len(a); i++ {
		out += int(a[i])
		out *= 17
		out = out % 256
	}
	return out
}

func solution1(data []string) int {
	data = preprocess(data)

	out := 0
	for i := 0; i < len(data); i++ {
		out += hashString(data[i])
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

