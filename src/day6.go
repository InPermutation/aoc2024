package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type direction rune
type pos struct {
	x   int
	y   int
	dir direction
}

type state struct {
	pos
	board [][]rune

	distinct int
	dejaVu   map[pos]bool
}

func NewStateFromFile(fname string) (rval *state) {
	fmt.Println("NewStateFromFile", fname)

	file, err := os.Open(fname)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	rval = &state{
		board:  [][]rune{},
		dejaVu: map[pos]bool{},
	}

	for scanner.Scan() {
		line := scanner.Text()
		bytes := make([]rune, len(line))

		for i, b := range line {
			bytes[i] = b
			if b == '^' {
				rval.x = i
				rval.y = len(rval.board)
				rval.dir = direction(b)
			}
		}

		rval.board = append(rval.board, bytes)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return
}
func (s state) something_in_front() bool {
	switch s.dir {
	case '^':
		return s.y > 0 && s.board[s.y-1][s.x] == '#'
	case '<':
		return s.x > 0 && s.board[s.y][s.x-1] == '#'
	case 'v':
		return s.y+1 < len(s.board) && s.board[s.y+1][s.x] == '#'
	case '>':
		return s.x+1 < len(s.board[s.y]) && s.board[s.y][s.x+1] == '#'
	default:
		fmt.Println("err", s.dir)
		panic(s)
	}
}
func (p *pos) turn90() {
	switch p.dir {
	case '^':
		p.dir = '>'
	case '>':
		p.dir = 'v'
	case 'v':
		p.dir = '<'
	case '<':
		p.dir = '^'
	default:
		panic(*p)
	}
}
func (s *state) step() bool {
	if s.dejaVu[s.pos] {
		return false
	}
	s.dejaVu[s.pos] = true

	if s.board[s.y][s.x] != 'X' {
		s.distinct++
	}
	s.board[s.y][s.x] = 'X'
	switch s.dir {
	case '^':
		s.y--
		return s.y >= 0
	case '>':
		s.x++
		return s.x < len(s.board[s.y])
	case 'v':
		s.y++
		return s.y < len(s.board)
	case '<':
		s.x--
		return s.x >= 0
	default:
		panic(*s)
	}
}
func (s *state) DeepCopy() (rval *state) {
	rval = &state{
		pos:   s.pos,
		board: make([][]rune, len(s.board)),
		// Doesn't copy dejaVu:
		dejaVu: map[pos]bool{},
	}

	for y, row0 := range s.board {
		row := make([]rune, len(row0))

		for x, o := range row0 {
			row[x] = o
		}

		rval.board[y] = row
	}

	return rval
}

func main() {
	for _, fname := range []string{
		"input/06/sample",
		"input/06/input",
	} {
		state0 := NewStateFromFile(fname)
		if state0 == nil {
			continue
		}

		state := state0.DeepCopy()

		for true {
			if state.something_in_front() {
				state.turn90()
			} else {
				if !state.step() {
					break
				}
			}
		}

		// Part 1:
		fmt.Println("distinct positions: ")
		fmt.Println(state.distinct)

		// Part 2:
		c := 0
		originalBoard := state.board
		for oy, r0 := range state0.board {
			for ox, v := range r0 {
				if v != '.' || originalBoard[oy][ox] != 'X' {
					continue
				}
				state := state0.DeepCopy()
				state.board[oy][ox] = '#'
				for true {
					if state.something_in_front() {
						state.turn90()
					} else {
						if !state.step() {
							break
						}
					}
				}

				if state.dejaVu[state.pos] {
					c++
				}
			}
		}
		fmt.Println("Obstruction positions: ")
		fmt.Println(c)
	}
}
