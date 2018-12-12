package main

import (
	"fmt"
)

const (
	puzzleInput = "#.#..#..###.###.#..###.#####...########.#...#####...##.#....#.####.#.#..#..#.#..###...#..#.#....##."
)

var (
	rules = map[string]rune{
		"#.###": '.',
		"###.#": '#',
		".##..": '.',
		"..###": '.',
		"..##.": '.',
		"##...": '#',
		"###..": '#',
		".#...": '#',
		"##..#": '#',
		"#....": '.',
		".#.#.": '.',
		"####.": '.',
		"#.#..": '.',
		"#.#.#": '.',
		"#..##": '#',
		".####": '#',
		"...##": '.',
		"#..#.": '#',
		".#.##": '#',
		"..#.#": '#',
		"##.#.": '#',
		"#.##.": '#',
		"#####": '.',
		"..#..": '#',
		"....#": '.',
		"##.##": '.',
		".###.": '#',
		".....": '.',
		"...#.": '#',
		".##.#": '.',
		"#...#": '.',
		".#..#": '#',
	}
)

func match(ind int, state []byte) (byte, bool) {
	left := ind - 2
	right := ind + 3
	if val, ok := rules[string(state[left:right])]; ok {
		return byte(val), true
	}
	return 0, false
}

func nextState(state []byte) []byte {
	next := make([]byte, len(state))
	next[0] = byte('.')
	next[1] = byte('.')
	next[len(next)-2] = byte('.')
	next[len(next)-1] = byte('.')
	for i := 2; i < len(next)-2; i++ {
		if val, ok := match(i, state); ok {
			next[i] = val
		} else {
			next[i] = state[i]
		}
	}
	return next
}

func score(zeroInd int, state []byte) int {
	count := 0
	for n, i := range state {
		if i == byte('#') {
			count += n - zeroInd
		}
	}
	return count
}

const (
	iterations int64 = 50000000000
)

func main() {
	state := make([]byte, len(puzzleInput)*16)
	for i := range state {
		state[i] = byte('.')
	}
	zeroInd := len(puzzleInput)
	copy(state[zeroInd:], []byte(puzzleInput))

	prevDelta := 0
	prevScore := score(zeroInd, state)
	for i := 0; i < 200; i++ {
		state = nextState(state)
		score := score(zeroInd, state)
		delta := score - prevScore
		if i == 19 {
			fmt.Println(score)
		}
		if delta == prevDelta {
			fmt.Println(int64(score) + (iterations-int64(i)-1)*int64(delta))
			return
		}
		prevScore = score
		prevDelta = delta
	}
}
