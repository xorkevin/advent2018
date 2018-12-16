package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	puzzleInput  = "input.txt"
	puzzleInput2 = "input2.txt"
)

const (
	addr = 0
	addi = 1
	mulr = 2
	muli = 3
	banr = 4
	bani = 5
	borr = 6
	bori = 7
	setr = 8
	seti = 9
	gtir = 10
	gtri = 11
	gtrr = 12
	eqir = 13
	eqri = 14
	eqrr = 15
)

func ExecOp(opcode int, a, b int, state []int) int {
	ra := state[a]
	rb := state[b]
	switch opcode {
	case addr:
		return ra + rb
	case addi:
		return ra + b
	case mulr:
		return ra * rb
	case muli:
		return ra * b
	case banr:
		return ra & rb
	case bani:
		return ra & b
	case borr:
		return ra | rb
	case bori:
		return ra | b
	case setr:
		return ra
	case seti:
		return a
	case gtir:
		if a > rb {
			return 1
		}
		return 0
	case gtri:
		if ra > b {
			return 1
		}
		return 0
	case gtrr:
		if ra > rb {
			return 1
		}
		return 0
	case eqir:
		if a == rb {
			return 1
		}
		return 0
	case eqri:
		if ra == b {
			return 1
		}
		return 0
	case eqrr:
		if ra == rb {
			return 1
		}
		return 0
	default:
		panic("invalid op")
	}
}

func ExecInstr(opcode int, a, b, c int, state []int) []int {
	out := make([]int, len(state))
	copy(out, state)
	out[c] = ExecOp(opcode, a, b, state)
	return out
}

type (
	TestCase struct {
		before []int
		op     []int
		after  []int
	}
)

func ParseOp(line string) []int {
	var s1, s2, s3, s4 int
	fmt.Sscanf(line, "%d %d %d %d", &s1, &s2, &s3, &s4)
	return []int{s1, s2, s3, s4}
}

func ParseState(flag, line string) []int {
	var s1, s2, s3, s4 int
	fmt.Sscanf(line, flag+": [%d, %d, %d, %d]", &s1, &s2, &s3, &s4)
	return []int{s1, s2, s3, s4}
}

func StateEqual(s1, s2 []int) bool {
	for i := 0; i < 4; i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

type (
	TranslationTable [][]int
)

func NewTranslationTable() TranslationTable {
	k := make(TranslationTable, 16)
	for i := range k {
		k[i] = make([]int, 16)
	}
	return k
}

func (tr TranslationTable) ProvideEx(testCase *TestCase) int {
	count := 0
	for i := 0; i < 16; i++ {
		op := testCase.op
		after := ExecInstr(i, op[1], op[2], op[3], testCase.before)
		if !StateEqual(after, testCase.after) {
			tr[op[0]][i] = 1
			count++
		}
	}
	return count
}

func (tr TranslationTable) GenTranslation() []int {
	k := make(map[int]int, len(tr))
	for len(k) < len(tr) {
		for op, opts := range tr {
			if _, ok := k[op]; ok {
				continue
			}
			count := 0
			last := 0
			for n, i := range opts {
				if i != 1 {
					count++
					last = n
				}
			}
			if count == 1 {
				for i := range tr {
					tr[i][last] = 1
				}
				k[op] = last
				break
			}
		}
	}

	f := make([]int, len(k))
	for i := range f {
		f[i] = k[i]
	}
	return f
}

func part1() []int {
	file, err := os.Open(puzzleInput)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	table := NewTranslationTable()
	part1Count := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		before := ParseState("Before", scanner.Text())
		scanner.Scan()
		op := ParseOp(scanner.Text())
		scanner.Scan()
		after := ParseState("After", scanner.Text())
		scanner.Scan()

		if count := table.ProvideEx(&TestCase{
			before: before,
			op:     op,
			after:  after,
		}); 16-count > 2 {
			part1Count++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(part1Count)

	return table.GenTranslation()
}

func part2(translator []int) {
	file, err := os.Open(puzzleInput2)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	state := []int{0, 0, 0, 0}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		op := ParseOp(scanner.Text())
		state = ExecInstr(translator[op[0]], op[1], op[2], op[3], state)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(state)
}

func main() {
	part2(part1())
}
