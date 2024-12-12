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
		"input/12/eshape",
		"input/12/mobius",
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
		coord{c.x - 1, c.y}, // LEFT
		coord{c.x + 1, c.y}, // RIGHT
		coord{c.x, c.y - 1}, // UP
		coord{c.x, c.y + 1}, // DOWN
	}
}

func (s *state) floodFill(c coord) map[coord]*plot {
	region := map[coord]*plot{}
	fringe := map[coord]bool{c: true}
	for len(fringe) > 0 {
		for c := range fringe {
			plot := s.plots[c]
			region[c] = plot
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
	s.nextRegion++
	return region
}

func rackRate(plots map[coord]*plot) int {
	perimeter, area := 0, 0
	seen := map[coord]bool{}
	for c := range plots {
		if !seen[c] {
			for _, n := range neighbors(c) {
				if _, ok := plots[n]; !ok {
					// no neigbor on that side - add perimeter
					perimeter++
				}
			}
			area++
			seen[c] = true
		}
	}

	return perimeter * area
}

func discount(plots map[coord]*plot) int {
	sides, area := 0, 0
	seen := map[coord]bool{}
	// only count the topmost/leftmost segments of each side
	for c := range plots {
		if !seen[c] {
			ns := neighbors(c)
			l, r, u, d := plots[ns[0]] != nil, plots[ns[1]] != nil, plots[ns[2]] != nil, plots[ns[3]] != nil
			uns := neighbors(ns[2])
			lns := neighbors(ns[0])
			if !l {
				// am i the top?
				if !u {
					// convex corner
					sides++
				} else if plots[uns[0]] != nil {
					//concave corner
					sides++
				}
			}
			if !r {
				// am i the top?
				if !u {
					// convex corner
					sides++
				} else if plots[uns[1]] != nil {
					// concave corner
					sides++
				}
			}
			if !u {
				// am i the left?
				if !l {
					// convex corner
					sides++
				} else if plots[lns[2]] != nil {
					// concave corner
					sides++
				}
			}
			if !d {
				// am i the left?
				if !l {
					// convex corner
					sides++
				} else if plots[lns[3]] != nil {
					//concave corner
					sides++
				}
			}

			area++
			seen[c] = true
		}
	}

	return sides * area
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)
		part1, part2 := 0, 0
		for coord, plot := range s.plots {
			if plot.region == 0 {
				region := s.floodFill(coord)
				part1 += rackRate(region)
				part2 += discount(region)
			}
		}

		fmt.Println("fence price: ", part1)
		fmt.Println("discount price: ", part2)
	}
}
