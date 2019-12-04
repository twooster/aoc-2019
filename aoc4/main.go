package main

import (
	"fmt"
	"strconv"
)

func bruteForceSearch(min, max, minRun, maxRun int) int {
	count := 0
outer:
	for i := min; i <= max; i += 1 {
		s := strconv.Itoa(i)

		run := 1
		p := s[0]
		d := s[1]
		if d < p {
			continue
		} else if d == p {
			run = 2
		}
		p = d

		foundRun := 0
		for j := 2; j < len(s); j += 1 {
			d = s[j]
			if d < p {
				continue outer
			} else if d == p {
				run += 1
			} else {
				if run >= minRun && run <= maxRun {
					foundRun = 1
				}
				run = 1
			}
			p = d
		}
		if run >= minRun && run <= maxRun {
			foundRun = 1
		}
		count += foundRun
	}
	return count
}

func main() {
	const min = 152085
	const max = 670283
	fmt.Printf("Part 1: %v\n", bruteForceSearch(min, max, 2, 1000))
	fmt.Printf("Part 2: %v\n", bruteForceSearch(min, max, 2, 2))
}
