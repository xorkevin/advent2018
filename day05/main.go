package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"unicode"
)

const (
	puzzleInput = "input.txt"
)

func transform(line []byte) ([]byte, bool) {
	b := bytes.Buffer{}
	changed := false
	for i := 1; i < len(line); i++ {
		if line[i-1] == line[i] || byte(unicode.ToUpper(rune(line[i-1]))) != byte(unicode.ToUpper(rune(line[i]))) {
			b.WriteByte(line[i-1])
		} else {
			changed = true
			i++
		}
	}
	l := len(line)
	if line[l-2] == line[l-1] || byte(unicode.ToUpper(rune(line[l-2]))) != byte(unicode.ToUpper(rune(line[l-1]))) {
		b.WriteByte(line[l-1])
	}
	return b.Bytes(), changed
}

func filter(line []byte, c byte) []byte {
	b := bytes.Buffer{}
	for _, i := range line {
		if i != c && byte(unicode.ToUpper(rune(i))) != c {
			b.WriteByte(i)
		}
	}
	return b.Bytes()
}

func iterate(line []byte) int {
	l := line
	changed := true
	for changed {
		l, changed = transform(l)
	}
	return len(l)
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

	line := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	l := []byte(line)
	changed := true
	for changed {
		l, changed = transform(l)
	}
	fmt.Println(len(l))

	shortestLength := 99999999
	for i := byte('A'); i <= byte('Z'); i++ {
		if size := iterate(filter([]byte(line), i)); size < shortestLength {
			shortestLength = size
		}
	}

	fmt.Println(shortestLength)
}
