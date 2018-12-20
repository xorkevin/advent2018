package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
)

const (
	puzzleInput = "input.txt"
)

type (
	Pos struct {
		x, y int
	}

	Board struct {
		grid map[Pos]map[Pos]struct{}
	}
)

func (p Pos) Up() Pos {
	return Pos{
		x: p.x,
		y: p.y - 1,
	}
}
func (p Pos) Down() Pos {
	return Pos{
		x: p.x,
		y: p.y + 1,
	}
}
func (p Pos) Left() Pos {
	return Pos{
		x: p.x - 1,
		y: p.y,
	}
}
func (p Pos) Right() Pos {
	return Pos{
		x: p.x + 1,
		y: p.y,
	}
}
func (p Pos) Less(other Pos) bool {
	if p.y == other.y {
		return p.x < other.x
	}
	return p.y < other.y
}

func NewBoard() Board {
	return Board{
		grid: map[Pos]map[Pos]struct{}{},
	}
}

func (b *Board) Update(pos Pos, other Pos) {
	if _, ok := b.grid[pos]; !ok {
		b.grid[pos] = map[Pos]struct{}{}
	}
	b.grid[pos][other] = struct{}{}
}

type (
	Path struct {
		pos    Pos
		branch int
	}
	Stack []Path
)

func (s *Stack) Push(a Path) {
	*s = append(*s, a)
}

func (s *Stack) Pop() Path {
	a := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return a
}

func (s *Stack) Peek() Path {
	return (*s)[len(*s)-1]
}

func (s *Stack) Len() int {
	return len(*s)
}

type (
	PosScore struct {
		pos   Pos
		score int
	}

	PosHeap struct {
		list   []Pos
		scores map[Pos]int
	}

	PosSet map[Pos]struct{}
)

func NewPosHeap() *PosHeap {
	return &PosHeap{
		list:   []Pos{},
		scores: map[Pos]int{},
	}
}

func (h PosHeap) Len() int {
	return len(h.list)
}
func (h PosHeap) Less(i, j int) bool {
	a := h.list[i]
	b := h.list[j]
	as := h.scores[a]
	bs := h.scores[b]
	if as == bs {
		return a.Less(b)
	}
	return as < bs
}
func (h PosHeap) Swap(i, j int) {
	h.list[i], h.list[j] = h.list[j], h.list[i]
}
func (h *PosHeap) Push(x interface{}) {
	h.list = append(h.list, x.(Pos))
}
func (h *PosHeap) Pop() interface{} {
	l := len(h.list)
	k := h.list[l-1]
	h.list = h.list[0 : l-1]
	return k
}
func (h *PosHeap) Add(g int, pos Pos) {
	h.scores[pos] = g
	heap.Push(h, pos)
}
func (h *PosHeap) Take() (Pos, int) {
	if h.Len() == 0 {
		return Pos{}, -1
	}
	k := heap.Pop(h).(Pos)
	ks := h.scores[k]
	delete(h.scores, k)
	return k, ks
}
func (h *PosHeap) Has(pos Pos) bool {
	_, ok := h.scores[pos]
	return ok
}

func (ps PosSet) Has(pos Pos) bool {
	_, ok := ps[pos]
	return ok
}
func (ps PosSet) Add(pos Pos) {
	ps[pos] = struct{}{}
}

func main() {
	file, err := os.Open(puzzleInput)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	board := NewBoard()
	path := Stack{}

	current := Pos{
		x: 0,
		y: 0,
	}
	for _, i := range line {
		switch i {
		case '^':
			continue
		case '$':
			break
		case '(':
			path.Push(Path{
				pos:    current,
				branch: 0,
			})
			path.Push(Path{
				pos:    current,
				branch: 1,
			})
		case '|':
			for prev := path.Pop(); prev.branch != 1; prev = path.Pop() {
			}
			current = path.Peek().pos
			path.Push(Path{
				pos:    current,
				branch: 1,
			})
		case ')':
			for prev := path.Pop(); prev.branch != 1; prev = path.Pop() {
			}
			current = path.Peek().pos
		case 'N':
			prev := current
			current = current.Up()
			path.Push(Path{
				pos:    prev,
				branch: 0,
			})
			board.Update(prev, current)
			board.Update(current, prev)
		case 'S':
			prev := current
			current = current.Down()
			path.Push(Path{
				pos:    prev,
				branch: 0,
			})
			board.Update(prev, current)
			board.Update(current, prev)
		case 'W':
			prev := current
			current = current.Left()
			path.Push(Path{
				pos:    prev,
				branch: 0,
			})
			board.Update(prev, current)
			board.Update(current, prev)
		case 'E':
			prev := current
			current = current.Right()
			path.Push(Path{
				pos:    prev,
				branch: 0,
			})
			board.Update(prev, current)
			board.Update(current, prev)
		}
	}

	closed := PosSet{}
	open := NewPosHeap()
	start := Pos{
		x: 0,
		y: 0,
	}
	open.Add(0, start)
	maxDist := 0
	lessHundred := PosSet{}
	for current, currentg := open.Take(); currentg > -1; current, currentg = open.Take() {
		if currentg > maxDist {
			maxDist = currentg
		}
		if currentg >= 1000 {
			lessHundred.Add(current)
		}
		closed.Add(current)
		for k, _ := range board.grid[current] {
			if !closed.Has(k) && !open.Has(k) {
				open.Add(currentg+1, k)
			}
		}
	}
	fmt.Println(maxDist)
	fmt.Println(len(lessHundred))
}
