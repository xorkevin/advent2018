package main

import (
	"bufio"
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
)

func (a *Pos) Down() Pos {
	return Pos{
		x: a.x,
		y: a.y + 1,
	}
}

func (a *Pos) Left() Pos {
	return Pos{
		x: a.x - 1,
		y: a.y,
	}
}

func (a *Pos) Right() Pos {
	return Pos{
		x: a.x + 1,
		y: a.y,
	}
}

type (
	SimLine struct {
		row        bool
		axis       int
		start, end int
	}

	Sim struct {
		board [][]byte
	}
)

func (s *Sim) At(pos Pos) byte {
	return s.board[pos.y][pos.x]
}

func (s *Sim) Emplace(pos Pos, c byte) {
	s.board[pos.y][pos.x] = c
}

func (s *Sim) IsInBounds(pos Pos) bool {
	return pos.y > -1 && pos.y < len(s.board) && pos.x > -1 && pos.x < len(s.board[0])
}

func (s *Sim) IsWall(pos Pos) bool {
	return s.board[pos.y][pos.x] == '#'
}

func (s *Sim) IsWater(pos Pos) bool {
	return s.board[pos.y][pos.x] == '|' || s.board[pos.y][pos.x] == '~'
}

func (s *Sim) IsStanding(pos Pos) bool {
	return s.board[pos.y][pos.x] == '#' || s.board[pos.y][pos.x] == '~'
}

func (s *Sim) Flow(source Pos) {
	if !s.IsInBounds(source) {
		return
	}
	if s.IsStanding(source) {
		return
	}
	if s.IsWater(source) {
		return
	}

	s.Emplace(source, '|')

	down := source.Down()
	s.Flow(down)
	if !s.IsInBounds(down) || !s.IsStanding(down) {
		return
	}

	left := source.Left()
	for !s.IsStanding(left) && s.IsStanding(left.Down()) {
		s.Emplace(left, '|')
		left = left.Left()
	}
	right := source.Right()
	for !s.IsStanding(right) && s.IsStanding(right.Down()) {
		s.Emplace(right, '|')
		right = right.Right()
	}

	if s.IsStanding(left) && s.IsStanding(right) {
		s.Emplace(source, '~')
		left := source.Left()
		for !s.IsStanding(left) && s.IsStanding(left.Down()) {
			s.Emplace(left, '~')
			left = left.Left()
		}
		right := source.Right()
		for !s.IsStanding(right) && s.IsStanding(right.Down()) {
			s.Emplace(right, '~')
			right = right.Right()
		}
		return
	}
	if !s.IsStanding(left) {
		s.Flow(left)
	}
	if !s.IsStanding(right) {
		s.Flow(right)
	}
}

func (s *Sim) Score() (int, int) {
	count := 0
	count2 := 0
	for _, row := range s.board {
		for _, i := range row {
			if i == '~' {
				count++
				count2++
			} else if i == '|' {
				count++
			}
		}
	}
	return count, count2
}

func (s *Sim) Print(nums bool) {
	if nums {
		fmt.Print("     ")
		for i := 0; i < len(s.board[0]); i++ {
			fmt.Print(i % 10)
		}
		fmt.Println()
	}
	for n, i := range s.board {
		if nums {
			fmt.Printf("%4d ", n)
		}
		fmt.Println(string(i))
	}
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

	lines := []SimLine{}
	miny := 99999999
	maxy := 0
	minx := 99999999
	maxx := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := false
		var a, s, e int
		if line[0] == 'x' {
			fmt.Sscanf(scanner.Text(), "x=%d, y=%d..%d", &a, &s, &e)
			if a < minx {
				minx = a
			}
			if a > maxx {
				maxx = a
			}
			if s < miny {
				miny = s
			}
			if e > maxy {
				maxy = e
			}
		} else {
			row = true
			fmt.Sscanf(scanner.Text(), "y=%d, x=%d..%d", &a, &s, &e)
			if a < miny {
				miny = a
			}
			if a > maxy {
				maxy = a
			}
			if s < minx {
				minx = s
			}
			if e > maxx {
				maxx = e
			}
		}
		lines = append(lines, SimLine{
			row:   row,
			axis:  a,
			start: s,
			end:   e,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	height := maxy - miny + 2
	width := maxx - minx + 3

	board := make([][]byte, height)
	for i := range board {
		k := make([]byte, width)
		for j := range k {
			k[j] = '.'
		}
		board[i] = k
	}

	for _, line := range lines {
		if line.row {
			for i := line.start; i <= line.end; i++ {
				board[line.axis-miny+1][i-minx+1] = '#'
			}
		} else {
			for i := line.start; i <= line.end; i++ {
				board[i-miny+1][line.axis-minx+1] = '#'
			}
		}
	}

	source := Pos{
		x: 500 - minx + 1,
		y: 0,
	}

	sim := Sim{
		board: board,
	}

	sim.Flow(source)
	sim.Emplace(source, '+')
	fmt.Println(sim.Score())
}
