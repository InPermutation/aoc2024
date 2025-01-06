package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type coord struct {
	x int
	y int
}

type state struct {
	fname string
	size  coord
	walls map[coord]bool
	start coord
	exit  coord
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
	s.walls = map[coord]bool{}

	y := 0
	for scanner.Scan() {
		line := scanner.Text()

		for x, v := range line {
			if v == '#' {
				s.walls[coord{x, y}] = true
			} else if v == 'S' {
				s.start = coord{x, y}
			} else if v == 'E' {
				s.exit = coord{x, y}
			}
			if x > s.size.x {
				s.size.x = x
			}
		}

		s.size.y = y
		y++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	s.size.x++
	s.size.y++

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/20/sample",
		"input/20/input",
	} {
		s := readFile(fname)
		if len(s.walls) == 0 {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func plus(a, b coord) coord {
	return coord{a.x + b.x, a.y + b.y}
}

var (
	LEFT  = coord{-1, 0}
	RIGHT = coord{1, 0}
	UP    = coord{0, -1}
	DOWN  = coord{0, 1}

	LRUD = []coord{LEFT, RIGHT, UP, DOWN}
)

func (s *state) neighbors(k coord) []coord {
	rv := []coord{}
	for _, dir := range LRUD {
		pos := plus(k, dir)
		if s.walls[pos] || pos.x < 0 || pos.y < 0 || pos.x >= s.size.x || pos.y >= s.size.y {
			continue
		}
		rv = append(rv, pos)
	}

	return rv
}

func (s *state) CostFrom(source coord) map[coord]int {
	frontier := map[coord]bool{source: true}
	scores := map[coord]int{source: 0}
	for len(frontier) > 0 {
		for k := range frontier {
			delete(frontier, k)
			some := false
			curr, ok := scores[k]
			if !ok {
				log.Fatal("no score for curr", k)
			}
			for _, v := range s.neighbors(k) {
				sc, ok := scores[v]
				if !ok || sc > curr+1 {
					frontier[v] = true
					scores[v] = curr + 1
					some = true
				}
			}
			if some {
				break
			}
		}
	}

	return scores
}

func (s *state) Cheats(maxCost int) map[int]int {
	fromStart := s.CostFrom(s.start)
	fromExit := s.CostFrom(s.exit)
	base := fromStart[s.exit]
	cheats := map[int]int{}

	for x := 0; x < s.size.x; x++ {
		for y := 0; y < s.size.y; y++ {
			start := coord{x, y}
			if s.walls[start] {
				continue
			}
			fmStart, ok := fromStart[start]
			if !ok {
				continue
			}
			for dy := -maxCost; dy <= maxCost; dy++ {
				for dx := -maxCost; dx <= maxCost; dx++ {
					cost := abs(dx) + abs(dy)
					if cost > maxCost {
						continue
					}
					end := plus(start, coord{dx, dy})
					if s.walls[end] {
						continue
					}
					fmExit, ok := fromExit[end]
					if !ok {
						continue
					}
					time := cost + fmStart + fmExit
					cheats[base-time]++
				}
			}
		}
	}

	return cheats
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)
		fromStart := s.CostFrom(s.start)
		base := fromStart[s.exit]

		cheats := s.Cheats(2)

		c := 0
		for i := 100; i <= base; i++ {
			c += cheats[i]
		}

		// Part 1
		fmt.Println("There are", c, "2ps cheats that save at least 100ps.")

		cheats = s.Cheats(20)

		c = 0
		for i := 100; i <= base; i++ {
			c += cheats[i]
		}

		// Part 2
		fmt.Println("There are", c, "≤20ps cheats that save at least 100ps.")
	}
}
