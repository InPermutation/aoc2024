package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	for _, fname := range []string{
		"input/04/sample",
		"input/04/input",
	} {
		fmt.Println(fname)

		file, err := os.Open(fname)
		if err != nil {
			log.Print(err)
			continue
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)

		lines := []string{}
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				lines = append(lines, line)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		// Part 1:
		found := 0

		safeGet := func(x, y int) byte {
			if y < 0 || y >= len(lines) {
				return '?'
			}
			line := lines[y]
			if x < 0 || x >= len(line) {
				return '?'
			}
			return line[x]
		}

		for y, line := range lines {
			for x, ch := range line {
				if ch != 'X' {
					continue
				}
				const Target = "XMAS"
				dir := func(dx, dy int) int {
					for o := 0; o < 4; o++ {
						it := safeGet(x+(o*dx), y+(o*dy))
						if it != Target[o] {
							return 0
						}
					}
					return 1
				}

				found += (dir(0, 1) + dir(0, -1) +
					dir(1, 0) + dir(-1, 0) +
					dir(1, 1) + dir(-1, -1) + dir(1, -1) + dir(-1, 1))
			}
		}

		fmt.Print("Found XMAS: ")
		fmt.Println(found)
	}
}
