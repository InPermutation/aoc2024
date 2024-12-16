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

type reindeer struct {
	pos coord
	dir coord
}

type state struct {
	fname    string
	size     coord
	walls    map[coord]bool
	reindeer reindeer
	exit     coord
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
				s.reindeer = reindeer{
					pos: coord{x, y},
					dir: coord{1, 0},
				}
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
		"input/16/sample",
		"input/16/sample2",
		"input/16/input",
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

func turns(a, b coord) int {
	if a == b {
		return 0
	}
	if a.x == b.x {
		return 2
	}
	if a.y == b.y {
		return 2
	}
	return 1
}

func (s *state) neighbors(k reindeer, curr int) map[reindeer]int {
	rv := map[reindeer]int{}
	for _, dir := range []coord{
		coord{-1, 0},
		coord{1, 0},
		coord{0, -1},
		coord{0, 1},
	} {
		pos := plus(k.pos, dir)
		if s.walls[pos] || pos.x < 0 || pos.y < 0 || pos.x >= s.size.x || pos.y >= s.size.y {
			continue
		}
		r := reindeer{pos: pos, dir: dir}
		s := 1000*turns(dir, k.dir) + 1 + curr
		rv[r] = s
	}

	return rv
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)

		frontier := map[reindeer]bool{s.reindeer: true}
		scores := map[coord]int{s.reindeer.pos: 0}
		for len(frontier) > 0 {
			for k := range frontier {
				delete(frontier, k)
				some := false
				curr, ok := scores[k.pos]
				if !ok {
					log.Fatal("no score for curr", k)
				}
				for r, s := range s.neighbors(k, curr) {
					sc, ok := scores[r.pos]
					if !ok || sc > s {
						frontier[r] = true
						scores[r.pos] = s
						some = true
					}
				}
				if some {
					break
				}
			}
		}

		// Part 1
		fmt.Println("min score: ", scores[s.exit])
	}
}
