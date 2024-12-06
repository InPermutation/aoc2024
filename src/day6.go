package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type page string
type set map[page]bool

func main() {
	for _, fname := range []string{
		"input/06/sample",
		"input/06/input",
	} {
		fmt.Println(fname)

		file, err := os.Open(fname)
		if err != nil {
			log.Print(err)
			continue
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)

		var x0, y0 int
		var dir0 rune

		board0 := [][]rune{}

		for scanner.Scan() {
			line := scanner.Text()
			bytes := make([]rune, len(line))

			for i, b := range line {
				bytes[i] = b
				if b == '^' {
					dir0 = b
					x0 = i
					y0 = len(board0)
				}
			}

			board0 = append(board0, bytes)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		x, y, dir := x0, y0, dir0
		board := deepcopy(board0)
		fmt.Println(x, y, dir)

		distinct := 0
		something_in_front := func() bool {
			switch dir {
			case '^':
				return y > 0 && board[y-1][x] == '#'
			case '<':
				return x > 0 && board[y][x-1] == '#'
			case 'v':
				return y+1 < len(board) && board[y+1][x] == '#'
			case '>':
				return x+1 < len(board[y]) && board[y][x+1] == '#'
			default:
				panic(dir)
			}
		}
		turn90 := func() {
			switch dir {
			case '^':
				dir = '>'
			case '>':
				dir = 'v'
			case 'v':
				dir = '<'
			case '<':
				dir = '^'
			default:
				panic(dir)
			}
		}
		step := func() bool {
			if board[y][x] != 'X' {
				distinct++
			}
			board[y][x] = 'X'
			switch dir {
			case '^':
				y--
				return y >= 0
			case '>':
				x++
				return x < len(board[y])
			case 'v':
				y++
				return y < len(board)
			case '<':
				x--
				return x >= 0
			default:
				panic(dir)
			}
		}

		for true {
			if something_in_front() {
				turn90()
			} else {
				if !step() {
					break
				}
			}
		}

		// Part 1:
		fmt.Println("distinct positions: ")
		fmt.Println(distinct)
	}
}

func deepcopy(board0 [][]rune) [][]rune {
	rval := make([][]rune, len(board0))
	for y, row0 := range board0 {
		row := make([]rune, len(row0))

		for x, o := range row0 {
			row[x] = o
		}

		rval[y] = row
	}

	return rval
}
