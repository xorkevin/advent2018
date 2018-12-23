package main

import (
	"container/heap"
	"fmt"
)

const (
	puzzleDepth   = 9171
	puzzleTargetX = 7
	puzzleTargetY = 721
	//puzzleDepth   = 510
	//puzzleTargetX = 10
	//puzzleTargetY = 10
)

type (
	Pos struct {
		x, y int
	}
)

func (p *Pos) GeoIndex(cache map[Pos]int) int {
	if v, ok := cache[*p]; ok {
		return v
	}
	if p.x == 0 {
		k := p.y * 48271
		cache[*p] = k
		return k
	} else if p.y == 0 {
		k := p.x * 16807
		cache[*p] = k
		return k
	} else if p.x == puzzleTargetX && p.y == puzzleTargetY {
		return 0
	}
	a := Pos{p.x - 1, p.y}
	b := Pos{p.x, p.y - 1}
	k := a.Erosion(cache) * b.Erosion(cache)
	cache[*p] = k
	return k
}

func (p *Pos) Erosion(cache map[Pos]int) int {
	return (p.GeoIndex(cache) + puzzleDepth) % 20183
}

const (
	regionRocky  = 0
	regionWet    = 1
	regionNarrow = 2
)

func (p *Pos) Region(cache map[Pos]int) int {
	return p.Erosion(cache) % 3
}

func absVal(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (p Pos) Mnhttn(other Pos) int {
	return absVal(other.x-p.x) + absVal(other.y-p.y)
}

func (p *Pos) Up() Pos {
	return Pos{
		x: p.x,
		y: p.y - 1,
	}
}
func (p *Pos) Down() Pos {
	return Pos{
		x: p.x,
		y: p.y + 1,
	}
}
func (p *Pos) Left() Pos {
	return Pos{
		x: p.x - 1,
		y: p.y,
	}
}
func (p *Pos) Right() Pos {
	return Pos{
		x: p.x + 1,
		y: p.y,
	}
}

type (
	Node struct {
		pos   Pos
		g, h  int
		item  int
		index int
	}
)

func (n *Node) f() int {
	return n.g + n.h
}

type (
	PosSet map[Pos]struct{}
)

func NewPosSet() PosSet {
	return PosSet{}
}
func (s PosSet) Has(a Pos) bool {
	_, ok := s[a]
	return ok
}
func (s PosSet) Add(a Pos) {
	s[a] = struct{}{}
}

type (
	PosHeap struct {
		list   []Pos
		scores map[Pos]*Node
	}
)

func NewPosHeap() PosHeap {
	return PosHeap{
		list:   []Pos{},
		scores: map[Pos]*Node{},
	}
}
func (h PosHeap) Len() int {
	return len(h.list)
}
func (h PosHeap) Less(i, j int) bool {
	return h.scores[h.list[i]].f() < h.scores[h.list[j]].f()
}
func (h PosHeap) Swap(i, j int) {
	h.list[i], h.list[j] = h.list[j], h.list[i]
	h.scores[h.list[i]].index = i
	h.scores[h.list[j]].index = j
}
func (h *PosHeap) Push(x interface{}) {
	h.list = append(h.list, x.(Pos))
}
func (h *PosHeap) Pop() interface{} {
	l := len(h.list)
	node := h.list[l-1]
	h.list = h.list[0 : l-1]
	return node
}
func (h *PosHeap) Update(n *Node) {
	n.index = h.scores[n.pos].index
	h.scores[n.pos] = n
	heap.Fix(h, n.index)
}
func (h *PosHeap) Add(n *Node) {
	n.index = len(h.list)
	h.scores[n.pos] = n
	heap.Push(h, n.pos)
}
func (h *PosHeap) Remove() (Pos, int, int) {
	if len(h.list) == 0 {
		return Pos{}, -1, -1
	}
	k := heap.Pop(h).(Pos)
	node := h.scores[k]
	delete(h.scores, k)
	return node.pos, node.g, node.item
}
func (h *PosHeap) Has(a Pos) bool {
	_, ok := h.scores[a]
	return ok
}

// 0: neither
// 1: torch
// 2: climbing gear
func canUseItem(ra, rb int, item int) bool {
	return ra != item && rb != item
}

func commonItem(ra, rb int) int {
	if canUseItem(ra, rb, 0) {
		return 0
	}
	if canUseItem(ra, rb, 1) {
		return 1
	}
	if canUseItem(ra, rb, 2) {
		return 2
	}
	panic("invalid state")
	return -1
}

func heuristic(a, b Pos, item int, cache map[Pos]int) int {
	k := a.Mnhttn(b)
	if canUseItem(a.Region(cache), b.Region(cache), item) {
		return k
	}
	return k + 7
}

func isInBounds(a Pos) bool {
	return a.x > -1 && a.y > -1
}
func neighbors(a Pos) []Pos {
	k := make([]Pos, 0, 4)
	up := a.Up()
	down := a.Down()
	left := a.Left()
	right := a.Right()
	if isInBounds(up) {
		k = append(k, up)
	}
	if isInBounds(down) {
		k = append(k, down)
	}
	if isInBounds(left) {
		k = append(k, left)
	}
	if isInBounds(right) {
		k = append(k, right)
	}
	return k
}

func astar(start, goal Pos, cache map[Pos]int) (int, int) {
	path := map[Pos]Pos{}
	closed := NewPosSet()
	open := NewPosHeap()
	open.Add(&Node{
		pos:  start,
		g:    0,
		h:    heuristic(start, goal, 1, cache),
		item: 1,
	})
	for current, cost, item := open.Remove(); cost > -1; current, cost, item = open.Remove() {
		if current == goal {
			routeCost := 0
			for k, ok := path[current]; ok; k, ok = path[k] {
				routeCost++
			}
			return cost, routeCost
		}
		closed.Add(current)
		for _, i := range neighbors(current) {
			if closed.Has(i) {
				continue
			}
			nitem := item
			ncost := cost + 1
			ra := current.Region(cache)
			rb := i.Region(cache)
			if ra != rb && !canUseItem(ra, rb, nitem) {
				nitem = commonItem(ra, rb)
				ncost += 7
			}
			n := &Node{
				pos:  i,
				g:    ncost,
				h:    heuristic(i, goal, nitem, cache),
				item: nitem,
			}
			if !open.Has(i) {
				path[i] = current
				open.Add(n)
			} else if ncost < open.scores[i].g {
				path[i] = current
				open.Update(n)
			}
		}
	}
	return -1, -1
}

func main() {
	part1 := 0
	cache := map[Pos]int{}
	for y := 0; y <= puzzleTargetY; y++ {
		for x := 0; x <= puzzleTargetX; x++ {
			k := Pos{x, y}
			part1 += k.Region(cache)
		}
	}
	fmt.Println(part1)
	fmt.Println(astar(Pos{0, 0}, Pos{puzzleTargetX, puzzleTargetY}, cache))
}
