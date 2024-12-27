package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type state struct {
	fname string
	Conn  [26 * 26][26 * 26]bool
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
		toks := strings.Split(line, "-")
		if len(toks) != 2 {
			log.Fatal(line)
		}
		if len(toks[0]) != 2 || len(toks[1]) != 2 {
			log.Fatal(line)
		}

		c1 := sToConn(toks[0])
		c2 := sToConn(toks[1])

		s.Conn[c1][c2] = true
		s.Conn[c2][c1] = true
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/23/sample",
		"input/23/input",
	} {
		s := readFile(fname)
		if len(s.Conn) == 0 {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func sToConn(s string) int {
	if len(s) != 2 {
		panic(s)
	}
	return int(s[0]-'a')*26 + int(s[1]-'a')
}

func main() {
	t := int('t' - 'a')
	for _, s := range states() {
		fmt.Println(s.fname)
		c := 0
		for i := 0; i < 26*26; i++ {
			for j := i + 1; j < 26*26; j++ {
				if !s.Conn[i][j] {
					continue
				}

				for k := j + 1; k < 26*26; k++ {
					if s.Conn[i][k] && s.Conn[j][k] {
						if i/26 == t || j/26 == t || k/26 == t {
							c++
						}
					}
				}
			}
		}
		fmt.Println(c)
	}
}
