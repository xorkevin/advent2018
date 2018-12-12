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

func match(ind int, state []byte) byte {
	left := ind - 2
	right := ind + 3
	if val, ok := rules[string(state[left:right])]; ok {
		return byte(val)
	}
	return state[ind]
}

func nextState(state []byte) []byte {
	next := make([]byte, 0, len(state))
	next = append(next, byte('.'))
	next = append(next, byte('.'))
	for i := 2; i < len(state)-2; i++ {
		next = append(next, match(i, state))
	}
	next = append(next, byte('.'))
	next = append(next, byte('.'))
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
	state := make([]byte, 0, len(puzzleInput)*16)
	for i := 0; i < len(puzzleInput)*16; i++ {
		state = append(state, byte('.'))
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
