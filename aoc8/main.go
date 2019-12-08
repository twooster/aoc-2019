package main

import (
	"fmt"
	"os"
)

func part1(ints []int, size int) int {
	minZeroes := size
	result := 0

	for i := 0; i < len(ints); i += size {
		slice := ints[i : i+size]
		zeroes := 0
		ones := 0
		twos := 0
		for _, j := range slice {
			if j == 0 {
				zeroes += 1
			} else if j == 1 {
				ones += 1
			} else if j == 2 {
				twos += 1
			}
		}

		if zeroes < minZeroes {
			minZeroes = zeroes
			result = ones * twos
		}
	}
	return result
}

func overlay(ints []int, size int) []int {
	result := make([]int, size)
	copy(result, ints)

	for i := size; i < len(ints); i += size {
		for j, v := range ints[i : i+size] {
			if result[j] == 2 {
				result[j] = v
			}
		}
	}

	return result
}

func printInRows(ints []int, width int) {
	for i := 0; i < len(ints); i += width {
		for _, j := range ints[i : i+width] {
			if j == 1 {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}

func main() {
	ints, err := readAsDigits(os.Stdin)
	if err != nil {
		fmt.Printf("Could not parse input: %v", err)
		os.Exit(1)
	}

	const width = 25
	const height = 6
	const size = width * height

	if len(ints)%size != 0 {
		fmt.Printf("Bad data, not a multiple of %v\n", size)
		os.Exit(2)
	}

	fmt.Printf("Part 1: %v\n", part1(ints, width*height))
	fmt.Println("Part 2:")
	printInRows(overlay(ints, width*height), width)
}
