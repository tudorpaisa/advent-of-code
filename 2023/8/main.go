package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"regexp"
	"os/exec"
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

func parseNodes(data []string) map[string][]string {
	out := make(map[string][]string)

	for _, i := range data {
		i = strings.ReplaceAll(i, " ", "")
		iSplit := strings.Split(i, "=")
		node := iSplit[0]
		re := regexp.MustCompile("\\(|\\)")
		nextNodes := re.ReplaceAll([]byte(iSplit[1]), []byte(""))
		nextSplit := strings.Split(string(nextNodes), ",")

		// fmt.Printf("Node=%s\tNextNodes=%s\n", node, nextSplit)

		out[node] = nextSplit
	}

	return out
}

func traverseInstructions(
	instructions string,
	nodes map[string][]string,
	startNode string,
	destNode string) int {

	i := 0
	nSteps := 0
	currentNode := startNode

	instructMap := make(map[string]int)
	instructMap["L"] = 0
	instructMap["R"] = 1

	for {
		// reset instructions
		if i >= len(instructions) {
			i = 0
		}

		// reached destination
		if currentNode == destNode {
			break
		}

		nextInstruct := string(instructions[i])
		nextIdx := instructMap[nextInstruct]

		currentNode = nodes[currentNode][nextIdx]
		i++
		nSteps++
	}

	return nSteps
}

func allNodesEndWithZ(nodes []string) bool {
	for _, i := range nodes {
		if string(i[len(i)-1]) != "Z" {
			return false
		}
	}
	return true
}


func traverseInstructionsUntilEndsWithZ(
	instructions string,
	nodes map[string][]string,
	startNodes []string) int{

	i := 0
	nSteps := 0
	currentNodes := startNodes

	instructMap := make(map[string]int)
	instructMap["L"] = 0
	instructMap["R"] = 1

	for {
		// reset instructions
		if i >= len(instructions) {
			i = 0
		}

		// reached destination
		if allNodesEndWithZ(currentNodes) {break}

		for idx, cNode := range currentNodes {
			nextInstruct := string(instructions[i])
			nextIdx := instructMap[nextInstruct]
			newNode := nodes[cNode][nextIdx]
			currentNodes[idx] = newNode
		}

		// if nSteps % 1000 == 0 {
		// 	printTraversalStatus(currentNodes, nSteps)
		// }

		i++
		nSteps++
		fmt.Printf("Steps=%d Nodes=%s\r", nSteps, currentNodes)
	}
	fmt.Printf("\n")

	return nSteps
}

func getAllNodesEndingWith(nodes map[string][]string, endsWith string) []string {
	out := *new([]string)

	for k, _ := range nodes {
		if string( k[len(k)-1] ) == endsWith { out = append(out, k) }
	}

	return out
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func printNodes(nodes []string) {
	maxElemPerRow := 3
	currNElem := 0
	for _, node := range nodes {
		fmt.Printf("[%s]", node)

		currNElem++

		if currNElem == maxElemPerRow {
			currNElem = 0
			fmt.Print("\n\n")
		} else {
			fmt.Print("\t")
		}
	}
	fmt.Print("\n")
}

func printTraversalStatus(nodes []string, nSteps int) {
	clearScreen()
	fmt.Printf("N. Steps = %d\n\n", nSteps)
	printNodes(nodes)
}

func solution1(data []string) int {
	instructions := data[0]
	nodesRaw := data[1:]

	// fmt.Println(instructions)
	// fmt.Println(nodes)

	nodes := parseNodes(nodesRaw)
	nSteps := traverseInstructions(instructions, nodes, "AAA", "ZZZ")

	return nSteps
}

func solution2(data []string) int {
	instructions := data[0]
	nodesRaw := data[1:]

	nodes := parseNodes(nodesRaw)
	startNodes := getAllNodesEndingWith(nodes, "A")
	fmt.Println(startNodes)
	nSteps := traverseInstructionsUntilEndsWithZ(instructions, nodes, startNodes)

	return nSteps
}

func main() {
	// data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(readLines("input1.txt")))
	fmt.Printf("Solution 2: %d\n", solution2(readLines("input1.txt")))
}

