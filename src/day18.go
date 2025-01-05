package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type pos struct {
	x int
	y int
}

type state struct {
	fname     string
	exit      pos
	corrupted []pos
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

	if fname == "input/18/sample" {
		s.exit = pos{6, 6}
	} else {
		s.exit = pos{70, 70}
	}

	for scanner.Scan() {
		line := scanner.Text()
		toks := strings.Split(line, ",")
		if len(toks) != 2 {
			log.Fatal("invalid line", line)
		}
		x, err := strconv.Atoi(toks[0])
		if err != nil {
			log.Fatal(err)
		}
		y, err := strconv.Atoi(toks[1])
		if err != nil {
			log.Fatal(err)
		}
		s.corrupted = append(s.corrupted, pos{x, y})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/18/sample",
		"input/18/input",
	} {
		s := readFile(fname)
		if len(s.corrupted) == 0 {
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
