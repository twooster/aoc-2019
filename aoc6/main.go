package main

import (
	"fmt"
	"os"
)

// a is orbited by b, c, d, ...
type OrbitedBy map[string][]string
type Orbits map[string]string

const (
	YourShip         = "YOU"
	Santa            = "SAN"
	CenterOfUniverse = "COM"
)

func main() {
	obs, err := readAsOrbits(os.Stdin)
	if err != nil {
		fmt.Printf("Could not parse input: %v", err)
		os.Exit(1)
	}

	orbits := make(Orbits)
	orbitedBy := make(OrbitedBy)
	for _, o := range obs {
		sats, _ := orbitedBy[o.Center]
		sats = append(sats, o.Satellite)
		orbitedBy[o.Center] = sats
		orbits[o.Satellite] = o.Center
	}

	c := countOrbits(orbitedBy, CenterOfUniverse, 0)
	fmt.Printf("Part 1: %v\n", c)

	transLen := findTransferPathLength(orbits, YourShip, Santa)
	fmt.Printf("Part 2: %v\n", transLen)
}

func countOrbits(o OrbitedBy, key string, dist int) int {
	sats, _ := o[key]
	sum := dist
	for _, s := range sats {
		sum += countOrbits(o, s, dist+1)
	}
	return sum
}

func getPath(o Orbits, key string) []string {
	parent, ok := o[key]
	if ok {
		return append(getPath(o, parent), key)
	}
	return []string{key}
}

func findTransferPathLength(o Orbits, key1, key2 string) int {
	path1 := getPath(o, key1)
	path2 := getPath(o, key2)

	if path1[0] != path2[0] {
		return -1
	}

	var i int
	for i = 0; i < len(path1); i += 1 {
		if path1[i] != path2[i] {
			break
		}
	}
	return len(path1) + len(path2) - i - i - 2
}
