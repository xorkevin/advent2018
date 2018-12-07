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
	numWorkers  = 5
)

type (
	Task struct {
		Id   byte
		Deps []byte
		Next []byte
	}

	Process struct {
		TaskID byte
		Cost   int
	}
)

func NewTask(id byte) *Task {
	return &Task{
		Id:   id,
		Deps: []byte{},
		Next: []byte{},
	}
}

func (t *Task) AddDep(dep byte) {
	t.Deps = append(t.Deps, dep)
}

func (t *Task) AddNext(next byte) {
	t.Next = append(t.Next, next)
}

func (t *Task) CanBegin(finished map[byte]struct{}) bool {
	for _, i := range t.Deps {
		if _, ok := finished[i]; !ok {
			return false
		}
	}
	return true
}

func NewProcess(id byte) *Process {
	return &Process{
		TaskID: id,
		Cost:   int(id) - int('A') + 61,
	}
}

type (
	Queue []byte
)

func (s *Queue) Push(a byte) {
	*s = append(*s, a)
}

func (s *Queue) Pop() byte {
	a := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return a
}

func (s *Queue) Peek() byte {
	return (*s)[len(*s)-1]
}

func (s Queue) Less(i, j int) bool {
	return s[i] > s[j]
}

func (s Queue) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Queue) Len() int {
	return len(s)
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

	isNext := map[byte]struct{}{}
	tasks := map[byte]*Task{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		link := strings.Fields(scanner.Text())
		dep := byte(link[1][0])
		next := byte(link[7][0])
		if task, ok := tasks[next]; ok {
			task.AddDep(dep)
		} else {
			k := NewTask(next)
			k.AddDep(dep)
			tasks[next] = k
		}
		if task, ok := tasks[dep]; ok {
			task.AddNext(next)
		} else {
			k := NewTask(dep)
			k.AddNext(next)
			tasks[dep] = k
		}
		isNext[next] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	openList := Queue{}
	for k, _ := range tasks {
		if _, ok := isNext[k]; !ok {
			openList.Push(k)
		}
	}
	openList2 := make(Queue, len(openList))
	copy(openList2, openList)
	part1(openList, tasks)
	part2(openList2, tasks)
}

func part1(openList Queue, tasks map[byte]*Task) {
	sort.Sort(openList)

	order := []byte{}

	closedList := map[byte]struct{}{}
	for len(openList) > 0 {
		top := openList.Pop()
		order = append(order, top)
		closedList[top] = struct{}{}
		task := tasks[top]
		for _, i := range task.Next {
			if _, ok := closedList[i]; ok {
				continue
			}
			next := tasks[i]
			if !next.CanBegin(closedList) {
				continue
			}
			openList.Push(i)
		}
		sort.Sort(openList)
	}

	fmt.Println(string(order))
}

func part2(openList Queue, tasks map[byte]*Task) {
	sort.Sort(openList)

	elapsedTime := 0

	closedList := map[byte]struct{}{}
	currentWork := []*Process{}

workloop:
	for {
		for len(currentWork) < numWorkers {
			if openList.Len() == 0 {
				if len(currentWork) == 0 {
					break workloop
				}
				break
			}
			top := openList.Pop()
			currentWork = append(currentWork, NewProcess(top))
		}
		minTime := 99999999
		for _, i := range currentWork {
			if i.Cost < minTime {
				minTime = i.Cost
			}
		}
		elapsedTime += minTime

		nextWork := []*Process{}
		nextTasks := map[byte]struct{}{}
		doneWork := []byte{}
		for _, i := range currentWork {
			i.Cost -= minTime
			if i.Cost == 0 {
				doneWork = append(doneWork, i.TaskID)
				closedList[i.TaskID] = struct{}{}
			} else {
				nextWork = append(nextWork, i)
				nextTasks[i.TaskID] = struct{}{}
			}
		}
		currentWork = nextWork

		for _, i := range doneWork {
			task := tasks[i]
			for _, j := range task.Next {
				if _, ok := nextTasks[j]; ok {
					continue
				}
				if _, ok := closedList[j]; ok {
					continue
				}
				next := tasks[j]
				if !next.CanBegin(closedList) {
					continue
				}
				openList.Push(j)
			}
			sort.Sort(openList)
		}
	}
	fmt.Println(elapsedTime)
}
