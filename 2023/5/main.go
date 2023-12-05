package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
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

func strArrToInt(data []string) []int {
    out := *new([]int)
    for _, i := range data {
        num, err := strconv.Atoi(i)
        if err != nil {
            panic(err)
        }
        out = append(out, num)
    }

    return out
}

type Range struct {
    start int
    end int
}

func initRange(start int, end int) Range {
    return Range{start, end}
}

type MapRange struct {
    destinationRange Range
    sourceRange Range
}

func initMapRange(destStart int, destEnd int, srcStart int, srcEnd int) MapRange {
    return MapRange{Range{destStart, destEnd}, Range{srcStart, srcEnd}}
}

func buildMapRange(arr []int) MapRange {
    dest := arr[0]
    src := arr[1]
    size := arr[2]

    return initMapRange( dest, dest + size -1, src, src + size - 1)
}

func extractSeeds(data string) []int {
    splitData := strings.Split(data, " ")

    out := strArrToInt(splitData[1:])
    return out
}

func isMapString(i string) bool {
    // we're being cheeky here; check if the string ends
    // with a `:`
    
    return string(i[len(i)-1]) == ":"
}

func isRange(i string) bool {
    return len(strings.Split(i, " ")) == 3
}

func extractMaps(data []string) map[string][]MapRange {
    out := make(map[string][]MapRange)
    currentMap := ""
    for _, i := range data {
        if isMapString(i) {
            currentMap = i
        } else if isRange(i) {
            if arr, ok := out[currentMap]; ok {
                arr = append(arr, buildMapRange(strArrToInt(strings.Split(i, " "))))
                out[currentMap] = arr
            } else {
                arr := *new([]MapRange)
                arr = append(arr, buildMapRange(strArrToInt(strings.Split(i, " "))))
                out[currentMap] = arr
            }
            
        } else {
            panic(
                fmt.Sprintf(
                    "Starting on the wrong foot! Expected a map name/range but got: '%s'",
                    i))
        }
    }

    return out
}

func printMap(a map[string][]MapRange) {
    for k, arr := range a {
        fmt.Printf("%s\n", k)
        for _, i := range arr {
            fmt.Printf(
                "    dest=[%d, %d] | src=[%d, %d]\n", 
                i.destinationRange.start,
                i.destinationRange.end,
                i.sourceRange.start,
                i.sourceRange.end)
        }
    }
}

func getNextValue(val int, mapRanges []MapRange) int {
    for _, i := range mapRanges {
        if val >= i.sourceRange.start && val <= i.sourceRange.end {
            inc := val - i.sourceRange.start
            return i.destinationRange.start + inc
        }
    }
    return val
}

func fromRangeToSeeds(arr []int) []int {
    out := make([]int, arr[1])
    for i := 0; i < arr[1]; i++ {
        out = append(out, arr[0] + i)
    }
    return out
}

func solution1(data []string) int {
    seeds := extractSeeds(data[0])
    mapNames := []string {"seed-to-soil map:", "soil-to-fertilizer map:", "fertilizer-to-water map:", "water-to-light map:", "light-to-temperature map:", "temperature-to-humidity map:", "humidity-to-location map:"}
    // fmt.Printf("%d\n", seeds)
    maps := extractMaps(data[1:])

    numMap := make(map[int]int)

    for _, seed := range seeds {
        // fmt.Printf("%d -> ", seed)
        var interim int = seed
        var mapRanges []MapRange;

        for i := range mapNames {
            mapName := mapNames[i]
            mapRanges = maps[mapName]
            interim = getNextValue(interim, mapRanges)
            // fmt.Printf("%d -> ", interim)
        }

        // fmt.Printf("\b\b\b   \n")
        numMap[seed] = interim
    }
    var out int = 1000000000000000000;
    for _, v := range numMap {
        if v < out { out = v }
    }
    return out
}

func goComputeSeedMappings(
    seedRange []int,
    maps map[string][]MapRange,
    ch chan int,
    wg *sync.WaitGroup) {

    mapNames := []string {"seed-to-soil map:", "soil-to-fertilizer map:", "fertilizer-to-water map:", "water-to-light map:", "light-to-temperature map:", "temperature-to-humidity map:", "humidity-to-location map:"}
    // fmt.Printf("%d\n", seeds)

    fmt.Printf("Generating seed ranges: [%d, %d)\n", seedRange[0], seedRange[0]+seedRange[1])
    seeds := fromRangeToSeeds(seedRange)

    // _currSeed := 1
    nSeeds := len(seeds)
    var out int = 1000000000000000000;

    fmt.Printf("Iterating over %d seeds...\n", nSeeds)
    for _, seed := range seeds {
        // fmt.Printf("Seed %d / %d\r", _currSeed, nSeeds)
        // _currSeed++

        // fmt.Printf("%d -> ", seed)
        var interim int = seed
        var mapRanges []MapRange;

        for i := range mapNames {
            mapName := mapNames[i]
            mapRanges = maps[mapName]
            interim = getNextValue(interim, mapRanges)
            // fmt.Printf("%d -> ", interim)
        }

        // fmt.Printf("\b\b\b   \n")
        if interim < out { out = interim }
    }

    fmt.Printf("Ranges [%d, %d) complete!\n", seedRange[0], seedRange[0]+seedRange[1])
    ch <- out
    wg.Done()
}

func solution2(data []string) int {
    seeds := extractSeeds(data[0])
    maps := extractMaps(data[1:])

    nSeedGroups := len(seeds)/2
    ch := make(chan int, nSeedGroups)
    wg := sync.WaitGroup{}

    for i := 0; i < nSeedGroups; i++ {
        fmt.Printf("Creating goroutine #%d\n", i)
        idx := i*2
        wg.Add(1)
        go goComputeSeedMappings(seeds[idx : idx+2], maps, ch, &wg)
        // running one goroutine at a time; else the RAM will explode
        wg.Wait()
    }

    close(ch)

    fmt.Printf("Computing Min\n")
    var out int = 1000000000000000000;
    for i := range ch {
        if i < out { out = i }
    }
    return out
}

func main() {
    // line 0 -> seeds
    // line 1...N -> maps
    data := readLines("input1.txt")
    fmt.Printf("Solution 1: %d\n", solution1(data))
    fmt.Printf("Solution 2: %d\n", solution2(data))

}

