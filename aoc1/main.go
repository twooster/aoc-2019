package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func readStdinIntoIntArray() ([]int, error) {
	var nums []int
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if num, err := strconv.Atoi(line); err == nil {
			nums = append(nums, num)
		} else {
			return nil, err
		}
	}
	return nums, nil
}

func calcFuel(w int) int {
	floor_div_3 := w / 3
	sub_2 := floor_div_3 - 2
	return sub_2
}

func part1(nums []int) int {
	total := 0
	for _, num := range nums {
		total += calcFuel(num)
	}
	return total
}

func calcTotalFuel(w int) int {
	fuel := calcFuel(w)
	if fuel <= 0 {
		return 0
	}
	return fuel + calcTotalFuel(fuel)
}

func part2(nums []int) int {
	total := 0
	for _, num := range nums {
		total += calcTotalFuel(num)
	}
	return total
}

func main() {
	nums, err := readStdinIntoIntArray()
	if err != nil {
		panic(fmt.Sprintf("Could not parse input: %v", err))
	}

	fmt.Printf("Part 1: %v\n", part1(nums))
	fmt.Printf("Part 2: %v\n", part2(nums))
}
