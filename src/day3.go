package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	re := regexp.MustCompile("do\\(\\)|don't\\(\\)|mul\\([0-9]{1,3},[0-9]{1,3}\\)")

	for _, fname := range []string{
		"input/03/sample",
		"input/03/sample2",
		"input/03/input",
	} {
		fmt.Println(fname)

		b, err := os.ReadFile(fname)
		if err != nil {
			log.Print(err)
			continue
		}
		str := string(b)

		matches := re.FindAllString(str, -1)

		var c, c2 int
		enabled := true
		for _, match := range matches {
			if match == "do()" {
				enabled = true
				continue
			}
			if match == "don't()" {
				enabled = false
				continue
			}
			match = match[4:]
			match = match[:len(match)-1]
			args := strings.Split(match, ",")
			if len(args) != 2 {
				log.Fatal(fmt.Errorf("can't parse %v", match))
			}
			l, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatal(err)
			}
			r, err := strconv.Atoi(args[1])
			if err != nil {
				log.Fatal(err)
			}
			c += l * r
			if enabled {
				c2 += l * r
			}
		}

		// Part 1:
		fmt.Print("sum: ")
		fmt.Println(c)

		// Part 2:
		fmt.Print("sum (conditionals): ")
		fmt.Println(c2)
	}
}
