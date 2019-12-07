package main

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func readAsCSVInts(input io.Reader) ([]int, error) {
	var nums []int
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		for _, numStr := range strings.Split(line, ",") {
			if num, err := strconv.Atoi(numStr); err == nil {
				nums = append(nums, num)
			} else {
				return nil, err
			}
		}
	}
	return nums, nil
}
