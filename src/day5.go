package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type set map[int]bool

func main() {
	for _, fname := range []string{
		"input/05/sample",
		"input/05/input",
	} {
		fmt.Println(fname)

		file, err := os.Open(fname)
		if err != nil {
			log.Print(err)
			continue
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)

		orderingRules := map[int]set{}
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				break
			}
			tokens := strings.Split(line, "|")
			if len(tokens) != 2 {
				panic(tokens)
			}
			l, err := strconv.Atoi(tokens[0])
			if err != nil {
				log.Fatal(err)
			}
			r, err := strconv.Atoi(tokens[1])
			if err != nil {
				log.Fatal(err)
			}
			if orderingRules[l] == nil {
				orderingRules[l] = set{}
			}
			orderingRules[l][r] = true
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		middlePageNumbersOk := 0
		for scanner.Scan() {
			tokens := strings.Split(scanner.Text(), ",")

			nums := make([]int, len(tokens))
			for i, v := range tokens {
				nums[i], err = strconv.Atoi(v)
				if err != nil {
					log.Fatal(err)
				}
			}

			ok := true
			for i, l := range nums {
				for _, r := range nums[i+1:] {
					if orderingRules[r][l] {
						ok = false
						break
					}

				}
			}

			if ok {
				if len(nums)%2 == 0 {
					log.Fatal("Assume line must have odd number of updates")
				}
				middlePageNumber := nums[len(nums)/2]
				middlePageNumbersOk += middlePageNumber
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		// Part 1:
		fmt.Print("Sum of middle page numbers of correctly-ordered updates: ")
		fmt.Println(middlePageNumbersOk)
	}
}
