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

func preproces(data []string) map[string][]string {
	out := make(map[string][]string)
	for _, i := range data {
		s := strings.Replace(i, ":", "", 1)
		sp := strings.Split(s, " ")
		out[sp[0]] = sp[1:]
	}
	return out
}

func format(m map[string][]string) []string {
	out := []string{}
	for k, v := range m {
		for _, i := range v {
			out = append(out, fmt.Sprintf("    %s -> %s;\n", k, i))
		}
	}
	return out
}

func delete(v string, a []string) []string {
	out := []string{}
	for _, i := range a {
		if i != v { out = append(out, i) }
	}
	return out
}

func deleteConnection(m map[string][]string, from, to string) map[string][]string {
	m[from] = delete(to, m[from])
	return m
}

func toBidirectional(m map[string][]string) map[string][]string {
	for k, v := range m {
		for _, i := range v {
			if l, ok := m[i]; ok {
				m[i] = append(l, k)
			} else {
				l = []string{k}
				m[i] = l
			}
		}
	}
	return m
}

func countNodes(m map[string][]string, start string) int {
	count := 0

	queue := []string{start}
	visited := make(map[string]bool)

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		// double-counting. whoops xD
		if vst, ok := visited[node]; !ok {
			queue = append(queue, node)
		} else {
			if vst == true {
				continue
			}
		}

		count++

		visited[node] = true

		targets := m[node]
		for _, i := range targets {
			if vst, ok := visited[i]; !ok {
				queue = append(queue, i)
			} else {
				if vst == true {
					continue
				}
			}
		}
	}

	return count
}

func solution1(data []string) int {
	d := preproces(data)
	s := format(d)

	f, err := os.Create("graph.dot")
	if err != nil {panic(err)}

	f.WriteString("digraph {\n")
	for _, i := range s {
		f.WriteString(i)
	}
	f.WriteString("}")

	f.Close()

	d = deleteConnection(d, "mxd", "glz")
	d = deleteConnection(d, "brd", "clb")
	d = deleteConnection(d, "jxd", "bbz")

	d = toBidirectional(d)

	//723 717 -> 518391
	fmt.Println(countNodes(d, "brd"), countNodes(d, "jxd"))
	return countNodes(d, "brd") * countNodes(d, "jxd")

	// fmt.Println(d)

	// return 0
}

func solution2(data []string) int {
	return 0
}


func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

