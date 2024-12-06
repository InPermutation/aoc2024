package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type direction rune
type state struct {
	x     int
	y     int
	dir   direction
	board [][]rune

	distinct int
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
		board: [][]rune{},
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
func (p *state) something_in_front() bool {
	switch p.dir {
	case '^':
		return p.y > 0 && p.board[p.y-1][p.x] == '#'
	case '<':
		return p.x > 0 && p.board[p.y][p.x-1] == '#'
	case 'v':
		return p.y+1 < len(p.board) && p.board[p.y+1][p.x] == '#'
	case '>':
		return p.x+1 < len(p.board[p.y]) && p.board[p.y][p.x+1] == '#'
	default:
		fmt.Println("err", p.dir)
		panic(*p)
	}
}
func (p *state) turn90() {
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
func (p *state) step() bool {
	if p.board[p.y][p.x] != 'X' {
		p.distinct++
	}
	p.board[p.y][p.x] = 'X'
	switch p.dir {
	case '^':
		p.y--
		return p.y >= 0
	case '>':
		p.x++
		return p.x < len(p.board[p.y])
	case 'v':
		p.y++
		return p.y < len(p.board)
	case '<':
		p.x--
		return p.x >= 0
	default:
		panic(*p)
	}
}
func (p *state) DeepCopy() (rval *state) {
	rval = &state{
		x:     p.x,
		y:     p.y,
		dir:   p.dir,
		board: make([][]rune, len(p.board)),
	}

	for y, row0 := range p.board {
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
	}
}
