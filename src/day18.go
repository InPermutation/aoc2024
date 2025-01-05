package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type pos struct {
	x int
	y int
}

type state struct {
	fname     string
	exit      pos
	corrupted []pos
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

	if fname == "input/18/sample" {
		s.exit = pos{6, 6}
	} else {
		s.exit = pos{70, 70}
	}

	for scanner.Scan() {
		line := scanner.Text()
		toks := strings.Split(line, ",")
		if len(toks) != 2 {
			log.Fatal("invalid line", line)
		}
		x, err := strconv.Atoi(toks[0])
		if err != nil {
			log.Fatal(err)
		}
		y, err := strconv.Atoi(toks[1])
		if err != nil {
			log.Fatal(err)
		}
		s.corrupted = append(s.corrupted, pos{x, y})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/18/sample",
		"input/18/input",
	} {
		s := readFile(fname)
		if len(s.corrupted) == 0 {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func (s *state) neighbors(p pos, corrupted map[pos]bool) []pos {
	rv := []pos{}
	ns := []pos{{p.x - 1, p.y}, {p.x + 1, p.y}, {p.x, p.y - 1}, {p.x, p.y + 1}}
	for _, n := range ns {
		if corrupted[n] {
			continue
		}
		if n.x >= 0 && n.x <= s.exit.x {
			if n.y >= 0 && n.y <= s.exit.y {
				rv = append(rv, n)
			}
		}
	}
	return rv
}

func (s *state) firstBytes() int {
	switch s.fname {
	case "input/18/sample":
		return 12
	case "input/18/input":
		return 1024
	default:
		panic(s.fname)
	}
}

func (s *state) firstCorrupted() map[pos]bool {
	corrupted := map[pos]bool{}
	for _, c := range s.corrupted[:s.firstBytes()] {
		corrupted[c] = true
	}

	return corrupted
}

func (s *state) Part1() int {
	origin := pos{0, 0}
	m := map[pos]int{
		origin: 0,
	}

	corrupted := s.firstCorrupted()

	fringe := []pos{origin}
	for len(fringe) > 0 {
		p := fringe[0]
		stepCost := m[p] + 1
		fringe = fringe[1:]
		for _, n := range s.neighbors(p, corrupted) {
			if currCost, ok := m[n]; !ok || currCost > stepCost {
				fringe = append(fringe, n)
				m[n] = stepCost
			}
		}
	}

	return m[s.exit]
}

func (s *state) Part2() pos {
	origin := pos{0, 0}

	// skip the firstCorrupted; we know from Part1 they are possible
	corrupted := s.firstCorrupted()
	// iteratively widen the solution
	for _, c := range s.corrupted[s.firstBytes():] {
		corrupted[c] = true
		// don't actually care about cost
		m := map[pos]bool{
			origin: true,
		}
		fringe := []pos{origin}

		for len(fringe) > 0 {
			p := fringe[0]
			fringe = fringe[1:]
			for _, n := range s.neighbors(p, corrupted) {
				if !m[n] {
					fringe = append(fringe, n)
					m[n] = true
				}
			}
		}
		if !m[s.exit] {
			return c
		}
	}

	panic("did not fail")
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)
		fmt.Println("part 1", s.Part1())
		fpt := s.Part2()
		fmt.Printf("part 2 (%v,%v)\n", fpt.x, fpt.y)
	}
}
