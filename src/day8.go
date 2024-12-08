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
		for _, ps := range nodes {
			for i, p := range ps {
				for _, p2 := range ps[i+1:] {
					for _, cand := range resonate(p, p2) {
						if boundCheck(cand, width, height) {
							antinodes[cand] = true
						}
					}
				}
			}
		}

		// Part 1
		fmt.Println("total unique antinode locations: ", len(antinodes))
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
