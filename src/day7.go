package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	for _, fname := range []string{
		"input/07/sample",
		"input/07/input",
	} {
		fmt.Println(fname)
		file, err := os.Open(fname)
		if err != nil {
			log.Print(err)
			continue
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)

		sum := 0
		sum3 := 0
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				log.Fatal("wrong number of parts", parts)
			}
			tgt, err := strconv.Atoi(parts[0])
			if err != nil {
				log.Fatal(err)
			}
			vals := []int{}
			for _, s := range strings.Split(parts[1], " ") {
				if s == "" {
					continue
				}
				i, err := strconv.Atoi(s)
				if err != nil {
					log.Fatal(err)
				}
				vals = append(vals, i)
			}
			if isPossible2(tgt, vals[0], vals[1:]) {
				sum += tgt
				sum3 += tgt
			} else if isPossible3(tgt, vals[0], vals[1:]) {
				sum3 += tgt
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		// Part 1
		fmt.Println("total calibration result: ", sum)

		// Part 2
		fmt.Println("new total calibration result: ", sum3)
	}
}

func isPossible2(tgt int, prev int, vals []int) bool {
	if len(vals) == 0 {
		return tgt == prev
	}
	curr := vals[0]
	rest := vals[1:]

	return (isPossible2(tgt, prev+curr, rest) ||
		isPossible2(tgt, prev*curr, rest))
}

func concat(a, b int) int {
	sa := strconv.Itoa(a)
	sb := strconv.Itoa(b)
	rv, err := strconv.Atoi(sa + sb)
	if err != nil {
		panic(err)
	}
	return rv
}

func isPossible3(tgt, prev int, vals []int) bool {
	if len(vals) == 0 {
		return tgt == prev
	}
	curr := vals[0]
	rest := vals[1:]

	return (isPossible3(tgt, prev+curr, rest) ||
		isPossible3(tgt, prev*curr, rest) ||
		isPossible3(tgt, concat(prev, curr), rest))
}
