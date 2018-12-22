package main

import (
	"fmt"
)

const (
	puzzleDepth   = 9171
	puzzleTargetX = 7
	puzzleTargetY = 721
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
}
