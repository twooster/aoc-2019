package main

import (
	"fmt"
	"math"
	"os"
	"sort"
)

type Point struct {
	x int
	y int
}

func absGCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	if a < 0 {
		return -a
	}
	return a
}

type ToShoot struct {
	dist  float64
	point Point
}

func main() {
	pts, err := readAsteroids(os.Stdin)
	if err != nil {
		fmt.Printf("Could not parse input: %v", err)
		os.Exit(1)
	}

	mostSeen := 0
	ptIdx := 0

	for i, a := range pts {
		seen := make(map[Point]bool)

		for j, b := range pts {
			if i == j {
				continue
			}
			relX := a.x - b.x
			relY := a.y - b.y
			div := absGCD(relX, relY)
			basisPt := Point{relX / div, relY / div}
			seen[basisPt] = true
		}

		seenCount := len(seen)
		if seenCount > mostSeen {
			mostSeen = seenCount
			ptIdx = i
		}
	}

	fmt.Printf("Part 1: %v\n", mostSeen)

	const pi64 = float64(math.Pi)
	const halfPi = pi64 / 2
	const twoPi = pi64 * 2
	const piAndHalf = halfPi + pi64

	rads := make(map[int][]ToShoot)

	a := pts[ptIdx]
	for j, b := range pts {
		if j == ptIdx {
			continue
		}
		relX := float64(b.x - a.x)
		relY := float64(a.y - b.y)

		dist := math.Sqrt(relX*relX + relY*relY)
		if math.Abs(dist) < 0.0000001 {
			fmt.Println("Overlapping points wat")
			os.Exit(2)
		}

		radsFromEast := math.Acos(relX / dist)
		if relY < 0 {
			radsFromEast = -radsFromEast
		}

		cwRadsFromNorth := twoPi - (radsFromEast - halfPi)
		for cwRadsFromNorth < 0 {
			cwRadsFromNorth += twoPi
		}
		for cwRadsFromNorth >= twoPi {
			cwRadsFromNorth -= twoPi
		}

		intRads := int(cwRadsFromNorth * 1000000)
		toShoot, _ := rads[intRads]
		rads[intRads] = append(toShoot, ToShoot{dist: dist, point: b})
	}

	// Sort
	radKeys := make([]int, len(rads))
	i := 0
	for rad, _ := range rads {
		radKeys[i] = rad
		toShoot := rads[rad]
		sort.SliceStable(toShoot, func(i, j int) bool {
			return toShoot[i].dist < toShoot[j].dist
		})
		i += 1
	}

	sort.Ints(radKeys)

	var twoHundredth Point
	i = 1
	for _, rad := range radKeys {
		toShoot := rads[rad]
		if len(toShoot) == 0 {
			continue
		}
		if i == 200 {
			twoHundredth = toShoot[0].point
			break
		}
		rads[rad] = toShoot[1:]
		i += 1
	}

	fmt.Printf("Part 2: %v\n", twoHundredth.x*100+twoHundredth.y)
}
