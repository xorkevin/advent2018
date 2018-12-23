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
		ipos  IPos
		g, h  int
		index int
	}
)

func (n *Node) f() int {
	return n.g + n.h
}

type (
	IPos struct {
		pos  Pos
		item int
	}
	PosSet map[IPos]struct{}
)

func NewPosSet() PosSet {
	return PosSet{}
}
func (s PosSet) Has(a IPos) bool {
	_, ok := s[a]
	return ok
}
func (s PosSet) Add(a IPos) {
	s[a] = struct{}{}
}

type (
	PosHeap struct {
		list   []IPos
		scores map[IPos]*Node
	}
)

func NewPosHeap() PosHeap {
	return PosHeap{
		list:   []IPos{},
		scores: map[IPos]*Node{},
	}
}
func (h PosHeap) Len() int {
	return len(h.list)
}
func (h PosHeap) Less(i, j int) bool {
	ki := h.scores[h.list[i]]
	kj := h.scores[h.list[j]]
	fi := ki.f()
	fj := kj.f()
	if fi == fj {
		return ki.g < kj.g
	}
	return fi < fj
}
func (h PosHeap) Swap(i, j int) {
	h.list[i], h.list[j] = h.list[j], h.list[i]
	h.scores[h.list[i]].index = i
	h.scores[h.list[j]].index = j
}
func (h *PosHeap) Push(x interface{}) {
	h.list = append(h.list, x.(IPos))
}
func (h *PosHeap) Pop() interface{} {
	l := len(h.list)
	node := h.list[l-1]
	h.list = h.list[0 : l-1]
	return node
}
func (h *PosHeap) Update(n *Node) {
	n.index = h.scores[n.ipos].index
	h.scores[n.ipos] = n
	heap.Fix(h, n.index)
}
func (h *PosHeap) Add(n *Node) {
	n.index = len(h.list)
	h.scores[n.ipos] = n
	heap.Push(h, n.ipos)
}
func (h *PosHeap) Remove() (IPos, int) {
	if len(h.list) == 0 {
		return IPos{}, -1
	}
	k := heap.Pop(h).(IPos)
	node := h.scores[k]
	delete(h.scores, k)
	return node.ipos, node.g
}
func (h *PosHeap) Has(a IPos) bool {
	_, ok := h.scores[a]
	return ok
}

// 0: neither
// 1: torch
// 2: climbing gear
func canUseItem(ra, rb int, item int) bool {
	return ra != item && rb != item
}

func otherItem(region, item int) int {
	if canUseItem(region, item, 0) {
		return 0
	}
	if canUseItem(region, item, 1) {
		return 1
	}
	return 2
}

func heuristic(a, b Pos, cache map[Pos]int) int {
	return a.Mnhttn(b)
}

func isInBounds(a Pos) bool {
	return a.x > -1 && a.y > -1
}
func neighbors(a IPos) []IPos {
	k := make([]IPos, 0, 4)
	up := a.pos.Up()
	down := a.pos.Down()
	left := a.pos.Left()
	right := a.pos.Right()
	if isInBounds(up) {
		k = append(k, IPos{
			pos:  up,
			item: a.item,
		})
	}
	if isInBounds(down) {
		k = append(k, IPos{
			pos:  down,
			item: a.item,
		})
	}
	if isInBounds(left) {
		k = append(k, IPos{
			pos:  left,
			item: a.item,
		})
	}
	if isInBounds(right) {
		k = append(k, IPos{
			pos:  right,
			item: a.item,
		})
	}
	return k
}

func astar(start, goal Pos, cache map[Pos]int) (int, int) {
	path := map[IPos]IPos{}
	closed := NewPosSet()
	open := NewPosHeap()
	open.Add(&Node{
		ipos: IPos{
			pos:  start,
			item: 1,
		},
		g: 0,
		h: heuristic(start, goal, cache),
	})
	for current, cost := open.Remove(); cost > -1; current, cost = open.Remove() {
		if current.pos == goal {
			routeCost := 0
			for k, ok := path[current]; ok; k, ok = path[k] {
				routeCost++
			}
			return cost, routeCost
		}
		closed.Add(current)
		ra := current.pos.Region(cache)
		for _, i := range neighbors(current) {
			if closed.Has(i) {
				continue
			}
			rb := i.pos.Region(cache)
			if !canUseItem(ra, rb, current.item) {
				continue
			}
			n := &Node{
				ipos: i,
				g:    cost + 1,
				h:    heuristic(i.pos, goal, cache),
			}
			if !open.Has(i) {
				path[i] = current
				open.Add(n)
			} else if cost+1 < open.scores[i].g {
				path[i] = current
				open.Update(n)
			}
		}
		k := IPos{
			pos:  current.pos,
			item: otherItem(ra, current.item),
		}
		if !closed.Has(k) {
			n := &Node{
				ipos: k,
				g:    cost + 7,
				h:    heuristic(current.pos, goal, cache),
			}
			if !open.Has(k) {
				path[k] = current
				open.Add(n)
			} else if cost+7 < open.scores[k].g {
				path[k] = current
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
