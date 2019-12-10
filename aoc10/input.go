package main

import (
	"bufio"
	"io"
)

func readAsteroids(input io.Reader) (pts []Point, err error) {
	scanner := bufio.NewScanner(input)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x, c := range line {
			if c == '#' {
				pts = append(pts, Point{x, y})
			}
		}
		y += 1
	}
	return pts, nil
}
