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

func (p *Point) Less(o *Point) bool {
	if p.y < o.y {
		return true
	}
	return p.x < o.x
}

func (p *Point) Ordered(o *Point) (*Point, *Point) {
	if p.Less(o) {
		return p, o
	}
	return o, p
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
			if _, hasSeen := seen[basisPt]; !hasSeen {
				seen[basisPt] = true
			}
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

	buckets := make(map[int][]ToShoot)

	a := pts[ptIdx]
	for j, b := range pts {
		if j == ptIdx {
			continue
		}
		relX := float64(b.x - a.x)
		relY := float64(a.y - b.y)

		dist := math.Sqrt(relX*relX + relY*relY)
		if math.Abs(dist) < 0.0000001 {
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
		bucket, _ := buckets[intRads]
		buckets[intRads] = append(bucket, ToShoot{dist: dist, point: b})
	}

	keys := make([]int, len(buckets))
	i := 0
	for key, _ := range buckets {
		keys[i] = key
		asteroids := buckets[key]
		sort.SliceStable(asteroids, func(i, j int) bool {
			return asteroids[i].dist < asteroids[j].dist
		})
		i += 1
	}

	sort.Ints(keys)

	var asteroid Point
	i = 1
	for _, k := range keys {
		asteroids := buckets[k]
		if len(asteroids) == 0 {
			continue
		}
		rec := asteroids[0]
		asteroid = rec.point
		//fmt.Printf("Pew pew: %v %.2f %.2f [%d, %d] (%d, %d)\n", i, float64(k)/1000000, rec.dist, asteroid.x, asteroid.y, asteroid.x-a.x, a.y-asteroid.y)
		if i == 200 {
			break
		}
		buckets[k] = asteroids[1:]
		i += 1
	}

	fmt.Printf("Part 2: %v\n", asteroid.x*100+asteroid.y)
}
