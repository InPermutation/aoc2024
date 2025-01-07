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

func expected(hiZ int, s0 *state) *state {
	s := state{
		wires: s0.wires,
		gates: map[string]gate{},
	}
	s.gates["z00"] = gate{"x00", "XOR", "y00"}
	s.gates["z01"] = gate{"v01", "XOR", "c01"}
	s.gates["v01"] = gate{"x01", "XOR", "y01"}
	s.gates["c01"] = gate{"x00", "AND", "y00"}

	for i := 2; i <= hiZ; i++ {
		suffix := fmt.Sprintf("%02d", i)
		s.gates["z"+suffix] = gate{"v" + suffix, "XOR", "c" + suffix}
		s.gates["v"+suffix] = gate{"x" + suffix, "XOR", "y" + suffix}
		s.gates["c"+suffix] = gate{"p" + suffix, "OR", "q" + suffix}
		z0 := s.gates[fmt.Sprintf("z%02d", i-1)]
		s.gates["p"+suffix] = gate{z0.a, "AND", z0.b}
		v0 := s.gates[fmt.Sprintf("v%02d", i-1)]
		s.gates["q"+suffix] = gate{v0.a, "AND", v0.b}
	}

	return &s
}

var empty gate

func descend(ex, s *state, ez, sz string) error {
	exg := ex.gates[ez]
	sg := s.gates[sz]
	if exg == empty && sg == empty {
		if ez != sz {
			return fmt.Errorf("wire mismatch: %q != %q", ez, sz)
		}
		return nil
	}
	if exg == empty && sg != empty {
		return fmt.Errorf("empty exg, nonempty %q", sg)
	}
	if exg != empty && sg == empty {
		return fmt.Errorf("empty " + sz + " nonempty exg")
	}
	if exg.op != sg.op {
		return fmt.Errorf("op mismatch: %q != %q. (%v, %v)", exg.op, sg.op, ez, sz)
	}

	eaa := descend(ex, s, exg.a, sg.a)
	ebb := descend(ex, s, exg.b, sg.b)
	if eaa == ebb && eaa == nil {
		return nil
	}
	eab := descend(ex, s, exg.a, sg.b)
	eba := descend(ex, s, exg.b, sg.a)
	if eab == nil && eba == nil {
		return nil
	}
	if eaa == nil {
		return ebb
	}
	if ebb == nil {
		return eaa
	}
	if eab == nil {
		return eba
	}
	if eba == nil {
		return eab
	}
	return fmt.Errorf("mismatch at %q %v", ez, sz)
}

func main() {
	for _, s := range states() {
		//fmt.Println(s.fname)
		hiZ, res, bit := s.hiZ(), 0, 1
		ex := expected(hiZ, &s)

		for z := 0; z <= hiZ; z++ {
			sz := fmt.Sprintf("z%02d", z)
			if s.eval(sz) {
				res |= bit
			}
			bit <<= 1
			// z16 = AND
			// z20 = AND
			// z33 = OR
			// z45 = OR
			// dkb =
			exd := ex.Debug(sz)
			dbg := s.Debug(sz)

			fmt.Println(dbg)
			if exd == "Quuahuhs" && exd != dbg {
				err := descend(ex, &s, sz, sz)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		//fmt.Println(res)
	}
}

func (s *state) Debug(w string) string {
	if g, ok := s.gates[w]; ok {
		if w[0] == 'Z' {
			panic(w)
			if g.op != "XOR" {
				return "must be XOR: " + w
			}
			if w == "z00" {
				return ""
			}
			a, b := g.a, g.b
			ga, gb := s.gates[a], s.gates[b]
			if ga.op != "XOR" && gb.op == "XOR" || gb.op != "OR" && ga.op == "OR" {
				a, b = b, a
				ga, gb = gb, ga
			}
			rval := ""
			if ga.op != "XOR" {
				rval += "must have XOR child: " + w + " / " + a + " " + ga.op + "\n"
			}
			if gb.op != "OR" {
				rval += "must have OR child: " + w + "/ " + b + " " + gb.op + "\n"
			}

			if len(rval) != 0 {
				return rval
			}

		}

		sa := s.Debug(g.a)
		sb := s.Debug(g.b)
		if sa < sb {
			sa, sb = sb, sa
		}
		return fmt.Sprintf("{%q: [%v, %v]}", g.op, sa, sb)
	}
	return fmt.Sprintf("%q", w)
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
