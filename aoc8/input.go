package main

import (
	"bufio"
	"io"
	"strconv"
)

func readAsDigits(input io.Reader) (nums []int, err error) {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		for _, c := range line {
			i, err := strconv.Atoi(string(c))
			if err != nil {
				return nil, err
			}
			nums = append(nums, i)
		}
	}
	return nums, nil
}
