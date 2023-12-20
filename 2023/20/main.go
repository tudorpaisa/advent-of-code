package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	// "time"
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

//---CONST---//
const BUTTON = "button"
const BROADCASTER = "broadcaster"

//---PULSE---//

type Pulse int
const (
	HighPulse = iota
	LowPulse
	NoPulse
)

func initPulseCounter() map[Pulse]int {
	out := make(map[Pulse]int)
	out[HighPulse] = 0
	out[LowPulse] = 0
	return out
}

//---STATES---//

type State int
const (
	On State = iota
	Off
)

//---MODULES--//

type Module interface {
	process(string, Pulse)
	sendPulse() Pulse
	refresh()
	getState() State
}

// flip-flop
type FlipFlopModule struct {
	isOn bool
	canFire bool
}

func (m FlipFlopModule) getState() State {if m.isOn { return On } else { return Off}}

func (m *FlipFlopModule) process(key string, p Pulse) {
	if p == HighPulse {
		m.canFire = false
	} else if p == LowPulse {
		m.canFire = true
	}
}

func (m *FlipFlopModule) sendPulse() Pulse {
	var out Pulse = NoPulse
	if m.canFire {
		m.isOn = !m.isOn
		if m.isOn {
			out = HighPulse
		} else {
			out = LowPulse
		}
	}
	return out
}

func (m FlipFlopModule) refresh() {
	// noop
}

func initFlipFlopModule() FlipFlopModule {
	return FlipFlopModule{false, false}
}

// conjunction
type ConjunctionModule struct {
	// this might have to know _all_ inputs beforehand
	inputHistory map[string]Pulse
}
func (m ConjunctionModule) getState() State {return On}

func (m *ConjunctionModule) process(key string, p Pulse) {
	m.inputHistory[key] = p
}

func (m *ConjunctionModule) sendPulse() Pulse {
	nHighPulses := 0
	for _, i := range m.inputHistory { if i == HighPulse { nHighPulses++ } }
	if nHighPulses == len(m.inputHistory) {
		return LowPulse
	}
	return HighPulse
}

func (m ConjunctionModule) refresh() {
	// noop
}

func initConjunctionModule() ConjunctionModule {
	return ConjunctionModule{make(map[string]Pulse)}
}

// broadcaster
type BroadcasterModule struct {
	lastPulse Pulse
}
func (m BroadcasterModule) getState() State {return On}

func (m *BroadcasterModule) process(key string, p Pulse) {
	m.lastPulse = p
}

func (m *BroadcasterModule) sendPulse() Pulse {
	return m.lastPulse
}

func (m BroadcasterModule) refresh() {
	// noop
}

func initBroadcasterModule() BroadcasterModule {
	return BroadcasterModule{NoPulse}
}

// untyped
type UntypedModule struct {}
func (m UntypedModule) getState() State {return On}

func (m UntypedModule) process(key string, p Pulse) {}

func (m UntypedModule) sendPulse() Pulse {return NoPulse}

func (m UntypedModule) refresh() {
	// noop
}

func initUntypedModule() UntypedModule {
	return UntypedModule{}
}

// button
type ButtonModule struct {}
func (m ButtonModule) getState() State {return On}

func (m ButtonModule) process(key string, p Pulse) {}

func (m ButtonModule) sendPulse() Pulse {return LowPulse}

func (m ButtonModule) refresh() {
	// noop
}

func initButtonModule() ButtonModule {
	return ButtonModule{}
}

//---FUNCTIONS---//

func parseConfiguration(data []string) (map[string]Module, map[string][]string) {
	destinationMap := make(map[string][]string)
	moduleMap := make(map[string]Module)
	// keeps track of all modules, and which are assigned a type
	assignmentMap := make(map[string]bool)
	// keeps track of what modules (values) send to a conjunction module (key)
	conjunctionMap := make(map[string][]string)

	for _, row := range data {
		rowSplit := strings.Split(row, " -> ")
		if len(rowSplit) != 2 {panic("Parsing Error")}

		origin := rowSplit[0]
		destinations := strings.Split(rowSplit[1], ", ")

		if origin == BROADCASTER {
			destinationMap[origin] = destinations
			m := initBroadcasterModule()
			moduleMap[origin] = &m
			assignmentMap[origin] = true
		} else if origin[0] == '%' {
			name := origin[1:]
			m := initFlipFlopModule()
			moduleMap[name] = &m
			destinationMap[name] = destinations
			assignmentMap[name] = true
		} else if origin[0] == '&' {
			name := origin[1:]
			m := initConjunctionModule()
			moduleMap[name] = &m
			destinationMap[name] = destinations
			assignmentMap[name] = true
			conjunctionMap[name] = []string{}
		}

		for _, i := range destinations {
			// add dest as unassigned if not in the map
			if _, ok := assignmentMap[i]; !ok {
				assignmentMap[i] = false
			}
		}
	}


	for source, destinations := range destinationMap {
		for _, dest := range destinations {
			// if dest in conjunction map
			if lst, ok := conjunctionMap[dest]; ok {
				// update list
				// NOTE: will we have duplicates?
				lst = append(lst, source)
				conjunctionMap[dest] = lst
			}
		}
	}
	// init conjunction modules
	for conjName, sources := range conjunctionMap {
		for _, i := range sources {
			m := moduleMap[conjName]
			m.process(i, NoPulse)
			moduleMap[conjName] = m
		}
	}

	for k, v := range assignmentMap {
		// if not assigned
		if !v {
			moduleMap[k] = initUntypedModule()
			destinationMap[k] = []string{}
		}
	}

	// add button
	moduleMap[BUTTON] = initButtonModule()
	destinationMap[BUTTON] = []string{BROADCASTER}

	return moduleMap, destinationMap
}


func printModuleMap(m map[string]Module) {
	for k, v := range m {
		fmt.Printf("'%s' : %T\n", k, v)
	}
}

func printDestinationMap(m map[string][]string) {
	for k, v := range m {
		fmt.Printf("'%s' -> %s\n", k, v)
	}
}

func sameState(moduleMap map[string]Module, initStates map[string]State) bool {
	for k, v := range moduleMap {
		if initStates[k] != v.getState() { return false}
	}
	return true
}

type Message struct {
	sender string
	destination string
	pulse Pulse
}

func start(moduleMap map[string]Module, destinationMap map[string][]string, n int) map[Pulse]int {
	pulseCounts := make(map[Pulse]int)
	pulseCounts[NoPulse] = 0
	pulseCounts[HighPulse] = 0
	pulseCounts[LowPulse] = 0

	for it := 0; it<n; it++ {
		queue := []Message{ {BUTTON, BROADCASTER, LowPulse} }
		pulseCounts[LowPulse] = pulseCounts[LowPulse] + 1

		for len(queue) > 0 {
			msg := queue[0]
			queue = queue[1:]

			sName := msg.sender
			mName := msg.destination
			pulse := msg.pulse

			m := moduleMap[mName]
			m.process(sName, pulse)
			outPulse := m.sendPulse()
			if outPulse == NoPulse {
				pulseCounts[NoPulse] = pulseCounts[NoPulse] + 1
				continue
			}

			destinations, ok := destinationMap[mName]
			if !ok { panic("Wrong source!") }

			for _, i := range destinations {
				queue = append(queue, Message{mName, i, outPulse})
				pulseCounts[outPulse] = pulseCounts[outPulse] + 1
			}

		}


	}

	return pulseCounts
}

func in(a string, b []string) bool {
	for _, i := range b {
		if a == i { return true }
	}
	return false
}

func allExist(keys []string, m map[string]int) bool {
	if len(keys) != len(m) { return false }

	for _, k := range keys {
		if _, ok := m[k]; !ok {
			return false
		}
	}
	return true
}

func start2(moduleMap map[string]Module, destinationMap map[string][]string, targets []string) map[string]int {
	targetIterations := make(map[string]int)

	it := 1
	for {
		queue := []Message{ {BUTTON, BROADCASTER, LowPulse} }

		for len(queue) > 0 {
			msg := queue[0]
			queue = queue[1:]

			sName := msg.sender
			mName := msg.destination
			pulse := msg.pulse

			m := moduleMap[mName]
			m.process(sName, pulse)
			outPulse := m.sendPulse()
			if outPulse == NoPulse {
				continue
			}

			destinations, ok := destinationMap[mName]
			if !ok { panic("Wrong source!") }

			for _, i := range destinations {

				if in(mName, targets) && outPulse == HighPulse {
					if _, ok := targetIterations[mName]; !ok {
						targetIterations[mName] = it
					}
				}

				queue = append(queue, Message{mName, i, outPulse})
			}

		}

		if allExist(targets, targetIterations) {
			break
		}

		it++

	}

	return targetIterations
}

func solution1(data []string) int {
	moduleMap, destinationMap := parseConfiguration(data)

	pulseCount := start(moduleMap, destinationMap, 1000)
	return pulseCount[LowPulse] * pulseCount[HighPulse]
}

func gcd(a int, b int) int {
	if a < 1 || b < 1 { panic("a/b is less than 1") }
	for b != 0 {
		a, b = b, a % b
	}
	return a
}

func lcm(a int, b int) int {
	if a == 0 || b == 0 { return 0 }
	return (a*b) / gcd(a, b)
}

func solution2(data []string) int {
	moduleMap, destinationMap := parseConfiguration(data)

	sourceMap := make(map[string][]string)
	for src, dests := range destinationMap {
		for _, dest := range dests {
			if l, ok := sourceMap[dest]; ok {
				l = append(l, src)
				sourceMap[dest] = l
			} else {
				l := []string{src}
				sourceMap[dest] = l
			}
		}
	}

	rxSrc := sourceMap["rx"]
	if len(rxSrc) != 1 { panic("PANIK!") }
	targets := sourceMap[rxSrc[0]]

	out := 1
	targetIts := start2(moduleMap, destinationMap, targets)
	for _, v := range targetIts {
		out = lcm(out, v)
	}

	return out
}


func main() {
	data := readLines("input1.txt")
	fmt.Printf("Solution 1: %d\n", solution1(data))
	fmt.Printf("Solution 2: %d\n", solution2(data))
}

