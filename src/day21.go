package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type state struct {
	fname string
	codes [5]string
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

	i := 0
	for scanner.Scan() {
		s.codes[i] = scanner.Text()
		i++
	}
	if i != 5 {
		log.Fatal("too few lines ", i)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/21/sample",
		"input/21/input",
	} {
		s := readFile(fname)
		if len(s.codes[0]) == 0 {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)
		sum := 0
		for _, v := range s.codes {
			numericPart, err := strconv.Atoi(v[:len(v)-1])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(v, "(", len(v), ")")
			d := driveNumeric(v)
			fmt.Println(d, "(", len(d), ")")
			d = driveDirectional(d)
			fmt.Println(d, "(", len(d), ")")
			d = driveDirectional(d)
			fmt.Println(d, "(", len(d), ")")
			sum += len(d) * numericPart
			fmt.Println()
		}
		fmt.Println()
		fmt.Println("sum complexity", sum)
	}
}

var dd map[string]string = map[string]string{
	"A^": "<",
	"A>": "v",
	"Av": "v<",
	"A<": "v<<",

	"<A": ">>^",
	"<^": ">^",
	"<v": ">",
	"<>": ">>",

	"v^": "^",
	"v>": ">",
	"v<": "<",
	"vA": ">^",

	"^A": ">",
	"^<": "v<",
	"^v": "v",
	"^>": "v>",

	">A": "^",
	">^": "^<",
	">v": "<",
	"><": "<<",

	"AA": "",
	"^^": "",
	"vv": "",
	"<<": "",
	">>": "",
}

func driveDirectional(s string) string {
	//  ^A
	// <v>
	pos := "A"
	rv := ""
	for _, v := range s {
		if i, ok := dd[pos+string(v)]; !ok {
			log.Fatal("dd fail ", s, ":", pos, string(v))
		} else {
			rv += i + "A"
			pos = string(v)
		}
	}

	return rv
}

var nd map[string]string = map[string]string{
	"A0": "<",
	"A1": "^<<",
	"A2": "^<",
	"A3": "^",
	"A4": "^^<<",
	"A5": "^^<",
	"A6": "^^",
	"A7": "^^^<<",
	"A8": "^^^<",
	"A9": "^^^",

	"0A": ">",
	"01": "^<",
	"02": "^",
	"03": "^>",
	"04": "^^<",
	"05": "^^",
	"06": "^^>",
	"07": "^^^<",
	"08": "^^^",
	"09": "^^^>",

	"17": "^^",

	"2A": "v>",
	"29": ">^^",

	"3A": "v",
	"37": "^^<<",
	"38": "^^<",
	"39": "^^",

	"45": ">",
	"46": ">>",

	"56": ">",

	"6A": "vv",
	"63": "v",

	"74": "v",
	"75": "v>",
	"76": "v>>",
	"78": ">",
	"79": ">>",

	"80": "vvv",
	"81": "vv<",
	"82": "vv",
	"83": "vv>",
	"87": "<",
	"89": ">",

	"9A": "vvv",
	"96": "v",
	"97": "<<",
	"98": "<",
}

func driveNumeric(s string) string {
	// 789
	// 456
	// 123
	//  0A
	pos := "A"
	rv := ""
	for _, c := range s {
		if s, ok := nd[pos+string(c)]; ok {
			rv += s + "A"
			pos = string(c)
		} else {
			log.Fatal(pos, string(c), ":", s)
		}
	}

	return rv
}
