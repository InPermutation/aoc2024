package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type pos struct {
	x                int
	y                int
	reachableSummits map[*pos]bool
	rating           int
}
type state struct {
	fname   string
	heights [10][]*pos
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

	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x, v := range line {
			h := int(v - '0')
			s.heights[h] = append(s.heights[h], &pos{
				x:                x,
				y:                y,
				reachableSummits: map[*pos]bool{},
			})
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
		"input/10/sample",
		"input/10/larger",
		"input/10/input",
	} {
		s := readFile(fname)
		if len(s.heights[0]) == 0 {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func neighbors(a, b *pos) bool {
	if a.x == b.x {
		return a.y == b.y-1 || a.y == b.y+1
	} else if a.y == b.y {
		return a.x == b.x-1 || a.x == b.x+1
	}
	return false
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)

		for _, summit := range s.heights[9] {
			summit.reachableSummits = map[*pos]bool{summit: true}
			summit.rating = 1
		}

		sum, rsum := 0, 0
		for i := 8; i >= 0; i-- {
			for _, p := range s.heights[i] {
				for _, t := range s.heights[i+1] {
					if neighbors(p, t) {
						for s := range t.reachableSummits {
							p.reachableSummits[s] = true
						}
						p.rating += t.rating
					}
				}
				if i == 0 {
					sum += len(p.reachableSummits)
					rsum += p.rating
				}
			}
		}

		// Part 1
		fmt.Println("sum of scores: ", sum)
		// Part 2
		fmt.Println("sum of ratings: ", rsum)
	}
}
