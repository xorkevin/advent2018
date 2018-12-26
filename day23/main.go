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
	maxPos := Pos{}
	bots := []Bot{}
	scanner := bufio.NewScanner(file)
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
			maxPos = pos
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	part1 := 0
	for _, i := range bots {
		if maxPos.Mnhttn(i.pos) <= maxRadius {
			part1++
		}
	}
	fmt.Println(part1)
}
