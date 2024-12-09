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
	chk := 0
	id := 0
	fsm := File
	for i := 0; i < s.width && len(s.files) > 0; i++ {
		if fsm == File {
			s.files[0]--
			chk += (i * id)
			if s.files[0] == 0 {
				id++
				fsm = Empty
				s.files = s.files[1:]
				if s.empty[0] == 0 {
					fsm = File
					s.empty = s.empty[1:]
				}
			}
		} else {
			last := len(s.files) - 1
			nid := id + last
			s.files[last]--
			if s.files[last] == 0 {
				s.files = s.files[:last]
			}
			chk += (i * nid)
			s.empty[0]--
			if s.empty[0] == 0 {
				fsm = File
				s.empty = s.empty[1:]
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
