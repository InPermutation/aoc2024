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
			if isPossible(tgt, vals) {
				sum += tgt
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		// Part 1
		fmt.Println("total calibration result: ", sum)
	}
}

func isPossible(tgt int, vals []int) bool {
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
