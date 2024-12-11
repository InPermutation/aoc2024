package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var dp map[int]map[int]int = map[int]map[int]int{}

type state struct {
	fname  string
	stones []int
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
		nums := strings.Fields(line)
		for _, num := range nums {
			i, err := strconv.Atoi(num)
			if err != nil {
				log.Fatal(err)
			}
			s.stones = append(s.stones, i)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/11/sample",
		"input/11/input",
	} {
		s := readFile(fname)
		if s.stones == nil {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func onelen(v, depth int) int {
	if depth == 0 {
		return 1
	}
	if v == 0 {
		return onelen(1, depth-1)
	}

	if atd, found := dp[v]; found {
		if rv, found := atd[depth]; found {
			return rv
		}
	} else {
		dp[v] = map[int]int{}
	}

	str := strconv.Itoa(v)
	sum := 0
	if len(str)%2 == 0 {
		il, err := strconv.Atoi(str[:len(str)/2])
		if err != nil {
			log.Fatal(err)
		}
		ir, err := strconv.Atoi(str[len(str)/2:])
		if err != nil {
			log.Fatal(err)
		}
		sum = onelen(il, depth-1) + onelen(ir, depth-1)
	} else {
		sum = onelen(v*2024, depth-1)
	}

	dp[v][depth] = sum
	return sum
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)
		for _, depth := range []int{25, 75} {
			sum := 0
			for _, v := range s.stones {
				sum += onelen(v, depth)
			}
			fmt.Println("len(stones) @ ", depth, " : ", sum)
		}
	}
}
