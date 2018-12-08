package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

type (
	Queue []int
)

func (s *Queue) Push(a int) {
	*s = append(*s, a)
}

func (s *Queue) Pop() int {
	a := (*s)[0]
	*s = (*s)[1:]
	return a
}

func (s Queue) Len() int {
	return len(s)
}

type (
	Node struct {
		Children []*Node
		Metadata []int
	}
)

func NewNode() *Node {
	return &Node{
		Children: []*Node{},
		Metadata: []int{},
	}
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

func (n *Node) AddMetadata(data int) {
	n.Metadata = append(n.Metadata, data)
}

func (n *Node) Sum() int {
	if len(n.Children) == 0 {
		s := 0
		for _, i := range n.Metadata {
			s += i
		}
		return s
	}
	s := 0
	for _, i := range n.Metadata {
		if i < len(n.Children)+1 {
			s += n.Children[i-1].Sum()
		}
	}
	return s
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

	strnums := strings.Split(line, " ")

	nums := make(Queue, 0, len(strnums))
	for _, i := range strnums {
		num, err := strconv.Atoi(i)
		if err != nil {
			log.Fatal(err)
		}
		nums.Push(num)
	}

	n, sum := processTree(&nums)
	fmt.Println(sum)
	fmt.Println(n.Sum())
}

func processTree(nums *Queue) (*Node, int) {
	numChildren := nums.Pop()
	numMetadata := nums.Pop()

	sum := 0

	node := NewNode()
	for i := 0; i < numChildren; i++ {
		child, metasum := processTree(nums)
		node.AddChild(child)
		sum += metasum
	}
	for i := 0; i < numMetadata; i++ {
		data := nums.Pop()
		node.AddMetadata(data)
		sum += data
	}
	return node, sum
}
