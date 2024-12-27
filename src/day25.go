package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type state struct {
	fname string
	keys  [][5]int
	locks [][5]int
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

	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			if len(lines) != 7 {
				log.Fatal(len(lines), lines)
			}
			if lines[0] == "#####" {
				lock := [5]int{}
				for i, v := range lines {
					for l := range lock {
						if v[l] == '#' {
							lock[l] = i
						}
					}
				}
				s.locks = append(s.locks, lock)
			} else {
				key := [5]int{}
				for i := range lines {
					v := lines[len(lines)-i-1]
					for l := range key {
						if v[l] == '#' {
							key[l] = i
						}
					}
				}
				s.keys = append(s.keys, key)
			}
			lines = lines[:0]
		} else {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/25/sample",
		"input/25/input",
	} {
		s := readFile(fname)
		if len(s.locks) == 0 || len(s.keys) == 0 {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)
		c := 0
		for _, key := range s.keys {
			for _, lock := range s.locks {
				ok := true
				for i := 0; i < 5; i++ {
					if key[i]+lock[i] > 5 {
						ok = false
					}
				}
				if ok {
					c++
				}
			}
		}
		fmt.Println("fits:", c)
	}
}
