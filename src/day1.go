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

func main() {
	for _, fname := range []string{
		"input/01/sample",
		"input/01/input",
	} {
		fmt.Println(fname)

		file, err := os.Open(fname)
		if err != nil {
			log.Print(err)
			continue
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)

		var leftList, rightList []int

		for scanner.Scan() {
			line := scanner.Text()
			split := strings.Split(line, " ")
			first, last := split[0], split[len(split)-1]
			l, err := strconv.Atoi(first)
			if err != nil {
				log.Fatal(err)
			}
			leftList = append(leftList, l)
			r, err := strconv.Atoi(last)
			if err != nil {
				log.Fatal(err)
			}
			rightList = append(rightList, r)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		slices.Sort(leftList)
		slices.Sort(rightList)

		// Part 1:
		sum_diff := 0

		for i, l := range leftList {
			r := rightList[i]
			diff := l - r
			if diff < 0 {
				diff = -diff
			}

			sum_diff += diff
		}

		fmt.Print("Sum of differences: ")
		fmt.Println(sum_diff)

		// Part 2:
		similarity_score := 0

		rcounts := map[int]int{}
		for _, v := range rightList {
			rcounts[v]++
		}

		for _, v := range leftList {
			similarity_score += v * rcounts[v]
		}

		fmt.Print("Similarity score: ")
		fmt.Println(similarity_score)
	}
}
