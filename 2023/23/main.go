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

type Direction int
const (
	Up Direction = iota
	Down
	Left
	Right
	NoDirection
)

type Pos struct {
	x int
	y int
}

func isWalkable(x, y int, data []string) bool  {
	// bounds checking
	if y < 0 || y >= len(data) || x < 0 || x >= len(data[0]) { return false }

	// tree check
	if data[y][x] == '#' {
		return false
	}
	return true
}

func in(a Pos, b []Pos) bool {
	for _, i := range b {
		if a.x == i.x && a.y == i.y {
			return true
		}
	}
	return false
}

func isSlope(s string) bool {
	if s == ">" || s == "<" || s == "^" || s == "v" {
		return true
	}
	return false
}

func slopeToCoords(s string) (int, int) {
	switch s {
	case "^":
		return 0, -1
	case "v":
		return 0, 1
	case ">":
		return 1, 0
	case "<":
		return -1, 0
	default:
		return 0, 0
	}
}

func printMap(data []string) {
	fmt.Print("\n")
	for _, row := range data {
		fmt.Printf("%s\n", row)
	}
}

func drawMap(path []Pos, data []string) {
	var nd []string = data

	for _, p := range path {
		nd[p.y] = nd[p.y][:p.x] + "O" + nd[p.y][p.x+1:]
	}
	printMap(nd)
}

func traverse(currentPos, endPos Pos, data []string, currentPath []Pos) [][]Pos {

	if currentPos.x == endPos.x && currentPos.y == endPos.y {
		return [][]Pos{currentPath}
	}

	out := [][]Pos{}

	goTos := []Pos{}

	currentTile := string(data[currentPos.y][currentPos.x])
	// force direction in case of a slope
	if isSlope(currentTile) {
		_x, _y := slopeToCoords(currentTile)
		nP := Pos{currentPos.x + _x, currentPos.y + _y}

		// fmt.Printf("%s -> %s\n", string(data[currentPos.y][currentPos.x]), string(data[nP.y][nP.x]))
		if !in(nP, currentPath) {goTos = append(goTos, nP)}
	} else {
		// check direction
		// up
		if isWalkable(currentPos.x, currentPos.y - 1, data) {
			nP := Pos{currentPos.x, currentPos.y - 1}
			if data[nP.y][nP.x] != 'v' {
				if !in(nP, currentPath) {goTos = append(goTos, nP)}
			}
		}
		// down
		if isWalkable(currentPos.x, currentPos.y + 1, data) {
			nP := Pos{currentPos.x, currentPos.y + 1}
			if data[nP.y][nP.x] != '^' {
				if !in(nP, currentPath) {goTos = append(goTos, nP)}
			}
		}
		// right
		if isWalkable(currentPos.x - 1, currentPos.y, data) {
			nP := Pos{currentPos.x - 1, currentPos.y}
			if data[nP.y][nP.x] != '>' {
				if !in(nP, currentPath) {goTos = append(goTos, nP)}
			}
		}
		// left
		if isWalkable(currentPos.x + 1, currentPos.y, data) {
			nP := Pos{currentPos.x + 1, currentPos.y}
			if data[nP.y][nP.x] != '<' {
				if !in(nP, currentPath) {goTos = append(goTos, nP)}
			}
		}
	}

	for _, i := range goTos {
		// fmt.Printf("%d\n", i)
		newPaths := traverse(i, endPos, data, append(currentPath, i))
		out = append(out, newPaths...)
	}

	return out
}

func solution1(data []string) int {
	paths := traverse(Pos{1, 0}, Pos{len(data[0])-2, len(data)-1}, data, []Pos{{1, 0}})
	// paths := traverse(Pos{1, 0}, Pos{7, 3}, data, []Pos{{1, 0}})

	max := 0
	for _, i := range paths {
		// fmt.Println(i)
		// if len(i) > 120 {
		// 	fmt.Printf("%d\n", i)
		// 	drawMap(i, data)
		// }
		// fmt.Printf("%d\n", len(i))

		if len(i)-1 > max {
			max = len(i)-1
		}
	}
	// fmt.Printf("%d\n", paths)
	return max
}

func traverseDry(currentPos, endPos Pos, data []string, currentPath []Pos, max *int) [][]Pos {

	if currentPos.x == endPos.x && currentPos.y == endPos.y {
		return [][]Pos{currentPath}
	}

	out := [][]Pos{}

	goTos := []Pos{}

	// check direction
	// up
	if isWalkable(currentPos.x, currentPos.y - 1, data) {
		nP := Pos{currentPos.x, currentPos.y - 1}
		if !in(nP, currentPath) {goTos = append(goTos, nP)}
	}
	// down
	if isWalkable(currentPos.x, currentPos.y + 1, data) {
		nP := Pos{currentPos.x, currentPos.y + 1}
		if !in(nP, currentPath) {goTos = append(goTos, nP)}
	}
	// right
	if isWalkable(currentPos.x - 1, currentPos.y, data) {
		nP := Pos{currentPos.x - 1, currentPos.y}
		if !in(nP, currentPath) {goTos = append(goTos, nP)}
	}
	// left
	if isWalkable(currentPos.x + 1, currentPos.y, data) {
		nP := Pos{currentPos.x + 1, currentPos.y}
		if !in(nP, currentPath) {goTos = append(goTos, nP)}
	}

	for _, i := range goTos {
		newPaths := traverseDry(i, endPos, data, append(currentPath, i), max)
		nMax := 0
		for _, j := range newPaths {
			if len(j) > nMax { nMax = len(j) }
		}
		if nMax > *max {
			max = &nMax
			fmt.Printf("New Max: %d\r", *max)
		}

		out = append(out, newPaths...)
	}

	return out
}

type Node struct {
	size int
	children []*Node
	parents []*Node
	path []Pos
}

func (s Node) getSize() int {
	return s.size
}

func (s Node) getChildren() []*Node {
	return s.children
}

func (s *Node) addChild(n *Node) {
	s.children = append(s.children, n)
}

func (s *Node) addParent(n *Node) {
	s.parents = append(s.parents, n)
}

func (s *Node) incrementSize() {
	s.size += 1
}

func (s *Node) getLastPos() Pos {
	return s.path[len(s.path)-1]
}

func (s *Node) addToPath(p Pos) {
	s.path = append(s.path, p)
}

func (s *Node) popLastPath() Pos {
	p := s.getLastPos()
	s.path = s.path[:len(s.path)-1]
	return p
}

type Graph struct {
	graph []*Node
}

func (g *Graph) addNode(n *Node) {
	g.graph = append(g.graph, n)
}

func isVisited(p Pos, graph *Graph) bool {
	for _, i := range graph.graph {
		if in(p, i.path) {
			return true
		}
	}
	return false
}

func buildGraph(data []string, currentNode *Node, graph *Graph) {
	lastPos := currentNode.getLastPos()

	goTos := []Pos{}

	// check direction
	// up
	if isWalkable(lastPos.x, lastPos.y - 1, data) {
		nP := Pos{lastPos.x, lastPos.y - 1}
		if !in(nP, currentNode.path) && !isVisited(nP, graph) {goTos = append(goTos, nP)}
	}
	// down
	if isWalkable(lastPos.x, lastPos.y + 1, data) {
		nP := Pos{lastPos.x, lastPos.y + 1}
		if !in(nP, currentNode.path) && !isVisited(nP, graph) {goTos = append(goTos, nP)}
	}
	// right
	if isWalkable(lastPos.x - 1, lastPos.y, data) {
		nP := Pos{lastPos.x - 1, lastPos.y}
		if !in(nP, currentNode.path) && !isVisited(nP, graph) {goTos = append(goTos, nP)}
	}
	// left
	if isWalkable(lastPos.x + 1, lastPos.y, data) {
		nP := Pos{lastPos.x + 1, lastPos.y}
		if !in(nP, currentNode.path) && !isVisited(nP, graph) {goTos = append(goTos, nP)}
	}

	if len(goTos) == 0 {
		// finished
	} else if len(goTos) == 1 {
		currentNode.addToPath(goTos[0])
		currentNode.incrementSize()
		buildGraph(data, currentNode, graph)
	} else {
		// fmt.Print("Building ")
		// we're at a xroad and last node is the xroad
		lastPos = currentNode.popLastPath()

		// create xroad node
		xroad := Node{1, []*Node{}, []*Node{}, []Pos{lastPos}}
		currentNode.addChild(&xroad)
		xroad.addParent(currentNode)

		// update graph
		graph.addNode(currentNode)
		graph.addNode(&xroad)

		// gotos are branching paths of xroad
		for _, i := range goTos {
			xroad.addChild(&Node{1, []*Node{}, []*Node{}, []Pos{i}})
		}

		// continue from children
		for _, i := range xroad.getChildren() {
			buildGraph(data, i, graph)
		}

	}

}

func getNodeWithPos(pos Pos, graph *Graph) *Node {
	for _, i := range graph.graph {
		if in(pos, i.path) {
			return i
		}
	}
	return nil
}

func solution2(data []string) int {

	// startPos := Pos{1, 0}
	// endPos := Pos{ len(data[0])-2, len(data)-1 }

	// startNode := &Node{1, []*Node{}, []*Node{}, []Pos{endPos}}
	// graph := Graph{[]*Node{}}
	// buildGraph(data, startNode, &graph)

	// lastNode := getNodeWithPos(startPos, &graph)
	// if lastNode == nil {
	// 	fmt.Println("adjhasjhd")
	// }

	// fmt.Printf("%s\n", graph.graph)


	// Going to have to go with brute force this time around.
	// Due to life, I didn't have time today to write a better algorithm.
	// My idea above was to create a graph consisting of all the crossroads
	// as 1 node and all single lines as 1 node as well and computing the
	// longest path based on that. Again, I don't have enough time to make
	// it work, so this will do.
	// After ~8 hours it `traverseDry` prints out the result with an off-by-1
	// error. The program does not finish executing though :|

	m := 0
	paths := traverseDry(Pos{1, 0}, Pos{len(data[0])-2, len(data)-1}, data, []Pos{{1, 0}}, &m)
	fmt.Print("\n")

	max := 0
	for _, i := range paths {

		if len(i)-1 > max {
			max = len(i)-1
		}
	}
	return max
}


func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

