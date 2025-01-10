package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
)

type coord struct {
	x int
	y int
}

type reindeer struct {
	pos coord
	dir coord
}

type state struct {
	fname    string
	size     coord
	walls    map[coord]bool
	reindeer reindeer
	exit     coord
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
	s.walls = map[coord]bool{}

	y := 0
	for scanner.Scan() {
		line := scanner.Text()

		for x, v := range line {
			if v == '#' {
				s.walls[coord{x, y}] = true
			} else if v == 'S' {
				s.reindeer = reindeer{
					pos: coord{x, y},
					dir: coord{1, 0},
				}
			} else if v == 'E' {
				s.exit = coord{x, y}
			}
			if x > s.size.x {
				s.size.x = x
			}
		}

		s.size.y = y
		y++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	s.size.x++
	s.size.y++

	return s
}

func states() (rval []state) {
	for _, fname := range []string{
		"input/16/sample",
		"input/16/sample2",
		"input/16/input",
	} {
		s := readFile(fname)
		if len(s.walls) == 0 {
			continue
		}
		rval = append(rval, s)
	}
	return
}

func plus(a, b coord) coord {
	return coord{a.x + b.x, a.y + b.y}
}

func (s *state) isTile(pos coord) bool {
	invalid := pos.x < 0 || pos.y < 0 || pos.x >= s.size.x || pos.y >= s.size.y
	return !invalid && !s.walls[pos]
}

var directions = []coord{
	coord{-1, 0},
	coord{1, 0},
	coord{0, -1},
	coord{0, 1},
}

func (s *state) neighbors(k reindeer) map[reindeer]int {
	rv := map[reindeer]int{}

	if pos := plus(k.pos, k.dir); s.isTile(pos) {
		rv[reindeer{pos: pos, dir: k.dir}] = 1
	}

	if k.dir.x == 0 {
		rv[reindeer{pos: k.pos, dir: directions[0]}] = 1000
		rv[reindeer{pos: k.pos, dir: directions[1]}] = 1000
	} else {
		rv[reindeer{pos: k.pos, dir: directions[2]}] = 1000
		rv[reindeer{pos: k.pos, dir: directions[3]}] = 1000
	}

	return rv
}

type Item struct {
	value    reindeer
	priority int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	it := pq[i]
	jt := pq[j]
	return it.priority < jt.priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *Item, value reindeer, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

func (s *state) vertices() (rv []reindeer) {
	for y := 0; y < s.size.y; y++ {
		for x := 0; x < s.size.x; x++ {
			pos := coord{x, y}
			if s.isTile(pos) {
				for _, d := range directions {
					rv = append(rv, reindeer{pos: pos, dir: d})
				}
			}
		}
	}

	return
}

func (s *state) dijkstras(
	source reindeer,
) (
	dist map[reindeer]int,
	prev map[reindeer][]reindeer,
) {
	dist = map[reindeer]int{}
	prev = map[reindeer][]reindeer{}

	vertices := s.vertices()
	Q := make(PriorityQueue, len(vertices))
	inQ := map[reindeer]*Item{}

	for i, v := range vertices {
		prev[v] = nil
		if v == s.reindeer {
			dist[v] = 0
		} else {
			dist[v] = math.MaxInt
		}
		item := &Item{value: v, priority: dist[v], index: i}
		Q[i] = item
		inQ[v] = item
	}
	heap.Init(&Q)

	for len(Q) > 0 {
		u := heap.Pop(&Q).(*Item).value
		delete(inQ, u)

		for neighbor, marginalCost := range s.neighbors(u) {
			item := inQ[neighbor]
			if item == nil {
				continue
			}

			alt := dist[u] + marginalCost
			if alt < dist[neighbor] {
				dist[neighbor] = alt
				Q.update(item, item.value, alt)
				prev[neighbor] = []reindeer{u}
			} else if alt == dist[neighbor] {
				prev[neighbor] = append(prev[neighbor], u)
			}
		}

	}

	return
}

func main() {
	for _, s := range states() {
		fmt.Println(s.fname)

		dist, _ := s.dijkstras(s.reindeer)

		exits := []int{}
		for _, d := range directions {
			r := reindeer{pos: s.exit, dir: d}
			if v, ok := dist[r]; ok {
				exits = append(exits, v)
			}
		}
		if len(exits) == 0 {
			log.Fatal("i need an exit")
		}
		part1 := exits[0]
		for _, v := range exits {
			if v < part1 {
				part1 = v
			}
		}

		// Part 1
		fmt.Println("min score: ", part1)
	}
}
