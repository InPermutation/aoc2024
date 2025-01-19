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

// Check if two gates are equivalent
func gateEquiv(actual, expected gate, equiv map[string]string) bool {
	if actual.op != expected.op {
		return false
	}
	e := gate{equiv[actual.a], actual.op, equiv[actual.b]}

	if e.a != expected.a && e.b != expected.a {
		return false
	}
	if e.a != expected.b && e.b != expected.b {
		return false
	}
	return true
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

		// Part 1
		fmt.Println(res)

		// Define an equivalence between the input and my expected notation.
		equiv := map[string]string{}

		// existing notation:
		// 'xNN', 'yNN' are inputs
		for z := 0; z <= hiZ; z++ {
			xg := fmt.Sprintf("x%02d", z)
			equiv[xg] = xg
			yg := fmt.Sprintf("y%02d", z)
			equiv[yg] = yg
		}
		// 'zNN' is output. equiv["zNN"] is set when it's completely equivalent

		// my notation:
		// z[NN-1].a AND z[NN-1].b -> pNN.
		// v[NN-1].a AND v[NN-1].b -> qNN.
		// cNN is carry. pNN OR qNN -> cNN.
		// vNN is half-add. xNN XOR yNN -> vNN.
		// vNN XOR cNN -> zNN.

		// ...except for z00,
		z00 := s.gates["z00"]
		if gateEquiv(z00, gate{"x00", "XOR", "y00"}, equiv) {
			equiv["z00"] = "z00"
		} else {
			fmt.Println("warn: z00 not equiv. expected ", gate{"x00", "XOR", "y00"}, "; got ", z00)
		}

		// ...and for z01, which should be vNN XOR (x00 AND y00)
		z01 := s.gates["z01"]

		if z01.op == "XOR" {
			na, nb := z01.a, z01.b
			a, b := s.gates[na], s.gates[nb]
			if a.op == "AND" && b.op == "XOR" {
				na, nb = nb, na
				a, b = b, a
			}
			if gateEquiv(a, gate{"x01", "XOR", "y01"}, equiv) {
				// we could set equiv[na]="v01" here, but it'll be set later
				if gateEquiv(b, gate{"x00", "AND", "y00"}, equiv) {
					equiv["z01"] = "z01"
				}
			}
		}
		if _, ok := equiv["z01"]; !ok {
			fmt.Println("warn: z01 not equiv. expected", gate{"v01", "XOR", "c01"}, "; got ", s.gates["z01"])
		}

		// Try to map all the vNNs:
		for name, g := range s.gates {
			if g.op != "XOR" {
				continue
			}
			if g.a[0] != 'x' && g.a[0] != 'y' {
				continue
			}
			if g.b[0] != 'x' && g.b[0] != 'y' {
				continue
			}
			if g.a[1:] != g.b[1:] {
				continue
			}
			equiv[name] = "v" + g.a[1:]
		}
		vs := map[string]string{}
		for k, v := range equiv {
			if v[0] == 'v' {
				vs[v] = k
			}
		}
		// There should be a 'vNN' for all NN -- except hiZ (zHI is exactly cHI)
		for z := 1; z < hiZ; z++ {
			vv := fmt.Sprintf("v%02d", z)
			if _, ok := vs[vv]; !ok {
				fmt.Println(vv, "not found")
			}
		}
		// zHI must be cHI
		zHI := s.gates[fmt.Sprintf("z%02d", hiZ)]
		if zHI.op != "OR" {
			fmt.Println("zHI must be OR; ", zHI)
			// TODO: finish validating zHI
		}

		// There should be a 'qNN' for all NN
		//deleteme:
		// v[NN-1].a AND v[NN-1].b -> qNN.
		// cNN is carry. pNN OR qNN -> cNN.
		// vNN is half-add. xNN XOR yNN -> vNN.
		// /deleteme

		for name, g := range s.gates {
			if g.op != "AND" {
				continue
			}
			if g.a[0] != 'x' && g.b[0] != 'x' {
				continue
			}
			if g.a[0] != 'y' && g.b[0] != 'y' {
				continue
			}
			if g.a[1:] != g.b[1:] {
				fmt.Println(name, g, "would be qNN but has mismatched xNN, yMM")
				continue
			}
			fmt.Println(name, g, " is probably a qNN")
		}
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

// find the highest 'z' so we know how big our machine is
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
