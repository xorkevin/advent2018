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

func transform(from, to *bytes.Buffer) bool {
	changed := false
	kprev, err := from.ReadByte()
	kcur, err := from.ReadByte()
	var next byte
	for {
		if react(kprev, kcur) {
			changed = true
			next, err = from.ReadByte()
			if err != nil {
				return true
			}
			kprev = kcur
			kcur = next
		} else {
			to.WriteByte(kprev)
		}
		next, err = from.ReadByte()
		if err != nil {
			if !react(kprev, kcur) {
				to.WriteByte(kcur)
			}
			return changed
		}
		kprev = kcur
		kcur = next
	}
}

func iterate(from *bytes.Buffer) int {
	to := &bytes.Buffer{}
	tolen := 0
	changed := true
	for changed {
		changed = transform(from, to)
		tolen = to.Len()
		t := from
		from = to
		to = t
		to.Reset()
	}
	return tolen
}

func filter(line []byte, c byte) int {
	b := bytes.Buffer{}
	for _, i := range line {
		if i != c && byteUpper(i) != c {
			b.WriteByte(i)
		}
	}
	return iterate(&b)
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

	b := bytes.Buffer{}
	b.Write(line)
	fmt.Println(iterate(&b))

	shortestLength := 99999999
	for i := byte('A'); i <= byte('Z'); i++ {
		if size := filter(line, i); size < shortestLength {
			shortestLength = size
		}
	}

	fmt.Println(shortestLength)
}
