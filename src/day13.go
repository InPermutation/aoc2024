package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type coord struct {
	x int
	y int
}

type machine struct {
	A     coord
	B     coord
	Prize coord
}

type state struct {
	fname    string
	machines []machine
}

func readFile(fname string) state {
	s := state{fname: fname}
	file, err := os.Open(fname)
	if err != nil {
		log.Print(err)
		return s
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var A, B coord
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.FieldsFunc(line, func(c rune) bool {
			return strings.IndexRune(":=,+ ", c) >= 0
		})
		if len(fields) < 3 {
			continue
		}
		x, err := strconv.Atoi(fields[len(fields)-3])
		if err != nil {
			fmt.Println(fields)
			log.Fatal(err)
		}
		y, err := strconv.Atoi(fields[len(fields)-1])
		if err != nil {
			fmt.Println(fields)
			log.Fatal(err)
		}
		if strings.HasPrefix(line, "Button ") {
			if line[7] == 'A' {
				A = coord{x, y}
			} else if line[7] == 'B' {
				B = coord{x, y}
			} else {
				log.Fatal("unknown button " + line)
			}
		} else if strings.HasPrefix(line, "Prize: ") {
			s.machines = append(s.machines, machine{A, B, coord{x, y}})
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/13/sample",
		"input/13/input",
	} {
		s := readFile(fname)
		if s.machines == nil {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func sub(a, b coord) coord {
	return coord{a.x - b.x, a.y - b.y}
}

func plus(a, b coord) coord {
	return coord{a.x + b.x, a.y + b.y}
}

func times(a coord, n int) coord {
	return coord{a.x * n, a.y * n}
}

func minTokens(m machine) (rval int) {
	for a := 0; a < 100; a++ {
		cost := 3 * a
		if rval != 0 && cost >= rval {
			break
		}
		for b := 0; b < 100; b++ {
			cost := 3*a + b
			if rval != 0 && cost >= rval {
				break
			}
			if plus(times(m.A, a), times(m.B, b)) == m.Prize {
				if rval == 0 || rval > cost {
					rval = cost
				}
			}
		}
	}
	return
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)

		// Part 1
		sum := 0
		for _, machine := range s.machines {
			sum += minTokens(machine)
		}
		fmt.Println()
		fmt.Println("min tokens: ", sum)
	}
}
