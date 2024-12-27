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

		for len(s.gates) > 0 {
			found := false
			for o, g := range s.gates {
				if a, ok := s.wires[g.a]; ok {
					if b, ok := s.wires[g.b]; ok {
						var res bool
						switch g.op {
						default:
							log.Fatal("unknown op", g)
						case "AND":
							res = a && b
						case "OR":
							res = a || b
						case "XOR":
							res = a != b
						}
						s.wires[o] = res
						found = true
						delete(s.gates, o)
						break
					}
				}
			}
			if !found {
				log.Fatal("no progress", s)
			}
		}

		res := 0
		for k, v := range s.wires {
			if v && k[0] == 'z' {
				i, err := strconv.Atoi(k[1:])
				if err != nil {
					panic(err)
				}
				res |= 1 << i
			}
		}
		fmt.Println(res)
	}
}
