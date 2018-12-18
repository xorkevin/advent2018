package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

type (
	Pos struct {
		x, y int
	}
)

func (p *Pos) Adjacent() []Pos {
	x := p.x
	y := p.y
	return []Pos{
		Pos{x - 1, y - 1},
		Pos{x, y - 1},
		Pos{x + 1, y - 1},
		Pos{x + 1, y},
		Pos{x + 1, y + 1},
		Pos{x, y + 1},
		Pos{x - 1, y + 1},
		Pos{x - 1, y},
	}
}

type (
	Board [][]byte
)

func (b Board) At(pos Pos) byte {
	if pos.x > -1 && pos.x < len(b[0]) && pos.y > -1 && pos.y < len(b) {
		return b[pos.y][pos.x]
	}
	return 0
}

func (b Board) CellNext(pos Pos) byte {
	switch b.At(pos) {
	case '.':
		count := 0
		for _, i := range pos.Adjacent() {
			if b.At(i) == '|' {
				count++
			}
		}
		if count > 2 {
			return '|'
		}
		return '.'
	case '|':
		count := 0
		for _, i := range pos.Adjacent() {
			if b.At(i) == '#' {
				count++
			}
		}
		if count > 2 {
			return '#'
		}
		return '|'
	case '#':
		countL := 0
		countT := 0
		for _, i := range pos.Adjacent() {
			switch b.At(i) {
			case '#':
				countL++
			case '|':
				countT++
			}
		}
		if countL > 0 && countT > 0 {
			return '#'
		}
		return '.'
	}
	return b.At(pos)
}

func (b Board) NextState() Board {
	next := make(Board, len(b))
	for y, row := range b {
		line := make([]byte, len(row))
		for x, _ := range row {
			line[x] = b.CellNext(Pos{x, y})
		}
		next[y] = line
	}
	return next
}

func (b Board) Print() {
	for _, i := range b {
		fmt.Println(string(i))
	}
}

func (b Board) Score() int {
	countL := 0
	countT := 0
	for _, row := range b {
		for _, i := range row {
			switch i {
			case '#':
				countL++
			case '|':
				countT++
			}
		}
	}
	return countL * countT
}

func (b Board) String() string {
	k := strings.Builder{}
	for _, row := range b {
		k.Write(row)
	}
	return k.String()
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

	board := Board{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		board = append(board, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	seen := map[string]int{}
	seen[board.String()] = 0

	for i := 1; i <= 10; i++ {
		board = board.NextState()
		seen[board.String()] = i
	}
	fmt.Println(board.Score())

	part2Iterations := 1000000000

	for i := 11; i <= part2Iterations; i++ {
		board = board.NextState()
		k := board.String()
		if n, ok := seen[k]; ok {
			delta := i - n
			i += (part2Iterations - i) / delta * delta
		}
		seen[k] = i
	}

	fmt.Println(board.Score())
}
