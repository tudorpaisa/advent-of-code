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

type Lens struct {
	label string
	slot int
	focalLength int
}

func printLensConf(a [][]Lens) {
	for i, lst := range a {
		if len(lst) == 0 {
			continue
		}
		fmt.Printf("Box %d:\n", i)
		for _, l := range lst {
			fmt.Printf("    %s = %d\n", l.label, l.focalLength)
		}
	}
	fmt.Print("\n")
}

func searchLabel(a string, l []Lens) (int, bool) {
	for i, j := range l {
		if a == j.label { return i, true}
	}
	return -1, false
}

func configureLenses(data []string) [][]Lens {
	m := make([][]Lens, 256)
	for i := 0; i < len(m); i++ { m[i] = []Lens{} }

	for _, i := range data {
		label := ""

		for j:=0; j<len(i); j++ {
			if i[j] == '=' {
				strFP := i[j+1:]
				focalLen, err := strconv.Atoi(strFP)
				if err != nil {panic(err)}

				hashCode := hashString(label)

				existingIdx, ok := searchLabel(label, m[hashCode])
				if ok {
					m[hashCode][existingIdx] = Lens{label, len(m[hashCode]), focalLen}
				} else {
					m[hashCode] = append(m[hashCode], Lens{label, len(m[hashCode]), focalLen})
				}

				break
			} else if i[j] == '-' {
				hashCode := hashString(label)
				for a, b := range m[hashCode] {
					if b.label == label {
						m[hashCode] = append(m[hashCode][:a], m[hashCode][a+1:]...)
					}
				}
			} else {
				label += string(i[j])
			}
		}
	}

	return m
}

func solution2(data []string) int {
	data = preprocess(data)
	m := configureLenses(data)

	out := 0
	for boxNum := range m {
		for i, lens := range m[boxNum] {
			out += (boxNum + 1) * (i+1) * lens.focalLength
		}
	}

	return out
}


func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

