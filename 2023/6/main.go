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

type RaceRecord struct {
	time int
	distance int
}

func computeDistance(holdTime int, raceTime int) int {
	// holdTime is also speed
	remainingTime := raceTime -  holdTime
	return holdTime * remainingTime
}

func isRaceWin(holdTime int, raceRecord RaceRecord) bool {
	// NOTE: maybe >= is not a good idea
	return computeDistance(holdTime, raceRecord.time) > raceRecord.distance
}

func searchLowerBound(holdTimes []int, raceRecord RaceRecord, pb RaceRecord) RaceRecord {
	for _, i := range holdTimes {
		if isRaceWin(i, raceRecord) {
			// return pb
			if i < pb.time {
				return RaceRecord{i, computeDistance(i, raceRecord.time)}
			}

		}
	}
	return pb
}

func searchUpperBound(holdTimes []int, raceRecord RaceRecord, pb RaceRecord) RaceRecord {
	for i := len(holdTimes)-1; i>=0; i-- {
	// for _, i := range holdTimes {
		if isRaceWin(holdTimes[i], raceRecord) {
			if i > pb.time {
				return RaceRecord{i, computeDistance(holdTimes[i], raceRecord.time)}
			}

		}
	}
	return pb
}

func genArr(a int) []int {
	out := *new([]int)
	for i := 0; i < a; i++ {
		out = append(out, i)
	}
	return out
}

func product(a []int) int {
	out := 1
	for _, i := range a {
		out *= i
	}
	return out
}

func solution1(data []RaceRecord) int {
	diffs := *new([]int)
	for _, i := range data {
		holdTimes := genArr(i.time + 1)
		lower := searchLowerBound(holdTimes, i, RaceRecord{i.time,0})
		upper := searchUpperBound(holdTimes, i, RaceRecord{0,0})
		diff := upper.time - lower.time + 1
		diffs = append(diffs, diff)
		fmt.Printf("Result %d, %d (%d)\n", lower, upper, diff)
	}

	return product(diffs)
}

func solution2(data []RaceRecord) []int {
	diffs := *new([]int)
	for _, i := range data {
		holdTimes := genArr(i.time + 1)
		lower := searchLowerBound(holdTimes, i, RaceRecord{i.time,0})
		highTime := i.time - lower.time
		upper := RaceRecord{highTime, computeDistance(highTime, i.time)}
		diff := upper.time - lower.time + 1
		diffs = append(diffs, diff)
		fmt.Printf("Result %d, %d (%d)\n", lower, upper, diff)
	}

	return diffs
}

func main() {
	data := []RaceRecord {
		{40, 277},
		{82, 1338},
		{91, 1349},
		{66, 1063},
		{71530, 940200},
		{40829166, 277133813491063},
	}

	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

