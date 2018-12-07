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

	order := []byte{}
	openList := Queue{}
	for k, _ := range tasks {
		if _, ok := isNext[k]; !ok {
			openList.Push(k)
		}
	}
	sort.Sort(openList)
	closedList := map[byte]struct{}{}

openloop:
	for len(openList) > 0 {
		top := openList.Pop()
		if _, ok := closedList[top]; ok {
			continue
		}
		task := tasks[top]
		for _, i := range task.Deps {
			if _, ok := closedList[i]; !ok {
				continue openloop
			}
		}
		closedList[top] = struct{}{}
		order = append(order, top)
		openList = append(openList, task.Next...)
		sort.Sort(openList)
	}

	fmt.Println(string(order))

	elapsedTime := 0
	openList = Queue{}
	for k, _ := range tasks {
		if _, ok := isNext[k]; !ok {
			openList.Push(k)
		}
	}
	sort.Sort(openList)
	closedList = map[byte]struct{}{}
	currentWork := []*Process{}

workloop:
	for {
		for len(currentWork) < 5 {
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
			} else {
				nextWork = append(nextWork, i)
				nextTasks[i.TaskID] = struct{}{}
			}
		}
		currentWork = nextWork

		for _, i := range doneWork {
			closedList[i] = struct{}{}
			task := tasks[i]
		nextloop:
			for _, j := range task.Next {
				if _, ok := nextTasks[j]; ok {
					continue
				}
				if _, ok := closedList[j]; ok {
					continue
				}
				next := tasks[j]
				for _, k := range next.Deps {
					if _, ok := closedList[k]; !ok {
						continue nextloop
					}
				}
				openList.Push(j)
			}
			sort.Sort(openList)
		}
	}
	fmt.Println(elapsedTime)
}
