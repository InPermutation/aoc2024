package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type freq rune

func (f freq) String() string {
	return "'" + string(f) + "'"
}

type pos struct {
	x int
	y int
}

func (p pos) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
}

func main() {
	for _, fname := range []string{
		"input/08/sample",
		"input/08/input",
	} {
		fmt.Println(fname)
		file, err := os.Open(fname)
		if err != nil {
			log.Print(err)
			continue
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)

		y := 0
		nodes := map[freq][]pos{}
		var height, width int
		for scanner.Scan() {
			line := scanner.Text()
			width = len(line)
			for x, f := range line {
				if f != '.' {
					nodes[freq(f)] = append(nodes[freq(f)], pos{x, y})
				}
			}

			y++
			// establish bounds:
			height, width = y, len(line)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		antinodes := map[pos]bool{}
		harmonics := map[pos]bool{}
		for _, ps := range nodes {
			for i, p := range ps {
				for _, p2 := range ps[i+1:] {
					for _, cand := range resonate(p, p2) {
						if boundCheck(cand, width, height) {
							antinodes[cand] = true
						}
					}
					for _, cand := range harmonize(p, p2, width, height) {
						if boundCheck(cand, width, height) {
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
	dx := a.x - b.x
	dy := a.y - b.y
	return []pos{
		pos{b.x - dx, b.y - dy},
		pos{a.x + dx, a.y + dy},
	}
}

func add(a, b pos) pos {
	return pos{
		a.x + b.x,
		a.y + b.y,
	}
}

func harmonize(a, b pos, width, height int) (rval []pos) {
	d := pos{a.x - b.x, a.y - b.y}
	for p := b; boundCheck(p, width, height); p = add(p, d) {
		rval = append(rval, p)
	}
	d = pos{-d.x, -d.y}
	for p := a; boundCheck(p, width, height); p = add(p, d) {
		rval = append(rval, p)
	}

	return
}
