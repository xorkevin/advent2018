package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

type (
	Event struct {
		Time  int
		Delta int
	}
	EventList []Event

	Timesheet map[int]map[int]int
)

func (e EventList) Len() int {
	return len(e)
}

func (e EventList) Less(i, j int) bool {
	return e[i].Time < e[j].Time
}

func (e EventList) Swap(i, j int) {
	t := e[i]
	e[i] = e[j]
	e[j] = t
}

func (t Timesheet) Get(guardid int) map[int]int {
	if v, ok := t[guardid]; ok {
		return v
	}
	v := map[int]int{}
	t[guardid] = v
	return v
}

func (t Timesheet) Add(guardid int, time int, delta int) {
	v := t.Get(guardid)
	v[time] += delta
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

	timesheet := Timesheet{}
	currentGuard := 0
	asleepTime := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "]")
		minute, err := strconv.Atoi(strings.Split(line[0], ":")[1])
		if err != nil {
			log.Fatal(err)
		}
		action := strings.Fields(line[1])
		switch action[0] {
		case "Guard":
			guardid, err := strconv.Atoi(action[1][1:])
			if err != nil {
				log.Fatal(err)
			}
			currentGuard = guardid
		case "falls":
			asleepTime = minute
		case "wakes":
			timesheet.Add(currentGuard, asleepTime, 1)
			timesheet.Add(currentGuard, minute, -1)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	guardtimes := map[int]EventList{}
	for guardid, v := range timesheet {
		events := make(EventList, 0, len(v))
		for time, delta := range v {
			events = append(events, Event{
				Time:  time,
				Delta: delta,
			})
		}
		sort.Sort(events)
		guardtimes[guardid] = events
	}

	// part 1
	totalMaxGuardid := 0
	totalMaxTime := 0
	totalMaxFreqMin := 0

	for guardid, events := range guardtimes {
		guardMaxFreq := 0
		guardMaxFreqMin := 0

		prevTime := 0
		prevCounter := 0
		cumulative := 0
		counter := 0
		for _, i := range events {
			prevCounter = counter
			counter += i.Delta
			cumulative += prevCounter * (i.Time - prevTime)
			prevTime = i.Time

			if counter > guardMaxFreq {
				guardMaxFreq = counter
				guardMaxFreqMin = i.Time
			}
		}

		if cumulative > totalMaxTime {
			totalMaxTime = cumulative
			totalMaxGuardid = guardid
			totalMaxFreqMin = guardMaxFreqMin
		}
	}

	fmt.Println("Part1: ", totalMaxGuardid*totalMaxFreqMin)

	// part 2
	maxGuardid := 0
	maxTime := 0
	maxFreqMin := 0

	for guardid, events := range guardtimes {
		counter := 0
		for _, i := range events {
			counter += i.Delta

			if counter > maxTime {
				maxTime = counter
				maxGuardid = guardid
				maxFreqMin = i.Time
			}
		}
	}

	fmt.Println("Part2: ", maxGuardid*maxFreqMin)
}
