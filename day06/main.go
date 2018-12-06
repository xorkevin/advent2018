package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	puzzleInput = "input.txt"
	gridEnd     = 720
	gridStart   = -360
)

type (
	Point struct {
		X, Y int
	}
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

	points := []Point{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var x, y int
		fmt.Sscanf(scanner.Text(), "%d, %d", &x, &y)
		points = append(points, Point{
			X: x,
			Y: y,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	inRegion := 0
	edge := map[int]struct{}{}
	counts := make([]int, len(points))
	for i := gridStart; i < gridEnd; i++ {
		for j := gridStart; j < gridEnd; j++ {
			k, tie := findClosest(j, i, points)
			if !tie {
				counts[k]++
			}
			if isEdge(j, i) {
				edge[k] = struct{}{}
			}
			if combinedDistance(j, i, points) < 10000 {
				inRegion++
			}
		}
	}

	max := 0
	for n, i := range counts {
		if _, ok := edge[n]; !ok && i > max {
			max = i
		}
	}

	fmt.Println(max)
	fmt.Println(inRegion)
}

func combinedDistance(x, y int, points []Point) int {
	p := Point{
		X: x,
		Y: y,
	}
	dist := 0
	for _, i := range points {
		dist += distance(&p, &i)
	}
	return dist
}

func isEdge(x, y int) bool {
	return x == gridStart || y == gridStart || x == gridEnd-1 || y == gridEnd-1
}

func findClosest(x, y int, points []Point) (int, bool) {
	p := Point{
		X: x,
		Y: y,
	}
	tie := false
	ind := 0
	dist := distance(&p, &points[0])
	for n, i := range points[1:] {
		k := distance(&p, &i)
		if k < dist {
			tie = false
			dist = k
			ind = n + 1
		} else if k == dist {
			tie = true
		}
	}
	return ind, tie
}

func distance(a, b *Point) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
