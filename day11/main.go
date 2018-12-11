package main

import (
	"fmt"
)

const (
	puzzleInput = 9810

	boardSize = 300
)

func powerLevel(x, y int) int {
	rackID := x + 10
	return (rackID*y+puzzleInput)*rackID/100%10 - 5
}

func power9(x, y int, board [][]int) int {
	power := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			power += board[y+i][x+j]
		}
	}
	return power
}

func generatePower() [][]int {
	board := make([][]int, boardSize)
	for i := range board {
		board[i] = make([]int, boardSize)
	}
	for y, i := range board {
		for x, _ := range i {
			board[y][x] = powerLevel(x+1, y+1)
		}
	}
	return board
}

func generatePartials(board [][]int) [][]int {
	partial := make([][]int, boardSize+1)
	for i := range partial {
		partial[i] = make([]int, boardSize+1)
	}
	for y := 1; y < boardSize+1; y++ {
		for x := 1; x < boardSize+1; x++ {
			partial[y][x] = board[y-1][x-1] + partial[y][x-1] + partial[y-1][x] - partial[y-1][x-1]
		}
	}
	return partial
}

func main() {
	board := generatePower()
	{
		maxPower := 0
		maxX := 0
		maxY := 0
		for y := 0; y < boardSize-2; y++ {
			for x := 0; x < boardSize-2; x++ {
				power := power9(x, y, board)
				if power > maxPower {
					maxPower = power
					maxX = x + 1
					maxY = y + 1
				}
			}
		}
		fmt.Println(maxX, maxY)
	}

	partial := generatePartials(board)
	{
		maxPower := 0
		maxX := 0
		maxY := 0
		maxSize := 0
		for i := 1; i < boardSize+1; i++ {
			for y := 0; y < boardSize-i+1; y++ {
				for x := 0; x < boardSize-i+1; x++ {
					power := partial[y+i][x+i] - partial[y][x+i] - partial[y+i][x] + partial[y][x]
					if power > maxPower {
						maxPower = power
						maxX = x + 1
						maxY = y + 1
						maxSize = i
					}
				}
			}
		}
		fmt.Println(maxPower, maxX, maxY, maxSize)
	}
}
