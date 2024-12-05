package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type page string
type set map[page]bool

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

		orderingRules := map[page]set{}
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				break
			}
			tokens := strings.Split(line, "|")
			if len(tokens) != 2 {
				panic(tokens)
			}
			l, r := page(tokens[0]), page(tokens[1])
			if orderingRules[l] == nil {
				orderingRules[l] = set{}
			}
			orderingRules[l][r] = true
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		middlePageNumbersOk := 0
		middlePageNumbersReordered := 0
		for scanner.Scan() {
			tokens := strings.Split(scanner.Text(), ",")

			ok := true
			for i, l := range tokens {
				for _, r := range tokens[i+1:] {
					if orderingRules[page(r)][page(l)] {
						ok = false
						break
					}

				}
			}

			if ok {
				if len(tokens)%2 == 0 {
					log.Fatal("Assume line must have odd number of updates")
				}
				middlePageNumber, err := strconv.Atoi(tokens[len(tokens)/2])
				if err != nil {
					log.Fatal(err)
				}
				middlePageNumbersOk += middlePageNumber
			} else {
				ordering := []page{}
				pages := set{}
				inserted := set{}

				for _, v := range tokens {
					pages[page(v)] = true
				}

				for len(pages) > 0 {
					// find the lowest allowable page number
					var lowest = page("INVALID")
					for k := range pages {
						isLowest := true
						for k2, v := range orderingRules {
							if k2 == k {
								// ignore yourself
								continue
							}
							if !pages[page(k2)] {
								// this rule doesn't apply
								continue
							}
							if inserted[page(k2)] {
								// this page was already inserted
								continue
							}
							if v[page(k)] {
								// breaks rule
								isLowest = false
								break
							}
						}
						if isLowest {
							lowest = k
							break
						}
					}

					if lowest == page("INVALID") {
						log.Fatal("couldn't find lowest")
					}

					// insert it into ordering
					ordering = append(ordering, lowest)
					inserted[lowest] = true

					// remove it from pages
					delete(pages, lowest)
				}

				middlePageNumber, err := strconv.Atoi(string(ordering[len(ordering)/2]))
				if err != nil {
					log.Fatal(err)
				}
				middlePageNumbersReordered += middlePageNumber

			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		// Part 1:
		fmt.Print("Sum of middle page numbers of correctly-ordered updates: ")
		fmt.Println(middlePageNumbersOk)

		// Part 2:
		fmt.Print("Sum of middle page numbers of reordered updates: ")
		fmt.Println(middlePageNumbersReordered)
	}
}
