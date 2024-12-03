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
		"input/02/sample",
		"input/02/input",
	} {
		fmt.Println(fname)

		file, err := os.Open(fname)
		if err != nil {
			log.Print(err)
			continue
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)

		reports := [][]int{}
		for scanner.Scan() {
			line := scanner.Text()
			split := strings.Split(line, " ")

			report := []int{}
			for _, v := range split {
				i, err := strconv.Atoi(v)
				if err != nil {
					log.Fatal(err)
				}
				report = append(report, i)
			}
			reports = append(reports, report)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		// Part 1:
		safes := 0
		for _, report := range reports {
			safe := isSafe(report)
			if safe {
				safes++
			}
		}

		fmt.Print("Safe: ")
		fmt.Println(safes)
	}
}

func isSafe(report []int) bool {
	asc, desc := true, true
	last := report[0]
	for i, v := range report {
		if i == 0 {
			continue
		}
		if last <= v {
			desc = false
		}
		if last >= v {
			asc = false
		}
		diff := last - v
		if diff < 0 {
			diff = -diff
		}
		if diff > 3 {
			asc = false
			desc = false
		}

		last = v
	}
	return asc || desc
}
