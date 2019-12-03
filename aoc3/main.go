package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Direction int

type Intersection struct {
	X    int
	Y    int
	Dist int
}

type Segment struct {
	X0        int
	Y0        int
	X1        int
	Y1        int
	StartDist int
	EndDist   int
}

func minMax(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (v *Segment) IntersectionPoint(h *Segment) *Intersection {
	if v.X0 == v.X1 && h.X0 == h.X1 {
		// both are vertical
		return nil
	} else if v.Y0 == v.Y1 && h.Y0 == h.Y1 {
		// both are horizontal
		return nil
	} else if v.Y0 == v.Y1 {
		// v is horizontal and h is vertical, thus inverted
		return h.IntersectionPoint(v)
	}

	minX, maxX := minMax(h.X0, h.X1)
	minY, maxY := minMax(v.Y0, v.Y1)

	if v.X0 >= minX && v.X0 <= maxX && h.Y0 >= minY && h.Y0 <= maxY {
		x := v.X0
		y := h.Y0
		vDist := v.StartDist + abs(y-v.Y0)
		hDist := h.StartDist + abs(x-h.X0)
		return &Intersection{
			X:    x,
			Y:    y,
			Dist: vDist + hDist,
		}
	}
	return nil
}

type Line = []Segment

func createSegment(s string, last *Segment) (*Segment, error) {
	if dist, err := strconv.Atoi(s[1:]); err != nil {
		return nil, err
	} else {
		var newX int
		var newY int
		switch s[0] {
		case 'U':
			newX = last.X1
			newY = last.Y1 + dist
		case 'R':
			newX = last.X1 + dist
			newY = last.Y1
		case 'D':
			newX = last.X1
			newY = last.Y1 - dist
		case 'L':
			newX = last.X1 - dist
			newY = last.Y1
		default:
			return nil, fmt.Errorf("Unknown direction: %v", s[0])
		}
		return &Segment{
			X0:        last.X1,
			Y0:        last.Y1,
			X1:        newX,
			Y1:        newY,
			StartDist: last.EndDist,
			EndDist:   last.EndDist + dist,
		}, nil
	}
}

func readAsPaths(input io.Reader) ([]Line, error) {
	var lines []Line
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		str := scanner.Text()
		last := &Segment{}
		line := Line{}
		for _, s := range strings.Split(str, ",") {
			var err error
			last, err = createSegment(s, last)
			if err != nil {
				return nil, err
			}
			line = append(line, *last)
		}
		lines = append(lines, line)
	}
	return lines, nil
}

func findIntersections(l1 []Segment, l2 []Segment) []Intersection {
	var intersections []Intersection
	for i := 0; i < len(l1); i += 1 {
		for j := 0; j < len(l2); j += 1 {
			xpt := l1[i].IntersectionPoint(&l2[j])
			if xpt != nil {
				intersections = append(intersections, *xpt)
			}
		}
	}
	return intersections
}

func findLowestManhattanDistanceToOrigin(c []Intersection) int {
	lowest := 10000000000
	for _, c := range c {
		var v int
		if c.X >= 0 {
			v += c.X
		} else {
			v -= c.X
		}

		if c.Y >= 0 {
			v += c.Y
		} else {
			v -= c.Y
		}

		if v > 0 && v < lowest {
			lowest = v
		}
	}
	return lowest
}

func findLowestLineDistance(c []Intersection) int {
	lowest := 10000000000
	for _, c := range c {
		if c.Dist > 0 && c.Dist < lowest {
			lowest = c.Dist
		}
	}
	return lowest
}

func main() {
	lines, err := readAsPaths(os.Stdin)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	for i := 0; i < len(lines); i += 1 {
		for j := i + 1; j < len(lines); j += 1 {
			intersections := findIntersections(lines[i], lines[j])
			fmt.Printf("Intersections: %v\n", intersections)
			lowestM := findLowestManhattanDistanceToOrigin(intersections)
			fmt.Printf("Lowest Manhattan: %v\n", lowestM)
			lowestLD := findLowestLineDistance(intersections)
			fmt.Printf("Lowest Line Distance: %v\n", lowestLD)
		}
	}
}
