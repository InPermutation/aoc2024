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

	swaps := map[string]string{
		// Found by observation:
		"fcd": "z33",
		"z33": "fcd",
		"hmk": "z16",
		"z16": "hmk",
		"fhp": "z20",
		"z20": "fhp",
		// TODO: 2 more...
	}

	s.gates = map[string]gate{}
	for scanner.Scan() {
		line := scanner.Text()

		toks := strings.Split(line, " ")
		if len(toks) != 5 || toks[3] != "->" {
			log.Fatal(line)
		}

		t4 := toks[4]
		if sw, ok := swaps[t4]; ok {
			t4 = sw
		}
		s.gates[t4] = gate{
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
	e := gate{equiv[actual.a], actual.op, equiv[actual.b]}
	return gateEqual(e, expected)
}

func gateEqual(actual, expected gate) bool {
	if actual.op != expected.op {
		return false
	}

	if actual.a != expected.a && actual.b != expected.a {
		return false
	}
	if actual.a != expected.b && actual.b != expected.b {
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
		// z[NN-1].a AND z[NN-1].b -> pNN.       # VALIDATED, ps
		// v[NN-1].a AND v[NN-1].b -> qNN.       # VALIDATED, qs
		// cNN is carry. pNN OR qNN -> cNN.      # PARTIAL VALIDATED, cs
		// vNN is half-add. xNN XOR yNN -> vNN.  # VALIDATED, vs
		// vNN XOR cNN -> zNN.                   # VALIDATED, zs

		// ...except for z00,
		z00 := s.gates["z00"]
		if gateEquiv(z00, gate{"x00", "XOR", "y00"}, equiv) {
			equiv["z00"] = "z00"
		} else {
			fmt.Println("warn: z00 not equiv. expected ", gate{"x00", "XOR", "y00"}, "; got ", z00)
		}

		// ...and for z01, which should be v01 XOR (x00 AND y00)
		z01 := s.gates["z01"]
		v01 := gate{"x01", "XOR", "y01"}
		// call (x00 AND y00) "c01" for the recursion to work
		c01 := gate{"x00", "AND", "y00"}

		if z01.op == "XOR" {
			na, nb := z01.a, z01.b
			a, b := s.gates[na], s.gates[nb]
			if a.op == "AND" && b.op == "XOR" {
				na, nb = nb, na
				a, b = b, a
			}
			if gateEquiv(a, v01, equiv) {
				// we could set equiv[na]="v01" here, but it'll be set later
				if gateEquiv(b, c01, equiv) {
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
				fmt.Println(vv, " vv not found")
			}
		}

		// Try to map all the qNNs:
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
			n, err := strconv.Atoi(string(g.a[1:]))
			if err != nil {
				fmt.Println("non-numeric; but would be qNN: ", name, g)
				continue
			}
			qname := fmt.Sprintf("q%02d", n+1)
			equiv[name] = qname
		}
		qs := map[string]string{}
		for k, v := range equiv {
			if v[0] == 'q' {
				qs[v] = k
			}
		}
		// There should be a 'qNN' for all NN
		for z := 1; z <= hiZ; z++ {
			qq := fmt.Sprintf("q%02d", z)
			if _, ok := qs[qq]; !ok {
				fmt.Println(qq, " qq not found")
			}
		}

		// Try to map all the cNNs:
		for name, g := range s.gates {
			if g.op != "OR" {
				continue
			}
			na, nb := g.a, g.b
			a, b := s.gates[na], s.gates[nb]
			if _, ok := equiv[na]; ok {
				na, nb = nb, na
				a, b = b, a
			}
			if qeq, ok := equiv[nb]; ok {
				nn, err := strconv.Atoi(qeq[1:])
				cnn := fmt.Sprintf("c%02d", nn)
				if err == nil && qeq[0] == 'q' {
					equiv[name] = cnn
				} else {
					fmt.Println(name, g, "would be", cnn, "but wrong qeq=", qeq, "/", nb, "; nn, err:=", nn, err)
				}
			} else {
				fmt.Println("no qNN for cNN", name, g, ", equivs: ", equiv[na], ",", equiv[nb])
			}
		}

		cs := map[string]string{}
		for k, v := range equiv {
			if v[0] == 'c' {
				cs[v] = k
			}
			if v == "c27" {
				fmt.Println("c27:", k, v, s.gates[k])
			}
		}
		if q01, ok := qs["q01"]; ok {
			cs["c01"] = q01
		} else {
			fmt.Println("no q01 to turn into c01")
		}
		// There should be a 'cNN' for all NN
		for z := 1; z <= hiZ; z++ {
			cc := fmt.Sprintf("c%02d", z)
			if _, ok := cs[cc]; !ok {
				fmt.Println(cc, " cc not found")
			}
		}

		for name, g := range s.gates {
			if name == "z01" {
				continue
			}
			if g.op != "XOR" {
				continue
			}
			if maybeV, ok := equiv[name]; ok {
				if maybeV[0] != 'v' {
					fmt.Println("weird equiv ", maybeV, name, "for", name, g)
				}
				// already a vNN
				continue
			}
			ea, eb := equiv[g.a], equiv[g.b]
			if ea == "" || eb == "" {
				fmt.Println("unresolved possible z", name, g, equiv[g.a], ", ", equiv[g.b])
				continue
			} else if ea[1:] == eb[1:] && name[0] == 'z' {
				if ea[0] == 'c' || ea[0] == 'v' && eb[0] == 'c' || eb[0] == 'v' {
					equiv[name] = name
				} else {
					fmt.Println("unknown likely z", name, g, "equivs: ", ea, ",", eb)
				}
			} else {
				fmt.Println("unknown possible z", name, g, "equivs: ", ea, ",", eb)
			}
		}

		zs := map[string]string{}
		for k, v := range equiv {
			if v[0] == 'z' {
				zs[v] = k
			}
		}
		// There should be a 'zNN' for all NN
		for z := 0; z <= hiZ; z++ {
			zz := fmt.Sprintf("z%02d", z)
			if _, ok := zs[zz]; !ok {
				fmt.Println(zz, " zz not found")
			}
		}

		for name, g := range s.gates {
			if g.op != "AND" {
				continue
			}
			if _, ok := equiv[name]; ok {
				continue
			}
			ea, eb := equiv[g.a], equiv[g.b]
			if ea == "" || eb == "" {
				continue
			}
			for zNN := range zs {
				eg := s.gates[zNN]

				ega := eg.a == g.a || eg.a == g.b
				egb := eg.b == g.a || eg.b == g.b

				if ega && egb {
					nn, err := strconv.Atoi(zNN[1:])
					if err != nil {
						fmt.Println(zNN, err)
						continue
					}
					equiv[name] = fmt.Sprintf("p%02d", nn+1)
					break
				} else if ega || egb {
					fmt.Println([]string{g.a, g.b, eg.a, eg.b}, ega, egb)
				}
			}

			if _, ok := equiv[name]; !ok {
				fmt.Println(name, g, "probably a pNN, but no equiv found", ea, eb)
			}
		}
		ps := map[string]string{}
		for k, v := range equiv {
			if v[0] == 'p' {
				ps[v] = k
			}
		}
		// There should be a 'pNN' for all NN

		for z := 2; z <= hiZ; z++ {
			pp := fmt.Sprintf("p%02d", z)
			if _, ok := ps[pp]; !ok {
				fmt.Println(pp, " pp not found")
			}
		}

		// zHI is equivalent to a carry (zHI===cHI)
		zHI := s.gates[fmt.Sprintf("z%02d", hiZ)]
		cHI := cs[fmt.Sprintf("c%02d", hiZ)]
		if !gateEqual(zHI, s.gates[cHI]) {
			fmt.Println("zHI", zHI, "must be cHI ", cHI, s.gates[cHI])
		}

		fmt.Println(equiv)
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
