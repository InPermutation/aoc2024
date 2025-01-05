package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type state struct {
	fname    string
	towels   []string
	patterns []string
	memo     map[string][]int
}

func readFile(fname string) state {
	s := state{
		fname: fname,
		memo:  map[string][]int{"": []int{1}},
	}
	file, err := os.Open(fname)
	if err != nil {
		log.Print(err)
		return s
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		switch i {
		case 0:
			toks := strings.Split(line, ",")
			for _, tok := range toks {
				s.towels = append(s.towels, strings.Trim(tok, " "))
			}
		case 1:
			if line != "" {
				log.Fatal(line)
			}
		default:
			s.patterns = append(s.patterns, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/19/sample",
		"input/19/input",
	} {
		s := readFile(fname)
		if len(s.patterns) == 0 || len(s.towels) == 0 {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func (s *state) Possible(design string) (rv []int) {
	if v, ok := s.memo[design]; ok {
		return v
	}

	for _, prefix := range s.towels {
		if strings.HasPrefix(design, prefix) {
			suffix := design[len(prefix):]
			sum := 0
			for _, p := range s.Possible(suffix) {
				sum += p
			}
			if sum != 0 {
				rv = append(rv, sum)
			}
		}
	}

	s.memo[design] = rv

	return
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)

		c, ways := 0, 0
		for _, design := range s.patterns {
			p := s.Possible(design)
			sum := 0
			for _, v := range p {
				sum += v
			}
			ways += sum
			if sum > 0 {
				c++
			}
		}

		fmt.Println(c)
		fmt.Println(ways)
	}
}
