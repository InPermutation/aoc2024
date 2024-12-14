package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type coord struct {
	x int
	y int
}

type robot struct {
	p coord
	v coord
}

type state struct {
	fname  string
	size   coord
	robots []robot
}

func parseInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return v
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

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.FieldsFunc(line, func(c rune) bool {
			return strings.IndexRune("pv=, ", c) >= 0
		})
		if len(fields) != 4 {
			log.Fatal(line)
		}
		x := parseInt(fields[0])
		y := parseInt(fields[1])
		vx := parseInt(fields[2])
		vy := parseInt(fields[3])
		s.robots = append(s.robots, robot{coord{x, y}, coord{vx, vy}})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if fname == "input/14/sample" {
		s.size = coord{11, 7}
	} else {
		s.size = coord{101, 103}
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/14/sample",
		"input/14/input",
	} {
		s := readFile(fname)
		if s.robots == nil {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func plus(a, b coord) coord {
	return coord{a.x + b.x, a.y + b.y}
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)

		// Part 1
		var q [4]int
		for _, r := range s.robots {
			x := (r.p.x + r.v.x*100) % s.size.x
			y := (r.p.y + r.v.y*100) % s.size.y

			if x < 0 {
				x += s.size.x
			}
			if y < 0 {
				y += s.size.y
			}

			if x < s.size.x/2 && y < s.size.y/2 {
				q[0]++
			} else if x < s.size.x/2 && y > s.size.y/2 {
				q[1]++
			} else if x > s.size.x/2 && y > s.size.y/2 {
				q[2]++
			} else if x > s.size.x/2 && y < s.size.y/2 {
				q[3]++
			}
		}
		sum := 1
		for _, c := range q {
			sum *= c
		}
		fmt.Println(q)
		fmt.Println("safety factor: ", sum)
	}
}
