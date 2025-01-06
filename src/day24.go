package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type gate struct {
	a  string
	op string
	b  string
}

type state struct {
	fname string
	wires map[string]bool
	gates map[string]gate
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

	s.wires = map[string]bool{}
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}

		toks := strings.Split(line, ":")
		if len(toks) != 2 {
			log.Fatal(line)
		}
		if len(toks[1]) != 2 || toks[1][0] != ' ' {
			log.Fatal(line)
		}

		s.wires[toks[0]] = (toks[1] == " 1")
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	s.gates = map[string]gate{}
	for scanner.Scan() {
		line := scanner.Text()

		toks := strings.Split(line, " ")
		if len(toks) != 5 || toks[3] != "->" {
			log.Fatal(line)
		}

		s.gates[toks[4]] = gate{
			a:  toks[0],
			op: toks[1],
			b:  toks[2],
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/24/sample",
		"input/24/larger",
		"input/24/input",
	} {
		s := readFile(fname)
		if len(s.wires) == 0 || len(s.gates) == 0 {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)

		hiZ, res, bit := s.hiZ(), 0, 1

		for z := 0; z <= hiZ; z++ {
			sz := fmt.Sprintf("z%02d", z)
			if s.eval(sz) {
				res |= bit
			}
			bit <<= 1
		}
		fmt.Println(res)
	}
}

func (s *state) eval(w string) bool {
	if g, ok := s.gates[w]; ok {
		a := s.eval(g.a)
		b := s.eval(g.b)
		switch g.op {
		default:
			panic(g)
		case "AND":
			return a && b
		case "OR":
			return a || b
		case "XOR":
			return a != b
		}
	} else if v, ok := s.wires[w]; ok {
		return v
	} else {
		fmt.Println(s.gates)
		fmt.Println(s.wires)
		panic(w)
	}
}

func (s *state) hiZ() (hi int) {
	for g := range s.gates {
		if g[0] == 'z' {
			z, err := strconv.Atoi(string(g[1:]))
			if err != nil {
				log.Fatal("bad z", g)
			}
			if z > hi {
				hi = z
			}
		}
	}

	return
}
