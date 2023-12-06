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
	return computeDistance(holdTime, raceRecord.time) >= raceRecord.distance
}

func searchLowerBound(holdTimes []int, raceRecord RaceRecord, pb RaceRecord) RaceRecord {
	if len(holdTimes) == 0 {
		return pb
	}
	// find min holdTime that gives a win
	var mid int = len(holdTimes) / 2

	if isRaceWin(mid, raceRecord) {
		if mid < pb.time {
			pb = RaceRecord{mid, computeDistance(mid, raceRecord.time)}
		}

		return searchLowerBound(holdTimes[:mid], raceRecord, pb)
	}

	return searchLowerBound(holdTimes[mid+1:], raceRecord, pb)
}

func searchUpperBound(holdTimes []int, raceRecord RaceRecord, pb RaceRecord) RaceRecord {
	if len(holdTimes) == 0 {
		return pb
	}
	// find max holdTime that gives a win
	var mid int = len(holdTimes) / 2

	if isRaceWin(mid, raceRecord) {
		if mid > pb.time {
			pb = RaceRecord{mid, computeDistance(mid, raceRecord.time)}
		}

		return searchUpperBound(holdTimes[mid+1:], raceRecord, pb)
	}

	return searchUpperBound(holdTimes[:mid], raceRecord, pb)
}

func genArr(a int) []int {
	out := *new([]int)
	for i := 0; i < a; i++ {
		out = append(out, i)
	}
	return out
}

func solution1(data []RaceRecord) int {
	for _, i := range data {
		holdTimes := genArr(i.time + 1)
		lower := searchLowerBound(holdTimes, i, RaceRecord{i.time,0})
		upper := searchUpperBound(holdTimes, i, RaceRecord{0,0})
		fmt.Printf("Result %d, %d (%d)\n", lower, upper, upper.time - lower.time)

	}
	return 0
}

func main() {
	data := []RaceRecord {
		{7, 9},
		{15, 40},
		{30, 200},
	}

	fmt.Printf("Solution 1: %d\n", solution1(data))
}

