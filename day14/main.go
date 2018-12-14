package main

import (
	"fmt"
)

const (
	puzzleInput  = 540391
	puzzleInput2 = "540391"
)

func split(b byte) []byte {
	if b > 9 {
		return []byte{b / 10, b % 10}
	}
	return []byte{b % 10}
}

func equalBytes(a, b []byte) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func cmp(state, target []byte) int {
	if len(state) < len(target)+1 {
		return -1
	}
	l := len(state) - len(target)
	if equalBytes(state[l-1:len(state)-1], target) {
		return l
	}
	if equalBytes(state[l:], target) {
		return l
	}
	return -1
}

func main() {
	count1 := 0
	count2 := 1
	state := []byte{3, 7}
	for len(state) < puzzleInput+10 {
		state = append(state, split(state[count1]+state[count2])...)
		count1 = (count1 + int(state[count1]) + 1) % len(state)
		count2 = (count2 + int(state[count2]) + 1) % len(state)
	}
	out := make([]byte, 10)
	copy(out, state[puzzleInput:puzzleInput+10])
	for i := range out {
		out[i] += byte('0')
	}
	fmt.Println(string(out))

	target := []byte(puzzleInput2)
	for n, i := range target {
		target[n] = i - byte('0')
	}
	state = []byte{3, 7}
	count1 = 0
	count2 = 1
	k := -1
	for ; k < 0; k = cmp(state, target) {
		state = append(state, split(state[count1]+state[count2])...)
		count1 = (count1 + int(state[count1]) + 1) % len(state)
		count2 = (count2 + int(state[count2]) + 1) % len(state)
	}
	fmt.Println(k)
}
