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
	Point struct {
		Posx int
		Posy int
		Velx int
		Vely int
	}

	PointList []*Point
)

func (p *Point) Step(t int) {
	p.Posx += p.Velx * t
	p.Posy += p.Vely * t
}

func (p PointList) Step(t int) {
	for _, i := range p {
		i.Step(t)
	}
}

func (p PointList) Size() (int, int, int, int) {
	minx := 99999999
	maxx := -99999999
	miny := 99999999
	maxy := -99999999
	for _, i := range p {
		if i.Posx < minx {
			minx = i.Posx
		}
		if i.Posx > maxx {
			maxx = i.Posx
		}
		if i.Posy < miny {
			miny = i.Posy
		}
		if i.Posy > maxy {
			maxy = i.Posy
		}
	}
	sizex := maxx - minx + 1
	sizey := maxy - miny + 1
	return sizex, sizey, minx, miny
}

func (p PointList) Print() {
	sizex, sizey, minx, miny := p.Size()

	baseline := make([]byte, 0, sizex)
	for i := 0; i < sizex; i++ {
		baseline = append(baseline, '.')
	}
	board := make([][]byte, sizey)
	for i := range board {
		board[i] = make([]byte, sizex)
		copy(board[i], baseline)
	}

	for _, i := range p {
		x := i.Posx - minx
		y := i.Posy - miny
		board[y][x] = '#'
	}

	for _, i := range board {
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

	points := PointList{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var posx, posy, velx, vely int
		fmt.Sscanf(scanner.Text(), "position=<%d, %d> velocity=<%d, %d>", &posx, &posy, &velx, &vely)
		points = append(points, &Point{
			Posx: posx,
			Posy: posy,
			Velx: velx,
			Vely: vely,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	step := 10656
	points.Step(step)
	points.Print()
	fmt.Println("Step: ", step)
	sizex, sizey, minx, miny := points.Size()
	fmt.Println(sizex, sizey, minx, miny)
}
