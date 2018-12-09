package main

import (
	"fmt"
)

const (
	playerCount = 486
	part1Val    = 70833
	lastValue   = part1Val * 100
)

type (
	Node struct {
		Val  int
		Next *Node
		Prev *Node
	}
)

func NewNode(val int) *Node {
	m := &Node{
		Val: val,
	}
	m.Next = m
	m.Prev = m
	return m
}

func (n *Node) Insert(node *Node) {
	next := n.Next
	n.Next = node
	node.Prev = n
	node.Next = next
	next.Prev = node
}

func (n *Node) RemoveNext() *Node {
	next := n.Next
	n.Next = next.Next
	n.Next.Prev = n
	return next
}

func main() {
	score := make([]int, playerCount)
	current := NewNode(0)
	for i := 0; i < lastValue; i++ {
		player := i % playerCount
		value := i + 1
		if value%23 == 0 {
			for j := 0; j < 8; j++ {
				current = current.Prev
			}
			node := current.RemoveNext()
			current = current.Next
			score[player] += value + node.Val
		} else {
			current = current.Next
			node := NewNode(value)
			current.Insert(node)
			current = current.Next
		}
	}

	maxScore := 0
	for _, i := range score {
		if i > maxScore {
			maxScore = i
		}
	}

	fmt.Println(maxScore)
}
