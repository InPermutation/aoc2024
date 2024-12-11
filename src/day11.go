package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type state struct {
	fname  string
	stones *list.List
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
	s.stones = &list.List{}

	for scanner.Scan() {
		line := scanner.Text()
		nums := strings.Fields(line)
		for _, num := range nums {
			i, err := strconv.Atoi(num)
			if err != nil {
				log.Fatal(err)
			}
			s.stones.PushBack(i)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/11/sample",
		"input/11/input",
	} {
		s := readFile(fname)
		if s.stones == nil {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)

		for i := 0; i < 25; i++ {
			for e := s.stones.Front(); e != nil; e = e.Next() {
				if e.Value == 0 {
					e.Value = 1
				} else {
					str := strconv.Itoa(e.Value.(int))
					if len(str)%2 == 0 {
						il, err := strconv.Atoi(str[:len(str)/2])
						if err != nil {
							log.Fatal(err)
						}
						ir, err := strconv.Atoi(str[len(str)/2:])
						if err != nil {
							log.Fatal(err)
						}
						s.stones.InsertBefore(il, e)
						e.Value = ir
					} else {
						e.Value = e.Value.(int) * 2024
					}
				}
			}
		}

		// Part 1
		fmt.Println("len(stones): ", s.stones.Len())
	}
}
