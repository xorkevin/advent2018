package main

import (
	"bytes"
	"fmt"
)

const (
	puzzleInput  = 540391
	puzzleInput2 = "540391"
)

func split(b byte) (byte, byte) {
	return b / 10, b % 10
}

func cmp(state, target []byte) int {
	sl := len(state)
	tl := len(target)
	if sl < tl+1 {
		return -1
	}
	l := sl - tl
	if bytes.Equal(state[l-1:sl-1], target) || bytes.Equal(state[l:], target) {
		return l
	}
	return -1
}

func main() {
	count1 := 0
	count2 := 1
	state := []byte{3, 7}
	for len(state) < puzzleInput+10 {
		a, b := split(state[count1] + state[count2])
		if a > 0 {
			state = append(state, a)
		}
		state = append(state, b)
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
		a, b := split(state[count1] + state[count2])
		if a > 0 {
			state = append(state, a)
		}
		state = append(state, b)
		count1 = (count1 + int(state[count1]) + 1) % len(state)
		count2 = (count2 + int(state[count2]) + 1) % len(state)
	}
	fmt.Println(k)
}
