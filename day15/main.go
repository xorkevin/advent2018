package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"sort"
)

const (
	puzzleInput = "input.txt"
)

func copyMap(m map[Pos]*Entity) map[Pos]*Entity {
	m2 := map[Pos]*Entity{}
	for k, v := range m {
		a := *v
		m2[k] = &a
	}
	return m2
}

type (
	Pos struct {
		x, y int
	}

	PosList []Pos
)

func (s PosList) Len() int {
	return len(s)
}
func (s PosList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s PosList) Less(i, j int) bool {
	return s[i].Less(s[j])
}

func (p Pos) Less(other Pos) bool {
	if p.y == other.y {
		return p.x < other.x
	}
	return p.y < other.y
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (p Pos) Manhattan(other Pos) int {
	return abs(other.x-p.x) + abs(other.y-p.y)
}

type (
	Entity struct {
		elf    bool
		pos    Pos
		health int
		attack int
		debug  bool
	}

	EntityByPos []*Entity
)

func NewEntity(pos Pos, elf bool) *Entity {
	return &Entity{
		pos:    pos,
		elf:    elf,
		health: 200,
		attack: 3,
	}
}

func (e *Entity) Dead() bool {
	return e.health < 1
}
func (e *Entity) Hit(damage int) bool {
	e.health -= damage
	return e.Dead()
}
func (e *Entity) Attack(other *Entity) bool {
	return other.Hit(e.attack)
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

type (
	Board [][]byte

	Game struct {
		elfs    map[Pos]*Entity
		goblins map[Pos]*Entity
		board   Board
		elfDied bool
	}
)

func NewGame(elfs, goblins map[Pos]*Entity, board [][]byte) *Game {
	return &Game{
		elfs:    elfs,
		goblins: goblins,
		board:   Board(board),
		elfDied: false,
	}
}

func (g *Game) IsFree(pos Pos) bool {
	_, hasElf := g.elfs[pos]
	_, hasGoblin := g.goblins[pos]
	return !hasElf && !hasGoblin && g.board[pos.y][pos.x] == byte('.')
}

func (g *Game) AdjacentFree(pos Pos) []Pos {
	k := []Pos{}
	i := Pos{
		x: pos.x,
		y: pos.y - 1,
	}
	if g.IsFree(i) {
		k = append(k, i)
	}
	i = Pos{
		x: pos.x - 1,
		y: pos.y,
	}
	if g.IsFree(i) {
		k = append(k, i)
	}
	i = Pos{
		x: pos.x + 1,
		y: pos.y,
	}
	if g.IsFree(i) {
		k = append(k, i)
	}
	i = Pos{
		x: pos.x,
		y: pos.y + 1,
	}
	if g.IsFree(i) {
		k = append(k, i)
	}
	return k
}

func (g *Game) IsEnemy(pos Pos, elf bool) bool {
	if elf {
		_, ok := g.goblins[pos]
		return ok
	}
	_, ok := g.elfs[pos]
	return ok
}

func (g *Game) AdjacentEnemy(pos Pos, elf bool) []Pos {
	k := []Pos{}
	i := Pos{
		x: pos.x,
		y: pos.y - 1,
	}
	if g.IsEnemy(i, elf) {
		k = append(k, i)
	}
	i = Pos{
		x: pos.x - 1,
		y: pos.y,
	}
	if g.IsEnemy(i, elf) {
		k = append(k, i)
	}
	i = Pos{
		x: pos.x + 1,
		y: pos.y,
	}
	if g.IsEnemy(i, elf) {
		k = append(k, i)
	}
	i = Pos{
		x: pos.x,
		y: pos.y + 1,
	}
	if g.IsEnemy(i, elf) {
		k = append(k, i)
	}
	return k
}

type (
	Score struct {
		g int
		h int
	}

	PosScore struct {
		pos   Pos
		score Score
	}

	PosHeap struct {
		start  Pos
		list   []Pos
		scores map[Pos]Score
	}

	PosSet map[Pos]struct{}
)

func (s Score) f() int {
	return s.g + s.h
}

func NewPosHeap(start Pos) *PosHeap {
	return &PosHeap{
		start:  start,
		list:   []Pos{},
		scores: map[Pos]Score{},
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
	af := as.f()
	bf := bs.f()
	if af == bf {
		if as.h == bs.h {
			return a.Less(b)
		}
		return as.h < bs.h
	}
	return af < bf
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
	h.scores[pos] = Score{
		g: g,
		h: 0,
	}
	heap.Push(h, pos)
}
func (h *PosHeap) Take() (Pos, int) {
	if h.Len() == 0 {
		return Pos{}, -1
	}
	k := heap.Pop(h).(Pos)
	ks := h.scores[k].g
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

func (g *Game) RemoveEntity(e *Entity) {
	if e.elf {
		delete(g.elfs, e.pos)
		g.elfDied = true
	} else {
		delete(g.goblins, e.pos)
	}
}

func (g *Game) EntityPath(start Pos, goal PosSet) (Pos, int) {
	closed := PosSet{}
	open := NewPosHeap(start)
	for k, _ := range goal {
		open.Add(0, k)
	}
	for current, currentg := open.Take(); currentg > -1; current, currentg = open.Take() {
		if start.Manhattan(current) < 2 {
			return current, currentg + 1
		}
		closed.Add(current)
		k := g.AdjacentFree(current)
		for _, i := range k {
			if !closed.Has(i) {
				if !open.Has(i) {
					open.Add(currentg+1, i)
				} else if currentg+1 < open.scores[i].g {
					for n, j := range open.list {
						if i == j {
							node := open.scores[i]
							node.g = currentg + 1
							open.scores[i] = node
							heap.Fix(open, n)
							break
						}
					}
				}
			}
		}
	}
	return Pos{}, -1
}

func (g *Game) EntityMove(e *Entity, enemies map[Pos]*Entity) {
	inRange := PosSet{}
	for k, _ := range enemies {
		for _, i := range g.AdjacentFree(k) {
			inRange.Add(i)
		}
	}

	next, cost := g.EntityPath(e.pos, inRange)
	if cost < 0 {
		return
	}

	if e.elf {
		delete(g.elfs, e.pos)
		g.elfs[next] = e
	} else {
		delete(g.goblins, e.pos)
		g.goblins[next] = e
	}
	e.pos = next
}

func (g *Game) EntityAttack(e *Entity, enemies []*Entity) {
	if len(enemies) == 0 {
		return
	}
	target := enemies[0]
	for _, i := range enemies[1:] {
		if i.health < target.health {
			target = i
		}
	}
	if dead := e.Attack(target); dead {
		g.RemoveEntity(target)
	}
}

func (g *Game) TickEntity(e *Entity, enemies map[Pos]*Entity) {
	adjacentEnemies := g.AdjacentEnemy(e.pos, e.elf)
	if len(adjacentEnemies) == 0 {
		g.EntityMove(e, enemies)
		adjacentEnemies = g.AdjacentEnemy(e.pos, e.elf)
	}
	adj := make([]*Entity, 0, len(adjacentEnemies))
	for _, i := range adjacentEnemies {
		adj = append(adj, enemies[i])
	}
	g.EntityAttack(e, adj)
}

func (g *Game) Tick(part2 bool) (bool, bool) {
	all := make([]*Entity, 0, len(g.elfs)+len(g.goblins))
	for _, v := range g.elfs {
		all = append(all, v)
	}
	for _, v := range g.goblins {
		all = append(all, v)
	}
	sort.Sort(EntityByPos(all))
	for _, i := range all {
		if i.Dead() {
			continue
		}

		if len(g.elfs) == 0 || len(g.goblins) == 0 {
			return true, false
		}

		if i.elf {
			g.TickEntity(i, g.goblins)
		} else {
			g.TickEntity(i, g.elfs)
		}
		if part2 && g.elfDied {
			return false, true
		}
	}
	return false, false
}

func (g *Game) Print() {
	board := make([][]byte, len(g.board))
	for i := range board {
		k := make([]byte, len(g.board[i]))
		copy(k, g.board[i])
		board[i] = k
	}
	for k := range g.elfs {
		board[k.y][k.x] = 'E'
	}
	for k := range g.goblins {
		board[k.y][k.x] = 'G'
	}
	fmt.Print("   ")
	for i := 0; i < len(board[0]); i++ {
		fmt.Print(i % 10)
	}
	fmt.Println()
	for n, i := range board {
		fmt.Printf("%2d ", n)
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
				board[y][x] = '.'
			case byte('G'):
				pos := Pos{
					x: x,
					y: y,
				}
				k := NewEntity(pos, false)
				goblins[pos] = k
				board[y][x] = '.'
			}
		}
	}

	debug := false

	game := NewGame(copyMap(elfs), copyMap(goblins), board)

	i := 0
	for done, _ := game.Tick(false); !done; done, _ = game.Tick(false) {
		i++
		if debug {
			fmt.Println(i)
			game.Print()
		}
	}
	totalHealth := 0
	for _, i := range game.elfs {
		totalHealth += i.health
	}
	for _, i := range game.goblins {
		totalHealth += i.health
	}
	fmt.Println("rounds:", i)
	fmt.Println(i * totalHealth)

	part2(elfs, goblins, board)
}

func part2(elfs, goblins map[Pos]*Entity, board [][]byte) {
	for attackVal := 4; ; attackVal++ {
		elfNum := len(elfs)
		e := copyMap(elfs)
		for _, v := range e {
			v.attack = attackVal
		}
		game := NewGame(e, copyMap(goblins), board)
		i := 0
		for done, elfDied := game.Tick(true); !done && !elfDied; done, elfDied = game.Tick(true) {
			i++
		}
		if len(game.elfs) != elfNum {
			continue
		}
		totalHealth := 0
		for _, i := range game.elfs {
			totalHealth += i.health
		}
		for _, i := range game.goblins {
			totalHealth += i.health
		}
		fmt.Println("attack:", attackVal)
		fmt.Println(i * totalHealth)
		break
	}
}
