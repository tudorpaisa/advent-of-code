package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
	"math"
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

type Block struct {
	startX int
	startY int
	startZ int
	endX int
	endY int
	endZ int
}

func (a Block) compare(b Block) int {
	if a.startZ == b.startZ {
		if a.startY == b.startY {
			return a.startX - b.startX
		}
		return a.startY - b.startY
	}
	return a.startZ - b.startZ
}

func (a Block) overlaps(b Block) bool {
	xOverlap := false
	if math.Max(float64(a.startX), float64(b.startX)) <= math.Min(float64(a.endX), float64(b.endX)) {
		xOverlap = true
	}
	yOverlap := false
	if math.Max(float64(a.startY), float64(b.startY)) <= math.Min(float64(a.endY), float64(b.endY)) {
		yOverlap = true
	}

	return xOverlap && yOverlap
}

func sortBlocks(a []Block) {
	if len(a) < 2 {
		return
	}

	p := len(a) / 2
	l, h := 0, len(a)-1

	a[p], a[h] = a[h], a[p]

	for i := 0; i < len(a); i++ {
		if a[i].compare(a[h]) < 0 {
			a[l], a[i] = a[i], a[l]
			l++
		}
	}

	a[h], a[l] = a[l], a[h]

	sortBlocks(a[:l])
	sortBlocks(a[l+1:])
}


type NamedBlock struct {
	name int
	block Block
}

func sortNamedBlocks(a []NamedBlock) {
	if len(a) < 2 {
		return
	}

	p := len(a) / 2
	l, h := 0, len(a)-1

	a[p], a[h] = a[h], a[p]

	for i := 0; i < len(a); i++ {
		if a[i].block.compare(a[h].block) < 0 {
			a[l], a[i] = a[i], a[l]
			l++
		}
	}

	a[h], a[l] = a[l], a[h]

	sortNamedBlocks(a[:l])
	sortNamedBlocks(a[l+1:])
}

func parseCoords(coords string) (int, int, int) {
	var x, y, z int
	var err error

	split := strings.Split(coords, ",")
	if len(split) != 3 { panic("Bad Coord Parsing") }

	x, err = strconv.Atoi(split[0])
	if err != nil { panic(err) }
	y, err = strconv.Atoi(split[1])
	if err != nil { panic(err) }
	z, err = strconv.Atoi(split[2])
	if err != nil { panic(err) }


	return x, y, z
}

func parse(data []string) []Block {
	out := []Block{}

	for _, i := range data {
		coords := strings.Split(i, "~")
		if len(coords) != 2 { panic("Bad Coords") }

		start, end := coords[0], coords[1]
		startX, startY, startZ := parseCoords(start)
		endX, endY, endZ := parseCoords(end)

		out = append(out, Block{startX, startY, startZ, endX, endY, endZ})
	}
	return out
}

func allZeros(a []int) bool {
	for _, i := range a {
		if i != 0 { return false }
	}
	return true
}

type BlockMap struct {
	xz [][]int
	yz [][]int

	width int
	depth int

	blocks map[int]Block
	blockCount int

	supportedBy map[int][]int
	supports map[int][]int
}

func (m *BlockMap) draw(b Block, value int) {
	for i := b.startZ; i <= b.endZ; i++ {
		for j := b.startX; j <= b.endX; j++ {
			m.xz[i][j] = value
		}
		for j := b.startY; j <= b.endY; j++ {
			m.yz[i][j] = value
		}
	}
}

func (m *BlockMap) addBlock(b Block) {

	m.blockCount += 1

	if b.startZ + (b.endZ - b.startZ)+1 >= len(m.xz) {
		diff := b.startZ + (b.endZ - b.startZ)+1 - len(m.xz)
		for i:=0; i<=diff+1; i++ {
			m.xz = append(m.xz, make([]int, m.width))
			m.yz = append(m.yz, make([]int, m.depth))
		}
	}

	m.draw(b, m.blockCount)

	m.blocks[m.blockCount] = b
}

func (m *BlockMap) collapse() {
	// doesn't work :-(
	namedBlocks := []NamedBlock{}
	for k, v := range m.blocks { namedBlocks = append(namedBlocks, NamedBlock{k, v}) }

	sortNamedBlocks(namedBlocks[:])

	for _, nb := range namedBlocks {
		if nb.block.startX == 0 && nb.block.startY == 0 && nb.block.startZ == 0 && nb.block.endX == 0 && nb.block.endY == 0 && nb.block.endZ == 0 {
			continue
		}
		if nb.name == -1 { continue }
		b := nb.block

		supportedBy := make(map[int]bool)
		canGoFurther := true
		// check if can go below
		// x
		for canGoFurther {
			below := m.xz[b.startZ-1][b.startX : b.endX + 1]
			nOks := 0
			for _, j := range below {
				if j == 0 { nOks++ } else { supportedBy[j] = true }
			}
			if nOks == len(below) {
				m.draw(b, 0)  // clear space
				b.startZ = b.startZ - 1  // update startZ
				b.endZ = b.endZ - 1  // update endZ
				m.draw(b, nb.name)  // redraw
			} else {
				canGoFurther = false
			}
		}
		// y
		canGoFurther = true
		for canGoFurther {
			below := m.yz[b.startZ-1][b.startY : b.endY + 1]
			nOks := 0
			for _, j := range below {
				if j == 0 { nOks++ } else { supportedBy[j] = true }
			}
			if nOks == len(below) {
				m.draw(b, 0)  // clear space
				b.startZ = b.startZ - 1  // update startZ
				b.endZ = b.endZ - 1  // update endZ
				m.draw(b, nb.name)  // redraw
			} else {
				canGoFurther = false
			}
		}

		sB := []int{}
		for k := range supportedBy { sB = append(sB, k) }

		// update
		m.blocks[nb.name] = b
		m.supportedBy[nb.name] = sB
	}
	m.strip()

	for t, b := range m.supportedBy {
		for _, i := range b {
			if arr, ok := m.supports[i]; ok {
				arr = append(arr, t)
				m.supports[i] = arr
			} else {
				m.supports[i] = []int{t}
			}
		}
	}
}

func (m *BlockMap) strip() {
	for i := len(m.xz)-1; i>=0; i-- {
		xzRow := m.xz[i]
		yzRow := m.yz[i]
		xzEmpty, yzEmpty := false, false
		if allZeros(xzRow) { xzEmpty = true }
		if allZeros(yzRow) { yzEmpty = true }
		if xzEmpty && yzEmpty {
			m.xz = m.xz[:i]
			m.yz = m.yz[:i]
		}
	}
}

func calculateXYDims(blocks []Block) (int, int) {
	maxX := 0
	maxY := 0

	for _, i := range blocks {
		width := i.startX + (i.endX - i.startX)
		if width > maxX {
			maxX = width
		}
		depth := i.startY + (i.endY - i.startY)
		if depth > maxY {
			maxY = depth
		}
	}

	// NOTE: maybe add 1?
	return maxX + 1, maxY + 1
}

func initBlockMap(width int, depth int) BlockMap {
	xz := [][]int{}
	xzInit := []int{}
	for i:=0; i<width; i++ { xzInit = append(xzInit, -1) }
	xz = append(xz, xzInit)

	yz := [][]int{}
	yzInit := []int{}
	for i:=0; i<depth; i++ { yzInit = append(yzInit, -1) }
	yz = append(yz, yzInit)

	m := BlockMap{xz, yz, width, depth, map[int]Block {-1: Block{} }, 0, make(map[int][]int), make(map[int][]int)}

	return m
}

func printSide(a [][]int) {
	fmt.Print("+")
	for i := 0; i < len(a[0]); i++ {
		fmt.Print("-----")
	}
	fmt.Print("+\n")

	for i := len(a)-1; i>=0; i-- {
		fmt.Printf("%5d | ", i)

		for _, v := range a[i] {
			if v == 0 {
				fmt.Print("      | ")
			} else if v == -1 {
				fmt.Print("..... | ")
			} else {
				fmt.Printf("%5d | ", v)
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("+")
	for i := 0; i < len(a[0]); i++ {
		fmt.Print("-----")
	}
	fmt.Print("+\n")
	fmt.Print("\n")
}

func isRemovable(key int, supportedBy, supports map[int][]int) int {
	sList, ok := supports[key]
	if !ok {
		return 1
	}

	for _, i := range sList {
		sbl := len(supportedBy[i])
		if sbl == 1 {
			return 0
		}
	}

	return 1
}

func solution1(data []string) int {
	out := 0
	blocks := parse(data)

	sortBlocks(blocks[:])
	// fall logic
	for idx, a := range blocks {
		maxZ := 1
		for _, b := range blocks[:idx] {
			if a.overlaps(b) {
				maxZ = int(math.Max(float64(maxZ), float64(b.endZ) + 1))
			}
		}
		a.endZ = a.endZ - (a.startZ - maxZ)
		a.startZ = maxZ

		blocks[idx] = a
	}
	sortBlocks(blocks[:])

	supportersMap := make(map[int]map[int]bool)
	supportedByMap := make(map[int]map[int]bool)
	for j, top := range blocks {
		for i, bot := range blocks[:j] {
			if bot.overlaps(top) && top.startZ == (bot.endZ + 1) {

				if val, ok := supportersMap[i]; ok {
					val[j] = true
					supportersMap[i] = val
				} else {
					supportersMap[i] = map[int]bool {j: true}
				}

				if val, ok := supportedByMap[j]; ok {
					val[i] = true
					supportedByMap[j] = val
				} else {
					supportedByMap[j] = map[int]bool {i: true}
				}
			}
		}
	}

	supporters := make(map[int][]int)
	for i := range blocks { supporters[i] = []int{} }
	for k, vals := range supportersMap {
		v := []int{}
		for i := range vals {
			v = append(v, i)
		}
		supporters[k] = v
	}
	supportedBy := make(map[int][]int)
	for i := range blocks { supportedBy[i] = []int{} }
	for k, vals := range supportedByMap {
		v := []int{}
		for i := range vals {
			v = append(v, i)
		}
		supportedBy[k] = v
	}

	// fmt.Printf("%d\n", blocks)
	// fmt.Printf("%d\n", supporters)
	// fmt.Printf("%d\n", supportedBy)

	for i := range blocks {
		s := supporters[i]
		if len(s) == 0 {
			out++
			continue
		}

		canRm := true
		for _, j := range s {
			if len(supportedBy[j]) < 2 {
				canRm = false
			}
		}
		if canRm {out++}
	}


	// width, depth := calculateXYDims(blocks)
	// fmt.Printf("%d, %d\n", width, depth)
	// m := initBlockMap(width, depth)

	// namedBlocks := []NamedBlock{}
	// for k, v := range blocks {
	// 	if k > 10 { continue }
	// 	namedBlocks = append(namedBlocks, NamedBlock{k, v})
	// }

	// fmt.Printf("%d\n", namedBlocks)
	// sortNamedBlocks(namedBlocks[:])
	// fmt.Printf("%d\n", namedBlocks)

	// for _, nb := range namedBlocks {
	// 	m.addBlock(nb.block)
	// 	m.collapse()
	// }
	// printSide(m.xz)
	// printSide(m.yz)
	// for i := 1; i <= m.blockCount; i++ {
	// 	fmt.Printf("%d : %d\n", i, m.supportedBy[i])
	// }
	// fmt.Printf("%d\n", m.supportedBy)
	// fmt.Printf("%d\n", m.supports)

	// for k := range m.supportedBy {
	// 	out += isRemovable(k, m.supportedBy, m.supports)
	// }

	return out
}

func intersection(a []int, b []int) []int {
	out := []int{}

	counts := map[int]int{}
	for _, i := range a {
		if v, ok := counts[i]; ok {
			v = v + 1
			counts[i] = v
		} else {
			counts[i] = 1
		}
	}

	for _, i := range b {
		if v, ok := counts[i]; ok {
			v = v + 1
			counts[i] = v
		} else {
			counts[i] = 1
		}
	}

	for k, v := range counts {
		if v >= 2 {
			out = append(out, k)
		}
	}

	return out
}

func set(a []int) []int {
	x := map[int]bool{}
	for _, i := range a {
		x[i] = true
	}
	out := []int{}
	for k := range x { out = append(out, k) }
	return out
}

func in(a int, b []int) bool {
	for _, i := range b {
		if i == a { return true }
	}
	return false
}

func aInB(a []int, b []int) bool {
	for _, i := range a {
		if !in(i, b) { return false }
	}
	return true
}

func solution2(data []string) int {
	out := 0
	blocks := parse(data)

	sortBlocks(blocks[:])
	// fall logic
	for idx, a := range blocks {
		maxZ := 1
		for _, b := range blocks[:idx] {
			if a.overlaps(b) {
				maxZ = int(math.Max(float64(maxZ), float64(b.endZ) + 1))
			}
		}
		a.endZ = a.endZ - (a.startZ - maxZ)
		a.startZ = maxZ

		blocks[idx] = a
	}
	sortBlocks(blocks[:])

	supportersMap := make(map[int]map[int]bool)
	supportedByMap := make(map[int]map[int]bool)
	for j, top := range blocks {
		for i, bot := range blocks[:j] {
			if bot.overlaps(top) && top.startZ == (bot.endZ + 1) {

				if val, ok := supportersMap[i]; ok {
					val[j] = true
					supportersMap[i] = val
				} else {
					supportersMap[i] = map[int]bool {j: true}
				}

				if val, ok := supportedByMap[j]; ok {
					val[i] = true
					supportedByMap[j] = val
				} else {
					supportedByMap[j] = map[int]bool {i: true}
				}
			}
		}
	}

	supporters := make(map[int][]int)
	for i := range blocks { supporters[i] = []int{} }
	for k, vals := range supportersMap {
		v := []int{}
		for i := range vals {
			v = append(v, i)
		}
		supporters[k] = v
	}
	supportedBy := make(map[int][]int)
	for i := range blocks { supportedBy[i] = []int{} }
	for k, vals := range supportedByMap {
		v := []int{}
		for i := range vals {
			v = append(v, i)
		}
		supportedBy[k] = v
	}


	for i := range blocks {
		queue := []int{}
		s := supporters[i]
		for _, j := range s {
			if len(supportedBy[j]) == 1 {
				queue = append(queue, j)
			}
		}
		tofall := set(queue)
		tofall = append([]int{i}, tofall...)

		for len(queue) > 0 {
			j := queue[0]
			queue = queue[1:]

			for _, x := range supporters[j] {
				if in(x, tofall) {
					continue
				}
				if aInB(supportedBy[x], tofall) {
					queue = append(queue, x)

					tofall = append(tofall, x)
					tofall = set(tofall)
				}

			}
		}

		out += len(tofall)-1
	}

	return out
}


func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

