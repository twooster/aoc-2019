package main

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func readAsCSVInts(input io.Reader) ([]int64, error) {
	var nums []int64
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		for _, numStr := range strings.Split(line, ",") {
			if num, err := strconv.ParseInt(numStr, 10, 64); err == nil {
				nums = append(nums, num)
			} else {
				return nil, err
			}
		}
	}
	return nums, nil
}
