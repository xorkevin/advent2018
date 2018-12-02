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

	lines := []string{}
	twoCount := 0
	threeCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		twice, thrice := getTwiceThrice(line)
		if twice {
			twoCount += 1
		}
		if thrice {
			threeCount += 1
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(twoCount * threeCount)

	for _, line1 := range lines {
		for _, line2 := range lines {
			if isSingleDiff(line1, line2) {
				fmt.Println(removeSingleDiff(line1, line2))
				return
			}
		}
	}
}

func getTwiceThrice(line string) (bool, bool) {
	seen := map[string]int{}
	chars := strings.Split(line, "")
	for _, char := range chars {
		if _, ok := seen[char]; ok {
			seen[char] += 1
		} else {
			seen[char] = 1
		}
	}
	isTwice := false
	isThrice := false
	for _, v := range seen {
		if v == 2 {
			isTwice = true
		} else if v == 3 {
			isThrice = true
		}
	}
	return isTwice, isThrice
}

func isSingleDiff(line1, line2 string) bool {
	if len(line1) != len(line2) {
		return false
	}
	diff := false
	for i := range line1 {
		if line1[i] != line2[i] {
			if !diff {
				diff = true
			} else {
				return false
			}
		}
	}
	return diff
}

func removeSingleDiff(line1, line2 string) string {
	if len(line1) != len(line2) {
		return ""
	}
	for i := range line1 {
		if line1[i] != line2[i] {
			return line1[:i] + line1[i+1:]
		}
	}
	return ""
}
