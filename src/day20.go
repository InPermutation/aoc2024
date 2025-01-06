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

func plus(a, b coord) coord {
	return coord{a.x + b.x, a.y + b.y}
}

func (s *state) neighbors(k coord) []coord {
	rv := []coord{}
	for _, dir := range []coord{
		coord{-1, 0},
		coord{1, 0},
		coord{0, -1},
		coord{0, 1},
	} {
		pos := plus(k, dir)
		if s.walls[pos] || pos.x < 0 || pos.y < 0 || pos.x >= s.size.x || pos.y >= s.size.y {
			continue
		}
		rv = append(rv, pos)
	}

	return rv
}

func (s *state) Cost() int {

	frontier := map[coord]bool{s.start: true}
	scores := map[coord]int{s.start: 0}
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

	return scores[s.exit]
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)

		// Part 0
		base := s.Cost()

		cheats := map[int]int{}
		for wall := range s.walls {
			s.walls[wall] = false

			wo := s.Cost()
			cheats[base-wo]++

			s.walls[wall] = true
		}

		c := 0
		for i := 100; i <= base; i++ {
			c += cheats[i]
		}

		fmt.Println("There are", c, "cheats that save at least 100 picoseconds.")
	}
}
