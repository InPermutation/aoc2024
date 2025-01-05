package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type state struct {
	fname   string
	A       int
	B       int
	C       int
	Ip      int
	Program []int
	Output  []int
}

func (s *state) DeepCopy() *state {
	return &state{
		s.fname,
		s.A,
		s.B,
		s.C,
		s.Ip,
		slices.Clone(s.Program),
		slices.Clone(s.Output),
	}
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

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Register ") {
			reg, err := strconv.Atoi(line[len("Register _: "):])
			if err != nil {
				log.Fatal(err)
			}
			switch line[len("Register ")] {
			case 'A':
				s.A = reg
			case 'B':
				s.B = reg
			case 'C':
				s.C = reg
			}
		} else if strings.HasPrefix(line, "Program: ") {
			toks := strings.Split(line[len("Program: "):], ",")
			s.Program = make([]int, len(toks))
			for i, str := range toks {
				v, err := strconv.Atoi(str)
				if err != nil {
					log.Fatal(err)
				}
				s.Program[i] = v
			}

		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/17/sample",
		"input/17/input",
	} {
		s := readFile(fname)
		if len(s.Program) == 0 {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func dv(n, b int) int {
	return n / (1 << b)
}

func (s *state) step() bool {
	op := s.Program[s.Ip]
	arg := s.Program[s.Ip+1]
	combo := arg
	if arg == 4 {
		combo = s.A
	} else if arg == 5 {
		combo = s.B
	} else if arg == 6 {
		combo = s.C
	}

	switch op {
	case 0: // adv
		s.A = dv(s.A, combo)
		s.Ip += 2
	case 1: // bxl
		s.B = s.B ^ arg
		s.Ip += 2
	case 2: // bst
		s.B = combo % 8
		s.Ip += 2
	case 3: // jnz
		if s.A == 0 {
			s.Ip += 2
		} else {
			s.Ip = arg
		}
	case 4: // bxc
		s.B = (s.B ^ s.C)
		s.Ip += 2
	case 5: // out
		s.Output = append(s.Output, combo%8)
		s.Ip += 2
	case 6: // bdv
		s.B = dv(s.A, combo)
		s.Ip += 2
	case 7: // cdv
		s.C = dv(s.A, combo)
		s.Ip += 2
	default:
		log.Fatal("unknown opcode ", s.Program[s.Ip])
	}

	return s.Ip < len(s.Program)
}
func (s *state) Part1() {
	for s.step() {
	}
	for i, v := range s.Output {
		if i != 0 {
			fmt.Print(",")
		}
		fmt.Print(v)
	}
	fmt.Println()
}

func (s *state) Part2(lm int) int {
	for a := 0; a <= 7; a++ {
		at := lm<<3 | a
		s1 := s.DeepCopy()
		s1.A = at
		for s1.step() {
		}

		// only need to check first output
		if s.Program[len(s.Program)-len(s1.Output)] != s1.Output[0] {
			continue
		}
		if len(s.Program) == len(s1.Output) {
			// win!
			return at
		}
		possible := s.Part2(at)
		if possible >= 0 {
			return possible
		}
	}
	return -1
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)
		s.DeepCopy().Part1()
		if s.fname == "input/17/input" {
			fmt.Println("quine:", s.Part2(0))
		}
	}
}

// :r input/17/input
// Register A: 65804993
// Register B: 0
// Register C: 0
// Program: 2,4,1,1,7,5,1,4,0,3,4,5,5,5,3,0
// --
// Decompiled:
// location
// |  opcode
// |  | operand
// |  | | disassembly
// v  v v v         ; comment
// 0  2 4 bst a     ; b = a
// 2  1 1 bxl 1     ; b = b ^ 1 -- toggle low bit of b
// 4  7 5 cdv (2^b) ; c = a / (2^b)
// 6  1 4 bxl 4     ; b = b ^ 4 -- toggle the 3rd bit of b
// 8  0 3 adv (2^3) ; a = a / 8 -- a>>=3
// 10 4 5 bxc       ; b = b ^ c
// 12 5 5 out b     ; print b
// 14 3 0 jnz 0     ; loop if a is nonzero

// analysis:
// - b and c are temporaries; only the value of a matters
// - at 8, a will shrink by 3 bits
// - at 14, loop back
// - NO INFINITE LOOPS
// - ALWAYS MAKES PROGRESS
// - therefore we can iterate from the end of the program and
//   see what the first, second, Nth triad is.
