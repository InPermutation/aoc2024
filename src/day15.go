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
		"input/15/smalldbl",
		"input/15/sample",
		"input/15/input",
		"input/15/m283xvt",
	} {
		s := readFile(fname)
		if len(s.walls) == 0 {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func doubled(s state) (ds state) {
	ds.fname = s.fname
	ds.size = coord{s.size.x * 2, s.size.y}
	ds.movements = s.movements
	ds.walls = map[coord]bool{}
	for w := range s.walls {
		ds.walls[coord{w.x * 2, w.y}] = true
		ds.walls[coord{w.x*2 + 1, w.y}] = true
	}
	ds.boxes = map[coord]bool{}
	for b := range s.boxes {
		ds.boxes[coord{b.x * 2, b.y}] = true
	}
	ds.robot = coord{s.robot.x * 2, s.robot.y}
	return
}

func plus(a, b coord) coord {
	return coord{a.x + b.x, a.y + b.y}
}

func dbg(s state, msg string, doubled bool) {
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
				if doubled {
					fmt.Print("[")
				} else {
					fmt.Print("O")
				}
			} else if doubled && s.boxes[coord{c.x - 1, c.y}] {
				fmt.Print("]")
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

		ds := doubled(s)

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

		fmt.Println("Part 1: ", sum)

		for _, m := range ds.movements {
			dir, ok := mdir[m]
			if !ok {
				log.Fatal("unexpected movement " + string(m))
			}
			c := plus(ds.robot, dir)
			if dir.x == 0 {
				ok := true
				if ds.walls[c] {
					continue
				}
				var row []coord
				var allb []coord
				c2 := coord{c.x - 1, c.y}
				if ds.boxes[c] {
					row = append(row, c)
				} else if ds.boxes[c2] {
					row = append(row, c2)
				}
				for ok && len(row) > 0 {
					var row2 []coord
					for _, b := range row {
						allb = append(allb, b)

						bl := plus(coord{b.x - 1, b.y}, dir)
						bn := plus(b, dir)
						br := plus(coord{b.x + 1, b.y}, dir)
						if ds.walls[bn] || ds.walls[br] {
							ok = false
							break
						}
						if ds.boxes[bn] {
							row2 = append(row2, bn)
						}
						if ds.boxes[br] {
							row2 = append(row2, br)
						}
						if ds.boxes[bl] {
							row2 = append(row2, bl)
						}
					}
					row = row2
				}
				if ok {
					ds.robot = plus(ds.robot, dir)
					// can be dupes, so let's do this in 2 phases
					for _, b := range allb {
						delete(ds.boxes, b)
					}
					for _, b := range allb {
						ds.boxes[plus(b, dir)] = true
					}
				}
			} else if dir.x == -1 {
				c2 := plus(c, dir)
				if ds.boxes[c] {
					log.Fatal(c, c2, dir)
				}
				if ds.boxes[c2] {
					firstBox := c2
					for {
						c2, c = plus(plus(c2, dir), dir), plus(plus(c, dir), dir)
						if ds.walls[c] {
							break
						}
						if !ds.boxes[c2] {
							// move everything 1
							for firstBox != c2 {
								delete(ds.boxes, firstBox)
								firstBox = plus(firstBox, dir)
								ds.boxes[firstBox] = true
								firstBox = plus(firstBox, dir)
							}
							ds.robot = plus(ds.robot, dir)

							break
						}
					}
				} else if !ds.walls[c] {
					ds.robot = c
				}
			} else if dir.x == 1 {
				movb := []coord{}
				for ds.boxes[c] && !ds.walls[c] {
					movb = append(movb, c)
					c = plus(c, dir)
					c = plus(c, dir)
				}
				if !ds.walls[c] {
					ds.robot = plus(ds.robot, dir)
					for _, b := range movb {
						delete(ds.boxes, b)
						ds.boxes[plus(b, dir)] = true
					}
				}
			} else {
				log.Fatal("unknown dir ", dir)
			}
		}
		sum = 0
		for b := range ds.boxes {
			sum += b.y*100 + b.x
		}

		fmt.Println("Part 2:", sum)

	}
}
