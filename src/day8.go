package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type pos struct {
	x int
	y int
}

type state struct {
	fname  string
	nodes  map[rune][]pos
	height int
	width  int
}

func readFile(fname string) state {
	file, err := os.Open(fname)
	if err != nil {
		log.Print(err)
		return state{fname: fname}
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	y := 0
	nodes := map[rune][]pos{}
	var height, width int
	for scanner.Scan() {
		line := scanner.Text()
		width = len(line)
		for x, f := range line {
			if f != '.' {
				nodes[f] = append(nodes[f], pos{x, y})
			}
		}

		y++
		// establish bounds:
		height, width = y, len(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return state{fname: fname, nodes: nodes, height: height, width: width}
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/08/sample",
		"input/08/input",
	} {
		s := readFile(fname)
		if s.width == 0 {
			continue
		}
		rval = append(rval, s)
	}
	return

}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)
		antinodes := map[pos]bool{}
		harmonics := map[pos]bool{}
		for _, ps := range s.nodes {
			for i, p := range ps {
				for _, p2 := range ps[i+1:] {
					for _, cand := range resonate(p, p2) {
						if boundCheck(cand, s.width, s.height) {
							antinodes[cand] = true
						}
					}
					for _, cand := range harmonize(p, p2, s.width, s.height) {
						if boundCheck(cand, s.width, s.height) {
							harmonics[cand] = true
						}
					}
				}
			}
		}

		// Part 1
		fmt.Println("total unique antinode locations: ", len(antinodes))
		// Part 2
		fmt.Println("total unique harmonic locations: ", len(harmonics))
	}
}

func boundCheck(cand pos, width, height int) bool {
	return cand.x >= 0 && cand.y >= 0 && cand.x < width && cand.y < height
}

func resonate(a, b pos) (rval []pos) {
	d := sub(a, b)
	return []pos{
		add(a, d),
		sub(b, d),
	}
}

func add(a, b pos) pos {
	return pos{
		a.x + b.x,
		a.y + b.y,
	}
}
func sub(a, b pos) pos {
	return pos{
		a.x - b.x,
		a.y - b.y,
	}
}

func harmonize(a, b pos, width, height int) (rval []pos) {
	d := sub(a, b)
	for p := b; boundCheck(p, width, height); p = add(p, d) {
		rval = append(rval, p)
	}
	for p := a; boundCheck(p, width, height); p = sub(p, d) {
		rval = append(rval, p)
	}

	return
}
