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

type plot struct {
	plant  rune
	region int
}

type state struct {
	fname      string
	plots      map[coord]*plot
	nextRegion int
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

	s.plots = map[coord]*plot{}
	s.nextRegion = 1

	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x, plant := range line {
			s.plots[coord{x: x, y: y}] = &plot{plant: plant}
		}

		y++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/12/tiny",
		"input/12/small",
		"input/12/sample",
		"input/12/input",
	} {
		s := readFile(fname)
		if s.plots == nil {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func neighbors(c coord) []coord {
	return []coord{
		coord{c.x - 1, c.y},
		coord{c.x + 1, c.y},
		coord{c.x, c.y - 1},
		coord{c.x, c.y + 1},
	}
}

func (s *state) floodFill(c coord) int {
	fringe := map[coord]bool{c: true}
	for len(fringe) > 0 {
		for c := range fringe {
			plot := s.plots[c]
			if plot.region == 0 {
				plot.region = s.nextRegion
				for _, n := range neighbors(c) {
					if p1, ok := s.plots[n]; ok && p1.region == 0 && p1.plant == plot.plant {
						fringe[n] = true
					}
				}
			}
			delete(fringe, c)
		}
	}
	perimeter, area := 0, 0
	fringe = map[coord]bool{c: true}
	seen := map[coord]bool{}
	for len(fringe) > 0 {
		found := false
		for c := range fringe {
			if !seen[c] {
				for _, n := range neighbors(c) {
					if p1, ok := s.plots[n]; ok && p1.region == s.nextRegion {
						fringe[n] = true
					} else {
						// no neigbor on that side - add perimeter
						perimeter++
					}
				}
				area++
				found = true
				seen[c] = true
				break
			}
		}
		if !found {
			break
		}
	}

	s.nextRegion++
	return perimeter * area
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)
		part1 := 0
		for coord, plot := range s.plots {
			if plot.region == 0 {
				part1 += s.floodFill(coord)
			}
		}

		fmt.Println("fence price: ", part1)
	}
}
