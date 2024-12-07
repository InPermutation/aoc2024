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
			if isPossible(tgt, vals, possibilities) {
				sum += tgt
				sum3 += tgt
			} else if isPossible(tgt, vals, possibilities3) {
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

func isPossible(tgt int, vals []int, possibilities func([]int) []int) bool {
	for _, v := range possibilities(vals) {
		if v == tgt {
			return true
		}
	}
	return false
}

func possibilities(vals []int) (rval []int) {
	if len(vals) == 0 {
		return
	}
	if len(vals) == 1 {
		return vals
	}

	rest := possibilities(vals[:len(vals)-1])
	for _, r := range rest {
		rval = append(rval, r+vals[len(vals)-1])
		rval = append(rval, r*vals[len(vals)-1])
	}
	return
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

func possibilities3(vals []int) []int {
	return possibilities3_internal(0, vals)
}

func possibilities3_internal(prev int, vals []int) (rval []int) {
	if len(vals) == 0 {
		return []int{prev}
	}

	curr := vals[0]
	rest := vals[1:]
	for _, r := range possibilities3_internal(prev+curr, rest) {
		rval = append(rval, r)
	}
	for _, r := range possibilities3_internal(prev*curr, rest) {
		rval = append(rval, r)
	}
	for _, r := range possibilities3_internal(concat(prev, curr), rest) {
		rval = append(rval, r)
	}
	return
}
