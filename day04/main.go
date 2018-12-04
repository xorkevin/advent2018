package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

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

	lines := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	part1(lines)
	part2(lines)
}

func part1(lines []string) {
	guardHours := map[int]int{}
	currentGuard := 0
	asleepTime := 0

	for _, line := range lines {
		var year, month, day, hour, minute int
		fmt.Sscanf(line, "[%d-%d-%d %d:%d]", &year, &month, &day, &hour, &minute)
		action := strings.Fields(strings.Split(line, "]")[1])
		if action[0] == "Guard" {
			var guardnum int
			fmt.Sscanf(action[1], "#%d", &guardnum)
			currentGuard = guardnum
		} else if action[0] == "falls" {
			asleepTime = minute
		} else if action[0] == "wakes" {
			if _, ok := guardHours[currentGuard]; ok {
				guardHours[currentGuard] += minute - asleepTime
			} else {
				guardHours[currentGuard] = minute - asleepTime
			}
		}
	}

	guardid := 0
	hoursAsleep := 0
	for k, v := range guardHours {
		if v > hoursAsleep {
			guardid = k
			hoursAsleep = v
		}
	}

	guardMinutes := make([]int, 60)

	for _, line := range lines {
		var year, month, day, hour, minute int
		fmt.Sscanf(line, "[%d-%d-%d %d:%d]", &year, &month, &day, &hour, &minute)
		action := strings.Fields(strings.Split(line, "]")[1])
		if action[0] == "Guard" {
			var guardnum int
			fmt.Sscanf(action[1], "#%d", &guardnum)
			currentGuard = guardnum
		} else if action[0] == "falls" {
			asleepTime = minute
		} else if action[0] == "wakes" {
			if currentGuard == guardid {
				for i := asleepTime; i < minute; i++ {
					guardMinutes[i]++
				}
			}
		}
	}

	minute := 0
	minuteAmount := 0
	for n, i := range guardMinutes {
		if i > minuteAmount {
			minute = n
			minuteAmount = i
		}
	}

	fmt.Println(minute * guardid)
}

func part2(lines []string) {
	minutes := map[int][]int{}
	currentGuard := 0
	asleepTime := 0

	for _, line := range lines {
		var year, month, day, hour, minute int
		fmt.Sscanf(line, "[%d-%d-%d %d:%d]", &year, &month, &day, &hour, &minute)
		action := strings.Fields(strings.Split(line, "]")[1])
		if action[0] == "Guard" {
			var guardnum int
			fmt.Sscanf(action[1], "#%d", &guardnum)
			currentGuard = guardnum
		} else if action[0] == "falls" {
			asleepTime = minute
		} else if action[0] == "wakes" {
			if _, ok := minutes[currentGuard]; !ok {
				minutes[currentGuard] = make([]int, 60)
			}
			arr := minutes[currentGuard]
			for i := asleepTime; i < minute; i++ {
				arr[i]++
			}
		}
	}

	maxGuard := 0
	maxMinute := 0
	maxAmount := 0

	for k, v := range minutes {
		for n, i := range v {
			if i > maxAmount {
				maxGuard = k
				maxMinute = n
				maxAmount = i
			}
		}
	}

	fmt.Println(maxGuard * maxMinute)
}
