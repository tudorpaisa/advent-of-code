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

func readLines(fname string) ([]string, []string ){
	var parts []string = strings.Split(readFileRaw(fname), "\n\n")
        p1 := strings.Split(parts[0], "\n")
        p2 := strings.Split(parts[1], "\n")
	return p1, p2
}

type Part struct {
    x int
    m int
    a int
    s int
}

func parseRule(r string) (string, []func(Part)(string, bool)) {
    // well, this was a nice exercise. too bad it won't work for part 2
    s := strings.Split(r, "{")
    if len(s) != 2 { panic("Rule parsing error") }

    key := s[0]
    ruleStr := strings.Replace(s[1], "}", "", 1)
    rulesSplit := strings.Split(ruleStr, ",")

    rules := []func(Part) (string, bool) {}

    for _, i := range rulesSplit {
        rSplit := strings.Split(i, ":")
        if len(rSplit) == 1 {
            rules = append(rules, func(p Part) (string, bool) {
                return rSplit[0], true
            })
            continue
        }

        nextKey := rSplit[1] 
        ruleText := rSplit[0]

        cat := ruleText[0]
        sign := ruleText[1]
        thresh, err := strconv.Atoi(string(ruleText[2:]))
        if err != nil { panic(err) }

        switch cat {
        case 'x':
            if sign == '>' {
                rules = append(rules,
                    func(p Part) (string, bool) {
                        if p.x > thresh { return nextKey, true} else { return "", false }
                    })
            } else {
                rules = append(rules,
                    func(p Part) (string, bool) {
                        if p.x < thresh { return nextKey, true} else { return "", false }
                    })
            }
        case 'm':
            if sign == '>' {
                rules = append(rules,
                    func(p Part) (string, bool) {
                        if p.m > thresh { return nextKey, true} else { return "", false }
                    })
            } else {
                rules = append(rules,
                    func(p Part) (string, bool) {
                        if p.m < thresh { return nextKey, true} else { return "", false }
                    })
            }
        case 'a':
            if sign == '>' {
                rules = append(rules,
                    func(p Part) (string, bool) {
                        if p.a > thresh { return nextKey, true} else { return "", false }
                    })
            } else {
                rules = append(rules,
                    func(p Part) (string, bool) {
                        if p.a < thresh { return nextKey, true} else { return "", false }
                    })
            }
        case 's':
            if sign == '>' {
                rules = append(rules,
                    func(p Part) (string, bool) {
                        if p.s > thresh { return nextKey, true} else { return "", false }
                    })
            } else {
                rules = append(rules,
                    func(p Part) (string, bool) {
                        if p.s < thresh { return nextKey, true} else { return "", false }
                    })
            }
        default:
            panic("Wrong category")
        }
    }

    return key, rules
}

func parseParts(data []string) []Part {
    out := []Part{}
    for _, i := range data {
        if i == "" { continue }
        pStr := strings.Replace(strings.Replace(i, "{", "", 1), "}", "", 1)
        pSplit := strings.Split(pStr, ",")
        part := Part{0, 0, 0, 0}

        for _, j := range pSplit {
            jSplit := strings.Split(j, "=")
            value, err := strconv.Atoi(jSplit[1])
            if err != nil { panic(err) }
            switch jSplit[0] {
            case "x":
                part.x = value
            case "m":
                part.m = value
            case "s":
                part.s = value
            case "a":
                part.a = value
            default:
                panic("Bad category")
            }
        }

        out = append(out, part)
    }

    return out
}

func computeRuleResult(parts []Part, rulesMap map[string][]func(Part)(string, bool)) []Part {
    out := []Part{}
    for _, p := range parts {
        key := "in"

        // fmt.Println(p)

        finished := false
        for !finished {
            if rules, ok := rulesMap[key]; ok {

                // fmt.Printf("%s:\n", key)
                for _, rule := range rules {

                    result, pass := rule(p)
                    // fmt.Printf("    -> %s (%t)\n", result, pass)
                    if pass {
                        if result == "A" {
                            out = append(out, p)
                            finished = true
                        } else if result == "R" {
                            finished = true
                        } else {
                            key = result
                        }
                        break
                    }

                }
            } else {
                panic("Key not found")
            }

        }
    }
    return out
}

func splitArrAfter(a []int, b int) ([]int, []int) {
    for i, j := range a {
        if j == b {
            return a[:i+1], a[i+1:]
        }
    }
    return a, []int{}
}

type RangedPart struct {
    x []int
    m []int
    a []int
    s []int
}

type Rule interface {
    execute(RangedPart) (string, bool, RangedPart, RangedPart)
    getValue() int
    getSign() string
    getKey() string
}

type AcceptRule struct { }

func (r AcceptRule) execute(p RangedPart) (string, bool, RangedPart, RangedPart) {
    return "A", true, p, RangedPart{}
}
func (r AcceptRule) getValue() int { return -1 }
func (r AcceptRule) getSign() string { return "" }
func (r AcceptRule) getKey() string { return "A" }

type RejectRule struct {
}

func (r RejectRule) execute(p RangedPart) (string, bool, RangedPart, RangedPart) {
    return "R", true, p, RangedPart{}
}
func (r RejectRule) getValue() int { return -1 }
func (r RejectRule) getSign() string { return "" }
func (r RejectRule) getKey() string { return "R" }

type GoToRule struct {
    key string
}
func (r GoToRule) execute(p RangedPart) (string, bool, RangedPart, RangedPart) {
    return r.key, true, p, RangedPart{}
}
func (r GoToRule) getValue() int { return -1 }
func (r GoToRule) getSign() string { return "" }
func (r GoToRule) getKey() string { return r.key }

type DecisionRule struct {
    category string
    sign string
    value int
    key string
}

func (r DecisionRule) getKey() string { return r.key }
func (r DecisionRule) execute(p RangedPart) (string, bool, RangedPart, RangedPart) {
    if r.sign == "<" {
        switch r.category {
        case "x":
            if p.x[len(p.x)-1] < r.value {
                return r.key, true, p, RangedPart{}
            }
            b, a := splitArrAfter(p.x, r.value)
            return r.key, false, RangedPart{b, p.m, p.a, p.s}, RangedPart{a, p.m, p.a, p.s}

        case "m":
            if p.m[len(p.m)-1] < r.value {
                return r.key, true, p, RangedPart{}
            }
            b, a := splitArrAfter(p.m, r.value)
            return r.key, false, RangedPart{p.x, b, p.a, p.s}, RangedPart{p.x, a, p.a, p.s}

        case "a":
            if p.a[len(p.a)-1] < r.value {
                return r.key, true, p, RangedPart{}
            }
            b, a := splitArrAfter(p.a, r.value)
            return r.key, false, RangedPart{p.x, p.m, b, p.s}, RangedPart{p.x, p.m, a, p.s}

        case "s":
            if p.s[len(p.s)-1] < r.value {
                return r.key, true, p, RangedPart{}
            }
            b, a := splitArrAfter(p.s, r.value)
            return r.key, false, RangedPart{p.x, p.m, p.a, b}, RangedPart{p.x, p.m, p.a, a}

        default:
            panic("Wrong Category")
        }
    } else {
        switch r.category {
        case "x":
            if p.x[0] > r.value {
                return r.key, true, p, RangedPart{}
            }
            b, a := splitArrAfter(p.x, r.value)
            return r.key, false, RangedPart{a, p.m, p.a, p.s}, RangedPart{b, p.m, p.a, p.s}

        case "m":
            if p.m[0] > r.value {
                return r.key, true, p, RangedPart{}
            }
            b, a := splitArrAfter(p.m, r.value)
            return r.key, false, RangedPart{p.x, a, p.a, p.s}, RangedPart{p.x, b, p.a, p.s}

        case "a":
            if p.a[0] > r.value {
                return r.key, true, p, RangedPart{}
            }
            b, a := splitArrAfter(p.a, r.value)
            return r.key, false, RangedPart{p.x, p.m, a, p.s}, RangedPart{p.x, p.m, b, p.s}

        case "s":
            if p.s[0] > r.value {
                return r.key, true, p, RangedPart{}
            }
            b, a := splitArrAfter(p.s, r.value)
            return r.key, false, RangedPart{p.x, p.m, p.a, a}, RangedPart{p.x, p.m, p.a, b}

        default:
            panic("Wrong Category")
        }
    }
}

func (r DecisionRule) getValue() int { return r.value }
func (r DecisionRule) getSign() string { return r.sign }

func ruleToObj(r string) (string, []Rule) {
    s := strings.Split(r, "{")
    if len(s) != 2 { panic("Rule parsing error") }

    key := s[0]
    ruleStr := strings.Replace(s[1], "}", "", 1)
    rulesSplit := strings.Split(ruleStr, ",")

    rules := []Rule{}

    for _, i := range rulesSplit {
        rSplit := strings.Split(i, ":")
        if len(rSplit) == 1 {
            if rSplit[0] == "A"{
                rules = append(rules, AcceptRule{})
            } else if rSplit[0] == "R" {
                rules = append(rules, RejectRule{})
            } else {
                rules = append(rules, GoToRule{rSplit[0]})
            }
            continue
        }

        nextKey := rSplit[1] 
        ruleText := rSplit[0]

        cat := string(ruleText[0])
        sign := string(ruleText[1])
        thresh, err := strconv.Atoi(string(ruleText[2:]))
        if err != nil { panic(err) }

        rules = append(rules, DecisionRule{cat, sign, thresh, nextKey})
    }

    return key, rules
}

func solution1(rulesStr, partsStr []string) int {
    ruleMap := make(map[string][]func(Part)(string, bool))
    for _, i := range rulesStr {
        key, rule := parseRule(i)
        ruleMap[key] = rule
    }
    parts := parseParts(partsStr)
    workingParts := computeRuleResult(parts, ruleMap)

    out := 0
    for _, i := range workingParts {
        out += i.x + i.m + i.a + i.s
    }
    return out
}

func printRangedPart(p RangedPart) {
    fmt.Printf("x: [%d, %d]\n", p.x[0], p.x[len(p.x)-1])
    fmt.Printf("m: [%d, %d]\n", p.m[0], p.m[len(p.m)-1])
    fmt.Printf("a: [%d, %d]\n", p.a[0], p.a[len(p.a)-1])
    fmt.Printf("s: [%d, %d]\n", p.s[0], p.s[len(p.s)-1])
    fmt.Print("\n")
}

func traverseRules(currentPart RangedPart, ruleMap map[string][]Rule, key string) []RangedPart {
    if key == "A" {
        return []RangedPart{currentPart}
    } else if key == "R" {
        return []RangedPart{}
    } 
    out := []RangedPart{}
    if rules, ok := ruleMap[key]; ok {

        var result string
        var pass bool
        var passedPart RangedPart
        var failedPart RangedPart

        for _, rule := range rules {
            result, pass, passedPart, failedPart = rule.execute(currentPart)
            fmt.Printf("%s -> %s (%t)\n", key, result, pass)
            // printRangedPart(currentPart)

            if pass {
                if result == "A" {
                    fmt.Println("ACCEPTED")
                    out = append(out, currentPart)
                } else if result == "R" {
                    // no-op
                } else {
                    out = append(out, traverseRules(passedPart, ruleMap, result)...)
                }
                break
            } else {
                out = append(out, traverseRules(passedPart, ruleMap, result)...)
            }
            currentPart = failedPart
        }

    } else {
        fmt.Printf("%s\n", key)
        panic("Bad key")
    }

    return out
}


func createRange(min, max int) []int {
    out := []int{}
    for i:=min; i<=max; i++ {
        out = append(out, i)
    }
    return out
}


func solution2(rulesStr, partsStr []string) int {
    ruleMap := make(map[string][]Rule)
    for _, i := range rulesStr {
        key, rules := ruleToObj(i)
        ruleMap[key] = rules
        // fmt.Printf("%s -> %s\n", key, rules)
    }

    out := 0
    part := RangedPart{createRange(1, 4000), createRange(1, 4000), createRange(1, 4000), createRange(1, 4000)}
    workingParts := traverseRules(part, ruleMap, "in")
    for _, i := range workingParts {
        out += (len(i.x)-1) * (len(i.m)-1) * (len(i.a)-1) * (len(i.s)-1)
        // printRangedPart(i)
    }

    return out
}


func main() {
	rules, parts := readLines("input2.txt")
	fmt.Printf("Solution 1: %d\n", solution1(rules, parts))
	fmt.Printf("Solution 2: %d\n", solution2(rules, parts))
}

