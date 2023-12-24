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

func preprocess(data []string) []Hail {
	out := []Hail{}
	for _, i := range data {
		row := strings.Replace(i, " @ ", ", ", 1)
		split := strings.Split(row, ", ")
		if len(split) != 6 { panic("Bad Data") }
		nums := []int{}
		for _, j := range split {
			n, err := strconv.Atoi(j)
			if err != nil {panic(err)}
			nums = append(nums, n)

		}

		// I hate math
		m := float64(nums[4]) / float64(nums[3])
		c := float64(nums[1]) - (m * float64(nums[0]))

		out = append(out, Hail{nums[0], nums[1], nums[2], nums[3], nums[4], nums[5], Line{m, -1.0, c}})
	}
	return out
}

func sum(a []float64) float64 {
	out := 0.0
	for _, i := range a { out += i }
	return out
}

func euclideanDistance(a, b []int) float64 {
	sqd := []float64{}
	for i := range a {
		sqd = append(sqd,  float64((a[i]-b[i]) * (a[i]-b[i])) )
	}
	return math.Sqrt( sum(sqd) )
}

type Line struct {
	 a float64
	 b float64
	 c float64
}

type Hail struct {
	x int
	y int
	z int
	vx int
	vy int
	vz int
	line Line
}

func (h Hail) compute(a int) (int, int, int) {
	return h.x + (h.vx * a), h.y + (h.vy * a), h.z + (h.vz * a)
}

func findXWithinXYBounds(h Hail, minB, maxB int) (int, bool) {
	it := 1
	minArr := []int{minB, minB}
	maxArr := []int{maxB, maxB}
	minDist := euclideanDistance([]int{h.x, h.y}, minArr)
	maxDist := euclideanDistance([]int{h.x, h.y}, maxArr)

	for {
		nx, ny, _ := h.compute(it)

		minD := euclideanDistance([]int{nx, ny}, minArr)
		maxD := euclideanDistance([]int{nx, ny}, maxArr)

		// if going towards the bounds
		if minD <= minDist || maxD <= maxDist {
			minDist = minD
			maxDist = maxD
		} else {
			break
		}
		it++
	}

	if it == 1 {
		return it, false
	}

	return it, true
}

func createValTable(h Hail, n int) [][]int {
	out := [][]int{}

	for i := 0; i <= n; i++ {
		x, y, z := h.compute(i)
		out = append(out, []int{x, y, z})
	}

	return out
}

func buildHailTables(h []Hail, minB, maxB int) [][][]int {
	out := [][][]int{}
	for _, i := range h {
		b, _ := findXWithinXYBounds(i, minB, maxB)
		t := createValTable(i, b)
		out = append(out, t)
	}
	return out
}

func withinBounds(v, minB, maxB int) bool {
	return v >= minB && v <= maxB
}

// doesn't work :\
func getNIntersects(hailTables [][][]int, minB, maxB int) int {
	nCross := 0
	for i := 0; i < len(hailTables)-1; i++ {
		a := hailTables[i]
		for j := i + 1; j < len(hailTables); j++ {
			b := hailTables[j]
			// fmt.Printf("%d (%d), %d (%d)\n", len(a), i, len(b), j)
			min := int(math.Min( float64(len(a)), float64(len(b)) ))

			cross := false

			for k := 0; k < min-1; k++ {

				acx := a[k][0]
				anx := a[k+1][0]
				acy := a[k][1]
				any := a[k+1][1]

				bcx := b[k][0]
				bnx := b[k+1][0]
				bcy := b[k][1]
				bny := b[k+1][1]

				// fmt.Printf("%3d, %3d v %3d. %3d\n", acx, acy, bcx, bcy)
				// fmt.Printf("%3d, %3d v %3d. %3d\n", anx, any, bnx, bny)
				// fmt.Print("\n")

				if withinBounds(acx, minB, maxB) && withinBounds(acy, minB, maxB) && withinBounds(bcx, minB, maxB) && withinBounds(bcy, minB, maxB) {
					// this is fugly!
					acxInB := withinBounds(acx, int(math.Min( float64(bcx) , float64(bnx))), int(math.Max( float64(bcx) , float64(bnx))))
					acyInB := withinBounds(acy, int(math.Min( float64(bcy) , float64(bny))), int(math.Max( float64(bcy) , float64(bny))))
					anxInB := withinBounds(anx, int(math.Min( float64(bcx) , float64(bnx))), int(math.Max( float64(bcx) , float64(bnx))))
					anyInB := withinBounds(any, int(math.Min( float64(bcy) , float64(bny))), int(math.Max( float64(bcy) , float64(bny))))

					bcxInB := withinBounds(bcx, int(math.Min( float64(acx) , float64(anx))), int(math.Max( float64(acx) , float64(anx))))
					bcyInB := withinBounds(bcy, int(math.Min( float64(acy) , float64(any))), int(math.Max( float64(acy) , float64(any))))
					bnxInB := withinBounds(bnx, int(math.Min( float64(acx) , float64(anx))), int(math.Max( float64(acx) , float64(anx))))
					bnyInB := withinBounds(bny, int(math.Min( float64(acy) , float64(any))), int(math.Max( float64(acy) , float64(any))))


					if acxInB || anxInB || bcxInB || bnxInB {
						if acyInB || anyInB || bcyInB || bnyInB {
							cross = true
							break
						}
					}
				}
			}

			if cross { nCross++ }
		}
	}
	return nCross

}


func countXYIntersects(hail []Hail, minB, maxB int) int {
	count := 0

	// I had to look this up; the math is new to me
	for i := 0; i < len(hail)-1; i++ {
		for j := i+1; j < len(hail); j++ {
			h1 := hail[i]
			h2 := hail[j]

			// intersection of two lines is when
			//   a1x + b1y + c1 = 0 AND a2x + b2y + c2 = 0
			// is calculated with
			//   ix := (b1 * c2 - b2 * c1) /  (a1 * b2 - a2 * b1)
			//   iy := (c1 * a2 - c2 * a1) /  (a1 * b2 - a2 * b1)

			a1, b1, c1 := h1.line.a, h1.line.b, h1.line.c
			a2, b2, c2 := h2.line.a, h2.line.b, h2.line.c

			den := (a1 * b2) - (a2 * b1)  // a1 * b2 - a2 * b1
			if den == 0 {continue}

			// intersection points
			ix := (b1 * c2 - b2 * c1) / den
			iy := (c1 * a2 - c2 * a1) / den

			// time should be > 0
			//   x(t) = x0 + vt
			//   t = x(t) - x0 / v
			//   x(t) >= t
			// and then this comes up...
			t1 := (ix - float64(h1.x)) / float64(h1.vx)
			t2 := (ix - float64(h2.x)) / float64(h2.vx)

			if t1 > 0 && t2 > 0 && ix >= float64(minB) && ix <= float64(maxB) && iy >= float64(minB) && iy <= float64(maxB) {
				count ++
			}
		}
	}

	return count
}

func solution1(data []string) int {
	minB, maxB := 200000000000000, 400000000000000
	hail := preprocess(data)
	return countXYIntersects(hail, minB, maxB)
	// hailTables := buildHailTables(hail, minB, maxB)
	// return getNIntersects(hailTables, minB, maxB)
}

func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Print("Solution 2: LOL\n")
}

