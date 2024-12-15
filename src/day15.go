package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type coord struct {
	x int
	y int
}

type state struct {
	fname string
	size  coord
	walls map[coord]bool
	boxes map[coord]bool
	robot coord

	movements string
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
	s.boxes = map[coord]bool{}

	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		for x, v := range line {
			if v == '#' {
				s.walls[coord{x, y}] = true
			} else if v == '@' {
				s.robot = coord{x, y}
			} else if v == 'O' {
				s.boxes[coord{x, y}] = true
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
	for scanner.Scan() {
		s.movements += scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	if s.movements == "" {
		log.Fatal("no movement")
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/15/small",
		"input/15/sample",
		"input/15/input",
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

func dbg(s state, msg string) {
	fmt.Println(msg + ":")
	errs := 0
	for y := 0; y < s.size.y; y++ {
		for x := 0; x < s.size.x; x++ {
			c := coord{x, y}
			if s.robot == c && s.walls[c] && s.boxes[c] {
				fmt.Print("!")
				errs++
			} else if s.robot == c && s.walls[c] {
				fmt.Print("*")
				errs++
			} else if s.robot == c && s.boxes[c] {
				fmt.Print("?")
				errs++
			} else if s.boxes[c] && s.walls[c] {
				fmt.Print("X")
				errs++
			} else if s.robot == c {
				fmt.Print("@")
			} else if s.walls[c] {
				fmt.Print("#")
			} else if s.boxes[c] {
				fmt.Print("O")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
	if errs > 0 {
		log.Fatal(strconv.Itoa(errs) + " errors!")
	}
}

var mdir = map[rune]coord{
	'<': coord{-1, 0},
	'^': coord{0, -1},
	'>': coord{1, 0},
	'v': coord{0, 1},
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)

		for _, m := range s.movements {
			dir, ok := mdir[m]
			if !ok {
				log.Fatal("unexpected movement " + string(m))
			}
			c := plus(s.robot, dir)
			if s.boxes[c] {
				firstBox := c
				for {
					c = plus(c, dir)
					if s.walls[c] {
						break
					} else if !s.boxes[c] {
						s.boxes[c] = true
						delete(s.boxes, firstBox)
						s.robot = firstBox
						break
					}
				}
			} else if !s.walls[c] {
				s.robot = c
			}
		}

		sum := 0
		for b := range s.boxes {
			sum += b.y*100 + b.x
		}

		fmt.Println(sum)
	}
}
