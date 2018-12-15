package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

const (
	puzzleInput = "input.txt"
)

type (
	Pos struct {
		x, y int
	}

	Entity struct {
		elf bool
		pos Pos
	}

	Board [][]byte

	Game struct {
		elfs    map[Pos]*Entity
		goblins map[Pos]*Entity
		board   Board
	}

	EntityByPos []*Entity
)

func (p Pos) Less(other Pos) bool {
	if p.y == other.y {
		return p.x < other.x
	}
	return p.y < other.y
}

func (s EntityByPos) Len() int {
	return len(s)
}
func (s EntityByPos) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s EntityByPos) Less(i, j int) bool {
	return s[i].pos.Less(s[j].pos)
}

func (b Board) IsFree(pos Pos) bool {
	return b[pos.y][pos.x] == byte('.')
}

func (b Board) AdjacentFree(pos Pos) []Pos {
	k := []Pos{}
	i := Pos{
		x: pos.x - 1,
		y: pos.y,
	}
	if b.IsFree(i) {
		k = append(k, i)
	}
	i = Pos{
		x: pos.x + 1,
		y: pos.y,
	}
	if b.IsFree(i) {
		k = append(k, i)
	}
	i = Pos{
		x: pos.x,
		y: pos.y - 1,
	}
	if b.IsFree(i) {
		k = append(k, i)
	}
	i = Pos{
		x: pos.x,
		y: pos.y + 1,
	}
	if b.IsFree(i) {
		k = append(k, i)
	}
	return k
}

func (b Board) IsEnemy(pos Pos, elf bool) bool {
	if elf {
		return b[pos.y][pos.x] == byte('G')
	}
	return b[pos.y][pos.x] == byte('E')
}

func (b Board) AdjacentEnemy(pos Pos, elf bool) []Pos {
	k := []Pos{}
	i := Pos{
		x: pos.x - 1,
		y: pos.y,
	}
	if b.IsEnemy(i, elf) {
		k = append(k, i)
	}
	i = Pos{
		x: pos.x + 1,
		y: pos.y,
	}
	if b.IsEnemy(i, elf) {
		k = append(k, i)
	}
	i = Pos{
		x: pos.x,
		y: pos.y - 1,
	}
	if b.IsEnemy(i, elf) {
		k = append(k, i)
	}
	i = Pos{
		x: pos.x,
		y: pos.y + 1,
	}
	if b.IsEnemy(i, elf) {
		k = append(k, i)
	}
	return k
}

func NewEntity(pos Pos, elf bool) *Entity {
	return &Entity{
		pos: pos,
		elf: elf,
	}
}

func (e *Entity) Tick(enemies map[Pos]*Entity, board Board) {
}

func NewGame(elfs, goblins map[Pos]*Entity, board [][]byte) *Game {
	return &Game{
		elfs:    elfs,
		goblins: goblins,
		board:   Board(board),
	}
}

func (g *Game) Tick() {
	all := make([]*Entity, 0, len(g.elfs)+len(g.goblins))
	for _, v := range g.elfs {
		all = append(all, v)
	}
	for _, v := range g.goblins {
		all = append(all, v)
	}
	sort.Sort(EntityByPos(all))
	for _, i := range all {
		// if i is dead, continue

		if i.elf {
			i.Tick(g.goblins, g.board)
		} else {
			i.Tick(g.elfs, g.board)
		}
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

	board := [][]byte{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		board = append(board, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	elfs := map[Pos]*Entity{}
	goblins := map[Pos]*Entity{}

	for y, row := range board {
		for x, i := range row {
			switch i {
			case byte('E'):
				pos := Pos{
					x: x,
					y: y,
				}
				elfs[pos] = NewEntity(pos, true)
			case byte('G'):
				pos := Pos{
					x: x,
					y: y,
				}
				goblins[pos] = NewEntity(pos, false)
			}
		}
	}

	game := NewGame(elfs, goblins, board)
	fmt.Println("len", len(game.elfs), len(game.goblins))
}
