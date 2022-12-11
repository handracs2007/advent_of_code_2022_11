package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Monkey represents a monkey.
type Monkey struct {
	Items             []int
	CalculateNewWorry func(old int) int
	GetNextMonkey     func(worry int) int
	inspectionCount   int
}

// Receive receives an item and store it into the list of items.
func (monkey *Monkey) Receive(item int) {
	monkey.Items = append(monkey.Items, item)
}

// Inspect returns the current item and remove it from the items bag. Returns -1 if there is no more item.
func (monkey *Monkey) Inspect() int {
	if len(monkey.Items) == 0 {
		return -1
	}

	item := monkey.Items[0]
	monkey.Items = monkey.Items[1:]
	monkey.inspectionCount++
	return item
}

func (monkey *Monkey) InspectionCount() int {
	return monkey.inspectionCount
}

// badMonkeysInAction simulates the bad monkeys in action. The useMod parameter is to determine whether a division
// or module should be performed on the new weight. This is to accommodate both part 1 and part 2 of the puzzle.
func badMonkeysInAction(monkeys []*Monkey, rounds, divFactor int, useMod bool) {
	for round := 0; round < rounds; round++ {
		for i := 0; i < len(monkeys); i++ {
			m := monkeys[i]

			for {
				w := m.Inspect()
				if w == -1 {
					break
				}

				nw := m.CalculateNewWorry(w)
				if useMod {
					nw %= divFactor
				} else {
					nw /= divFactor
				}

				nm := m.GetNextMonkey(nw)

				monkeys[nm].Receive(nw)
			}
		}
	}

	counts := make([]int, len(monkeys), len(monkeys))
	for i := 0; i < len(monkeys); i++ {
		counts[i] = monkeys[i].InspectionCount()
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i] > counts[j]
	})
	fmt.Println(counts[0] * counts[1])
}

func main() {
	// Read the input file.
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("unable to open input file: %s", err)
	}
	defer f.Close()

	var supermod = 1 // This is for part 2
	var currMonkey1 *Monkey
	var currMonkey2 *Monkey
	var monkeys1 []*Monkey
	var monkeys2 []*Monkey
	r := bufio.NewReader(f)
	for {
		l, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Fatalf("failed to read input file: %s", err)
		}

		if err == io.EOF {
			break
		}
		l = strings.TrimSpace(l)

		if strings.HasPrefix(l, "Monkey") {
			// Append a new monkey
			currMonkey1 = &Monkey{
				Items:             []int{},
				CalculateNewWorry: nil,
				GetNextMonkey:     nil,
			}
			monkeys1 = append(monkeys1, currMonkey1)

			currMonkey2 = &Monkey{
				Items:             []int{},
				CalculateNewWorry: nil,
				GetNextMonkey:     nil,
			}
			monkeys2 = append(monkeys2, currMonkey2)
		} else if strings.HasPrefix(l, "Starting items") {
			colonIdx := strings.Index(l, ":")
			strItems := l[colonIdx+2:]
			items := strings.Split(strItems, ", ")

			for _, item := range items {
				val, _ := strconv.Atoi(item)
				currMonkey1.Items = append(currMonkey1.Items, val)
				currMonkey2.Items = append(currMonkey2.Items, val)
			}
		} else if strings.HasPrefix(l, "Operation") {
			operator := l[21]
			factor, err := strconv.Atoi(l[23:])
			if err != nil {
				factor = 0
			}

			switch operator {
			case '+':
				calc := func(old int) int {
					if factor == 0 {
						return old + old
					}

					return old + factor
				}
				currMonkey1.CalculateNewWorry = calc
				currMonkey2.CalculateNewWorry = calc
				break

			case '-':
				calc := func(old int) int {
					if factor == 0 {
						return 0
					}

					return old - factor
				}
				currMonkey1.CalculateNewWorry = calc
				currMonkey2.CalculateNewWorry = calc
				break

			case '*':
				calc := func(old int) int {
					if factor == 0 {
						return old * old
					}

					return old * factor
				}
				currMonkey1.CalculateNewWorry = calc
				currMonkey2.CalculateNewWorry = calc
				break

			case '/':
				calc := func(old int) int {
					if factor == 0 {
						return 1
					}

					return old / factor
				}
				currMonkey1.CalculateNewWorry = calc
				currMonkey2.CalculateNewWorry = calc
				break
			}
		} else if strings.HasPrefix(l, "Test") {
			lastSpaceIdx := strings.LastIndex(l, " ")
			divFactor, _ := strconv.Atoi(l[lastSpaceIdx+1:])
			supermod *= divFactor // Used for part 2

			// Read the if true condition
			l, _ = r.ReadString('\n')
			l = strings.TrimSpace(l)
			lastSpaceIdx = strings.LastIndex(l, " ")
			trueMonkey, _ := strconv.Atoi(l[lastSpaceIdx+1:])

			// Read the if false condition
			l, _ = r.ReadString('\n')
			l = strings.TrimSpace(l)
			lastSpaceIdx = strings.LastIndex(l, " ")
			falseMonkey, _ := strconv.Atoi(l[lastSpaceIdx+1:])

			get := func(worry int) int {
				if worry%divFactor == 0 {
					return trueMonkey
				}

				return falseMonkey
			}
			currMonkey1.GetNextMonkey = get
			currMonkey2.GetNextMonkey = get
		}
	}

	badMonkeysInAction(monkeys1, 20, 3, false)
	badMonkeysInAction(monkeys2, 10000, supermod, true)
}
