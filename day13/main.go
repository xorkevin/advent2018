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

const (
	dirN = 0
	dirE = 1
	dirS = 2
	dirW = 3
)

type (
	Cart struct {
		x, y, dir int
		turn      int
	}

	CartList struct {
		carts CartSlice
		board [][]byte
	}

	CartSlice []*Cart
)

func (s CartSlice) Len() int {
	return len(s)
}
func (s CartSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s CartSlice) Less(i, j int) bool {
	if s[i].y == s[j].y {
		return s[i].x < s[j].x
	}
	return s[i].y < s[j].y
}

func cartCharToDir(b byte) int {
	switch b {
	case byte('^'):
		return dirN
	case byte('>'):
		return dirE
	case byte('v'):
		return dirS
	case byte('<'):
		return dirW
	default:
		return -1
	}
}

func cartDirToChar(dir int) byte {
	switch dir {
	case dirN:
		return byte('^')
	case dirE:
		return byte('>')
	case dirS:
		return byte('v')
	case dirW:
		return byte('<')
	default:
		panic("invalid dir")
	}
	return 8
}

func NewCart(x, y, dir int) *Cart {
	return &Cart{
		x:    x,
		y:    y,
		dir:  dir,
		turn: 0,
	}
}

func (c *Cart) Step(board [][]byte) {
	switch c.dir {
	case dirN:
		c.y--
	case dirE:
		c.x++
	case dirS:
		c.y++
	case dirW:
		c.x--
	}
	switch board[c.y][c.x] {
	case byte('\\'):
		switch c.dir {
		case dirN:
			c.dir = dirW
		case dirE:
			c.dir = dirS
		case dirS:
			c.dir = dirE
		case dirW:
			c.dir = dirN
		default:
			panic("invalid direction")
		}
	case byte('/'):
		switch c.dir {
		case dirN:
			c.dir = dirE
		case dirE:
			c.dir = dirN
		case dirS:
			c.dir = dirW
		case dirW:
			c.dir = dirS
		default:
			panic("invalid direction")
		}
	case byte('+'):
		switch c.turn {
		case 0:
			c.dir = (c.dir + 3) % 4
		case 1:
		case 2:
			c.dir = (c.dir + 1) % 4
		default:
			panic("invalid turn")
		}
		c.turn = (c.turn + 1) % 3
	}
}

func NewCartList(board [][]byte) *CartList {
	return &CartList{
		carts: CartSlice{},
		board: board,
	}
}

func (cl *CartList) AddCart(c *Cart) {
	cl.carts = append(cl.carts, c)
}

func (cl *CartList) Step() {
	crashed := map[int]struct{}{}
	for n, i := range cl.carts {
		if _, ok := crashed[n]; ok {
			continue
		}
		i.Step(cl.board)
		if crash := cl.HasCrash(i, crashed); crash > -1 {
			fmt.Println(i.x, i.y)
			crashed[n] = struct{}{}
			crashed[crash] = struct{}{}
		}
	}
	k := CartSlice{}
	for n, i := range cl.carts {
		if _, ok := crashed[n]; !ok {
			k = append(k, i)
		}
	}
	cl.carts = k
	sort.Sort(cl.carts)
}

func (cl *CartList) HasCrash(c *Cart, crashed map[int]struct{}) int {
	for n, i := range cl.carts {
		if _, ok := crashed[n]; ok {
			continue
		}
		if c != i && c.x == i.x && c.y == i.y {
			return n
		}
	}
	return -1
}

func (cl *CartList) Print() {
	board := make([][]byte, len(cl.board))
	for i := range board {
		board[i] = make([]byte, len(cl.board[i]))
		copy(board[i], cl.board[i])
	}
	for _, i := range cl.carts {
		if cartCharToDir(board[i.y][i.x]) > -1 {
			board[i.y][i.x] = byte('X')
		} else {
			board[i.y][i.x] = cartDirToChar(i.dir)
		}
	}
	for _, i := range board {
		fmt.Println(string(i))
	}
}

func (cl *CartList) PrintPos() {
	for _, i := range cl.carts {
		fmt.Printf("%d, %d; ", i.x, i.y)
	}
	fmt.Println()
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

	carts := NewCartList(board)
	for y, row := range board {
		for x, i := range row {
			if dir := cartCharToDir(i); dir > -1 {
				carts.AddCart(NewCart(x, y, dir))
				if dir == dirE || dir == dirW {
					board[y][x] = byte('-')
				} else if dir == dirN || dir == dirS {
					board[y][x] = byte('|')
				}
			}
		}
	}

	for i := 0; carts.carts.Len() > 1; i++ {
		carts.Step()
	}

	lastCart := carts.carts[0]
	fmt.Println(lastCart.x, lastCart.y)
}
