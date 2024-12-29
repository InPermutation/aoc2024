package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type state struct {
	fname   string
	secrets []int
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
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		s.secrets = append(s.secrets, num)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/22/sample",
		"input/22/sample2",
		"input/22/input",
	} {
		s := readFile(fname)
		if len(s.secrets) == 0 {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func prune(i int) int {
	return i % 16777216
}

func mix(i, j int) int {
	return i ^ j
}

func next(i int) int {
	i = prune(mix(i, i*64))
	i = prune(mix(i, i/32))
	i = prune(mix(i, i*2048))
	return i
}

func main() {
	if prune(100000000) != 16113920 {
		log.Fatal("prune test failed")
	}
	if mix(15, 42) != 37 {
		log.Fatal("mix test failed")
	}
	test := []int{
		123,
		15887950,
		16495136,
		527345,
		704524,
		1553684,
		12683156,
		11100544,
		12249484,
		7753432,
		5908254,
	}

	for i, v := range test {
		if i == 0 {
			continue
		}
		if next(test[i-1]) != v {
			log.Fatal(test[i-1], v, next(test[i-1]))
		}
	}
	fmt.Println("tests passed")

	for _, s := range states() {
		fmt.Println(s.fname)

		sum := 0
		total := map[[4]int]int{}
		for _, v := range s.secrets {
			//fmt.Print(v, ": ")

			curr := [4]int{-99, -99, -99, -99}
			d := map[[4]int]int{}
			for i := 0; i < 2000; i++ {
				curr[0] = curr[1]
				curr[1] = curr[2]
				curr[2] = curr[3]
				n := next(v)
				curr[3] = n%10 - v%10
				v = n
				if _, ok := d[curr]; !ok {
					price := v % 10
					d[curr] = price
				}
			}
			sum += v
			//fmt.Println(v)

			for k, v := range d {
				total[k] += v
			}
		}
		fmt.Println("sum of 2000th:", sum)

		most := 0
		for try := [4]int{-9, -9, -9, -9}; try != [4]int{10, -9, -9, -9}; try = inc(try) {
			if a := total[try]; a > most {
				most = a
			}
		}
		fmt.Println("most bananas:", most)
	}
}

func inc(try [4]int) [4]int {
	if try[3] == 9 {
		try[3] = -9
		if try[2] == 9 {
			try[2] = -9
			if try[1] == 9 {
				try[1] = -9
				try[0]++
			} else {
				try[1]++
			}
		} else {
			try[2]++
		}
	} else {
		try[3]++
	}
	return try
}
