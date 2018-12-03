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

	claims := make([][]int, 1000)
	for i := range claims {
		claims[i] = make([]int, 1000)
	}

	claims2 := make([][]int, 1000)
	for i := range claims2 {
		claims2[i] = make([]int, 1000)
	}

	bannedClaims := make([]bool, 1253)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		claimID := 0
		col := 0
		row := 0
		width := 0
		height := 0
		fmt.Sscanf(line, "#%d @ %d,%d: %dx%d", &claimID, &col, &row, &width, &height)
		for i := 0; i < height; i++ {
			for j := 0; j < width; j++ {
				claims[row+i][col+j]++
				c2 := claims2[row+i][col+j]
				if c2 > 0 {
					bannedClaims[claimID-1] = true
					bannedClaims[c2-1] = true
				}
				claims2[row+i][col+j] = claimID
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	count := 0
	for _, i := range claims {
		for _, j := range i {
			if j > 1 {
				count++
			}
		}
	}

	fmt.Println(count)

	for n, i := range bannedClaims {
		if !i {
			fmt.Println(n + 1)
			return
		}
	}
}
