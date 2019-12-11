package main

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func readAsCommaSeparatedWords(input io.Reader) ([]Word, error) {
	var nums []Word
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		for _, numStr := range strings.Split(line, ",") {
			if num, err := strconv.ParseInt(numStr, 10, WordWidth); err == nil {
				nums = append(nums, Word(num))
			} else {
				return nil, err
			}
		}
	}
	return nums, nil
}
