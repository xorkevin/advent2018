package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

const (
	puzzleInput = "input.txt"
)

func byteUpper(a byte) byte {
	if a-byte('a') > byte('z') {
		return a
	}
	return a - byte('a') + byte('A')
}

func react(a, b byte) bool {
	i := a
	j := b
	if i-byte('a') <= byte('z') {
		i -= byte('a') - byte('A')
	}
	if j-byte('a') <= byte('z') {
		j -= byte('a') - byte('A')
	}

	return a != b && i == j
}

type (
	Stack []byte
)

func (s *Stack) Push(a byte) {
	*s = append(*s, a)
}

func (s *Stack) Pop() byte {
	a := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return a
}

func (s *Stack) Peek() byte {
	return (*s)[len(*s)-1]
}

func (s *Stack) Len() int {
	return len(*s)
}

func reduce(line []byte) []byte {
	reduction := Stack{}
	for _, i := range line {
		if reduction.Len() == 0 {
			reduction.Push(i)
			continue
		}
		if react(reduction.Peek(), i) {
			reduction.Pop()
			continue
		}
		reduction.Push(i)
	}
	return []byte(reduction)
}

func filter(line []byte, c byte) int {
	b := bytes.Buffer{}
	for _, i := range line {
		if i != c && byteUpper(i) != c {
			b.WriteByte(i)
		}
	}
	return len(reduce(b.Bytes()))
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

	var line []byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = []byte(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(reduce(line)))

	shortestLength := 99999999
	for i := byte('A'); i <= byte('Z'); i++ {
		if size := filter(line, i); size < shortestLength {
			shortestLength = size
		}
	}

	fmt.Println(shortestLength)
}
