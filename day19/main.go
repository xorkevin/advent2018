package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

const (
	addr = iota
	addi
	mulr
	muli
	banr
	bani
	borr
	bori
	setr
	seti
	gtir
	gtri
	gtrr
	eqir
	eqri
	eqrr
)

func ExecOp(opcode int, a, b int, state []int) int {
	switch opcode {
	case addr:
		return state[a] + state[b]
	case addi:
		return state[a] + b
	case mulr:
		return state[a] * state[b]
	case muli:
		return state[a] * b
	case banr:
		return state[a] & state[b]
	case bani:
		return state[a] & b
	case borr:
		return state[a] | state[b]
	case bori:
		return state[a] | b
	case setr:
		return state[a]
	case seti:
		return a
	case gtir:
		if a > state[b] {
			return 1
		}
		return 0
	case gtri:
		if state[a] > b {
			return 1
		}
		return 0
	case gtrr:
		if state[a] > state[b] {
			return 1
		}
		return 0
	case eqir:
		if a == state[b] {
			return 1
		}
		return 0
	case eqri:
		if state[a] == b {
			return 1
		}
		return 0
	case eqrr:
		if state[a] == state[b] {
			return 1
		}
		return 0
	default:
		panic("invalid op")
	}
}

func ExecInstr(instr []int, state []int) []int {
	opcode := instr[0]
	a := instr[1]
	b := instr[2]
	c := instr[3]
	out := make([]int, len(state))
	copy(out, state)
	out[c] = ExecOp(opcode, a, b, state)
	return out
}

func parseInstrToCode(instr string) int {
	switch instr {
	case "addr":
		return addr
	case "addi":
		return addi
	case "mulr":
		return mulr
	case "muli":
		return muli
	case "banr":
		return banr
	case "bani":
		return bani
	case "borr":
		return borr
	case "bori":
		return bori
	case "setr":
		return setr
	case "seti":
		return seti
	case "gtir":
		return gtir
	case "gtri":
		return gtri
	case "gtrr":
		return gtrr
	case "eqir":
		return eqir
	case "eqri":
		return eqri
	case "eqrr":
		return eqrr
	default:
		panic("Invalid instr")
	}
}

func parseLine(line string) []int {
	k := strings.SplitN(line, " ", 2)
	var a, b, c int
	fmt.Sscanf(k[1], "%d %d %d", &a, &b, &c)
	return []int{parseInstrToCode(k[0]), a, b, c}
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

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		log.Fatal("Failed to read file")
	}
	var ipreg int
	fmt.Sscanf(scanner.Text(), "#ip %d", &ipreg)

	instructions := [][]int{}
	for scanner.Scan() {
		instructions = append(instructions, parseLine(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	debugger := bufio.NewScanner(os.Stdin)
	debug := false
	mode1 := false

	last := 0
	debugr := 1
	mode2 := false

	state := []int{0, 0, 0, 0, 0, 0}
	for state[ipreg] > -1 && state[ipreg] < len(instructions) {
		state = ExecInstr(instructions[state[ipreg]], state)
		state[ipreg]++
		if mode2 && last != state[debugr] {
			last = state[debugr]
			fmt.Println(state)
		}
		//if mode1 && state[4] == 2 && state[2] == 1 {
		//	debug = true
		//}
		if mode1 && state[0] == 502 {
			debug = true
		}
		if mode1 && debug {
			fmt.Print(state)
			debugger.Scan()
		}
	}
	fmt.Println(state[0])

	// shell: factor: 10551398
	// prime factors: 2, 11, 13, 79, 467
	k := 0
	for _, i := range permute([]int{2, 11, 13, 79, 467}) {
		k += i
	}
	fmt.Println(k)
}

func permute(primefactors []int) []int {
	factors := map[int]struct{}{}
	l := len(primefactors)
	for num := 0; num < (1 << uint(l)); num++ {
		k := 1
		for i := range primefactors {
			if num&(1<<uint(i)) != 0 {
				k *= primefactors[i]
			}
		}
		factors[k] = struct{}{}
	}
	f := make([]int, 0, len(factors))
	for k, _ := range factors {
		f = append(f, k)
	}
	sort.Ints(f)
	return f
}
