package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type state struct {
	fname string
	codes [5]string
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

	i := 0
	for scanner.Scan() {
		s.codes[i] = scanner.Text()
		i++
	}
	if i != 5 {
		log.Fatal("too few lines ", i)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/21/sample",
		"input/21/input",
	} {
		s := readFile(fname)
		if len(s.codes[0]) == 0 {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)
		sum := 0
		for _, v := range s.codes {
			numericPart, err := strconv.Atoi(v[:len(v)-1])
			if err != nil {
				log.Fatal(err)
			}
			d := driveNumeric(v)
			d = driveDirectional(d)
			d = driveDirectional(d)
			sum += len(d) * numericPart
		}
		fmt.Println()
		fmt.Println("sum complexity", sum)
	}
}

var directional map[rune]coord = map[rune]coord{
	//  ^A
	// <v>
	'^': coord{1, 0},
	'A': coord{2, 0},

	'<': coord{0, 1},
	'v': coord{1, 1},
	'>': coord{2, 1},
}

func driveDirectional(s string) string {
	return drivePad(directional, s)
}

type coord struct {
	x int
	y int
}

var numpad map[rune]coord = map[rune]coord{
	// 789
	// 456
	// 123
	//  0A
	'7': coord{0, 0},
	'8': coord{1, 0},
	'9': coord{2, 0},
	'4': coord{0, 1},
	'5': coord{1, 1},
	'6': coord{2, 1},

	'1': coord{0, 2},
	'2': coord{1, 2},
	'3': coord{2, 2},

	'0': coord{1, 3},
	'A': coord{2, 3},
}

func driveNumeric(s string) string {
	return drivePad(numpad, s)
}

func drivePad(pad map[rune]coord, s string) string {
	_, isDirectional := pad['<']
	pos := pad['A']
	sb := strings.Builder{}

	for _, c := range s {
		next, ok := pad[c]
		if !ok {
			log.Fatal("unknown digit ", string(c), " in ", s)
		}
		diff := coord{next.x - pos.x, next.y - pos.y}
		if isDirectional {
			// don't hover X
			if pos.y == 0 && next.x == 0 {
				for diff.y > 0 {
					sb.WriteRune('v')
					diff.y--
				}
			}
			if pos.x == 1 && pos.y == 1 {
				for diff.x > 0 {
					sb.WriteRune('>')
					diff.x--
				}
			}

		} else {
			// don't hover X
			if pos.x == 0 && next.y == 3 {
				for diff.x > 0 {
					sb.WriteRune('>')
					diff.x--
				}
			}
			if pos.y == 3 && next.x == 0 {
				for diff.y < 0 {
					sb.WriteRune('^')
					diff.y++
				}
			}
		}

		for diff.x < 0 {
			sb.WriteRune('<')
			diff.x++
		}
		for diff.y > 0 {
			sb.WriteRune('v')
			diff.y--
		}
		for diff.x > 0 {
			sb.WriteRune('>')
			diff.x--
		}
		for diff.y < 0 {
			sb.WriteRune('^')
			diff.y++
		}
		sb.WriteRune('A')
		pos = next
	}

	return sb.String()
}
