package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type OrbitRecord struct {
	Center    string
	Satellite string
}

func readAsOrbits(input io.Reader) ([]OrbitRecord, error) {
	var orbits []OrbitRecord

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		pieces := strings.Split(line, ")")
		if len(pieces) != 2 {
			return nil, fmt.Errorf("Malformed orbit line: %v", line)
		}

		orbits = append(orbits, OrbitRecord{
			Center:    pieces[0],
			Satellite: pieces[1],
		})
	}
	return orbits, nil
}
