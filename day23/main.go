package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
)

const (
	puzzleInput = "input.txt"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type (
	Pos struct {
		x, y, z int
	}

	Bot struct {
		pos    Pos
		radius int
	}
)

func (p Pos) Mnhttn(other Pos) int {
	return abs(other.x-p.x) + abs(other.y-p.y) + abs(other.z-p.z)
}

type (
	Cube struct {
		pos  Pos
		size int
	}
)

func axisDist(p, lower, upper int) int {
	if p < lower {
		return lower - p
	}
	if p > upper {
		return p - upper
	}
	return 0
}

func (c Cube) Dist(pos Pos) int {
	return axisDist(pos.x, c.pos.x, c.pos.x+c.size-1) + axisDist(pos.y, c.pos.y, c.pos.y+c.size-1) + axisDist(pos.z, c.pos.z, c.pos.z+c.size-1)
}

func (c Cube) Divide() []Cube {
	x := c.pos.x
	y := c.pos.y
	z := c.pos.z
	size := c.size / 2
	xm := x + size
	ym := y + size
	zm := z + size
	return []Cube{
		Cube{
			pos:  Pos{x, y, z},
			size: size,
		},
		Cube{
			pos:  Pos{xm, y, z},
			size: size,
		},
		Cube{
			pos:  Pos{x, ym, z},
			size: size,
		},
		Cube{
			pos:  Pos{xm, ym, z},
			size: size,
		},
		Cube{
			pos:  Pos{x, y, zm},
			size: size,
		},
		Cube{
			pos:  Pos{xm, y, zm},
			size: size,
		},
		Cube{
			pos:  Pos{x, ym, zm},
			size: size,
		},
		Cube{
			pos:  Pos{xm, ym, zm},
			size: size,
		},
	}
}

func botRange(c Cube, bots []Bot) int {
	k := 0
	for _, i := range bots {
		if c.Dist(i.pos) <= i.radius {
			k++
		}
	}
	return k
}

type (
	Score struct {
		cube  Cube
		score int
		index int
	}

	ScoreSlice []Score
)

func (s ScoreSlice) Len() int {
	return len(s)
}
func (s ScoreSlice) Less(i, j int) bool {
	si := s[i]
	sj := s[j]
	if si.score == sj.score {
		di := si.cube.pos.Mnhttn(Pos{})
		dj := sj.cube.pos.Mnhttn(Pos{})
		if di == dj {
			return si.cube.size < sj.cube.size
		}
		return di < dj
	}
	return si.score > sj.score
}
func (s ScoreSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
	s[i].index = i
	s[j].index = j
}
func (s *ScoreSlice) Push(x interface{}) {
	i := x.(Score)
	i.index = len(*s)
	*s = append(*s, i)
}
func (s *ScoreSlice) Pop() interface{} {
	l := len(*s)
	i := (*s)[l-1]
	i.index = -1
	*s = (*s)[0 : l-1]
	return i
}
func (s *ScoreSlice) Add(score Score) {
	heap.Push(s, score)
}
func (s *ScoreSlice) Remove() (Score, bool) {
	if s.Len() == 0 {
		return Score{}, false
	}
	return heap.Pop(s).(Score), true
}

func searchCube(start Cube, bots []Bot) Pos {
	nodeCount := 0
	open := ScoreSlice{}
	open.Add(Score{
		cube:  start,
		score: botRange(start, bots),
	})
	for current, ok := open.Remove(); ok; current, ok = open.Remove() {
		nodeCount++
		if current.cube.size < 2 {
			return current.cube.pos
		}
		for _, i := range current.cube.Divide() {
			open.Add(Score{
				cube:  i,
				score: botRange(i, bots),
			})
		}
	}
	return Pos{}
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

	maxRadius := 0
	maxRadBot := Pos{}
	bots := []Bot{}
	scanner := bufio.NewScanner(file)
	minPos := math.MaxInt32
	maxPos := math.MinInt32
	for scanner.Scan() {
		var x, y, z, r int
		fmt.Sscanf(scanner.Text(), "pos=<%d,%d,%d>, r=%d", &x, &y, &z, &r)
		pos := Pos{x, y, z}
		bots = append(bots, Bot{
			pos:    pos,
			radius: r,
		})
		if r > maxRadius {
			maxRadius = r
			maxRadBot = pos
		}
		if x < minPos {
			minPos = x
		}
		if y < minPos {
			minPos = y
		}
		if z < minPos {
			minPos = z
		}
		if x > maxPos {
			maxPos = x
		}
		if y > maxPos {
			maxPos = y
		}
		if z > maxPos {
			maxPos = z
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	part1 := 0
	for _, i := range bots {
		if maxRadBot.Mnhttn(i.pos) <= maxRadius {
			part1++
		}
	}
	fmt.Println(part1)

	fmt.Println(minPos, maxPos)
	corner := -1 << 28
	cornerSize := 1 << 29
	fmt.Println(corner, corner+cornerSize)

	start := Cube{
		pos:  Pos{corner, corner, corner},
		size: cornerSize,
	}
	fmt.Println(searchCube(start, bots).Mnhttn(Pos{}))
}
