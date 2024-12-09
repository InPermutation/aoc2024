package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type state struct {
	fname string
	src   string
	files []int
	empty []int
	width int
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
		if s.src != "" {
			log.Fatal("Too many lines")
		}
		s.src = line
		for i, r := range line {
			v := int(r - '0')
			if i%2 == 0 {
				s.width += v
				s.files = append(s.files, v)
			} else {
				s.empty = append(s.empty, v)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/09/small",
		"input/09/sample",
		"input/09/input",
	} {
		s := readFile(fname)
		if s.src == "" {
			continue
		}
		rval = append(rval, s)
	}
	return
}

const (
	File = iota
	Empty
)

func (s state) checksum() int {
	var lfile, rfile, lempty int // pointers into files / empty
	var subl, subr int           // how much we have consumed of (lfile, rfile)
	var sube int                 // how much we have consumed of lempty
	chk, id := 0, 0
	fsm := File
	rfile = len(s.files) - 1
	for i := 0; i < s.width && len(s.files) > 0; i++ {
		if fsm == File {
			subl++
			chk += (i * id)
			if subl == s.files[lfile] {
				id++
				fsm = Empty
				lfile++
				subl = 0
				if s.empty[lempty] == 0 {
					fsm = File
					lempty++
				}
			}
		} else {
			nid := id + (rfile - lfile)
			subr++
			if s.files[rfile] == subr {
				rfile--
				subr = 0
			}
			chk += (i * nid)
			sube++
			if s.empty[lempty] == sube {
				fsm = File
				lempty++
				sube = 0
			}
		}
	}

	return chk
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)
		// Part 1
		fmt.Println("checksum: ", s.checksum())
	}
}
