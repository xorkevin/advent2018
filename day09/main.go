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
	playerScore := make([]int, playerCount)
	current := NewNode(0)
	for i := 0; i < lastValue; i++ {
		currentPlayer := i % playerCount
		currentValue := i + 1
		if currentValue%23 == 0 {
			playerScore[currentPlayer] += currentValue
			for j := 0; j < 8; j++ {
				current = current.Prev
			}
			node := current.RemoveNext()
			current = current.Next
			playerScore[currentPlayer] += node.Val
		} else {
			current = current.Next
			node := NewNode(currentValue)
			current.Insert(node)
			current = node
		}
	}

	maxScore := 0
	for _, i := range playerScore {
		if i > maxScore {
			maxScore = i
		}
	}

	fmt.Println(maxScore)
}
