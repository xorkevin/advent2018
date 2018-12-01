package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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

	numlist := []int{}
	incr := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		numlist = append(numlist, num)
		incr += num
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(incr)

	visited := map[int]struct{}{}
	visited[0] = struct{}{}

	current := 0
	for {
		for _, i := range numlist {
			current += i
			if _, ok := visited[current]; ok {
				fmt.Println(current)
				return
			}
			visited[current] = struct{}{}
		}
	}
}
