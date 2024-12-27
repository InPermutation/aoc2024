package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
		fmt.Println(s)
	}
}
