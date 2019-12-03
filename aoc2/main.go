package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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

const Add = 1
const Mul = 2
const Stop = 99

func ensureBounds(p []int, addr int) error {
	if addr < 0 {
		return fmt.Errorf("Address < 0: %v", addr)
	} else if addr >= len(p) {
		return fmt.Errorf("Address > max len: %v > %v", addr, len(p))
	}
	return nil
}

func getValue(p []int, addr int) (int, error) {
	if err := ensureBounds(p, addr); err != nil {
		return 0, err
	}
	return p[addr], nil
}

func setValue(p []int, addr int, val int) error {
	if err := ensureBounds(p, addr); err != nil {
		return err
	}
	p[addr] = val
	return nil
}

func getIndirectValue(p []int, indirectAddr int) (int, error) {
	addr, err := getValue(p, indirectAddr)
	if err != nil {
		return 0, err
	}
	val, err := getValue(p, addr)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func setIndirectValue(p []int, indirectAddr int, val int) error {
	addr, err := getValue(p, indirectAddr)
	if err != nil {
		return err
	}
	return setValue(p, addr, val)
}

func runProgram(p []int) error {
	ip := 0
loop:
	for {
		inst, err := getValue(p, ip)
		if err != nil {
			return err
		}

		switch inst {
		case Add:
			a, err := getIndirectValue(p, ip+1)
			if err != nil {
				return err
			}

			b, err := getIndirectValue(p, ip+2)
			if err != nil {
				return err
			}

			if err := setIndirectValue(p, ip+3, a+b); err != nil {
				return err
			}

			ip += 4
		case Mul:
			a, err := getIndirectValue(p, ip+1)
			if err != nil {
				return err
			}

			b, err := getIndirectValue(p, ip+2)
			if err != nil {
				return err
			}

			if err := setIndirectValue(p, ip+3, a*b); err != nil {
				return err
			}

			ip += 4
		case Stop:
			break loop
		default:
			return fmt.Errorf("Unknown operation: %v at %v", inst, ip)
		}
	}
	return nil
}

func runProgramWithValues(p []int, pos1, pos2 int) (int, error) {
	if len(p) < 3 {
		return 0, fmt.Errorf("Program too short: %v", len(p))
	}
	c := make([]int, len(p))
	copy(c, p)
	c[1] = pos1
	c[2] = pos2
	if err := runProgram(c); err != nil {
		return 0, err
	}
	return c[0], nil
}

func main() {
	const part2Search = 19690720

	p, err := readAsCSVInts(os.Stdin)
	if err != nil {
		fmt.Printf("Could not parse input: %v", err)
		os.Exit(1)
	}

	if result, err := runProgramWithValues(p, 12, 2); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(2)
	} else {
		fmt.Printf("Part 1: %v\n", result)
	}

	var noun int
	var verb int
	found := false
outer:
	for noun = 0; noun <= 99; noun += 1 {
		for verb = 0; verb <= 99; verb += 1 {
			if result, err := runProgramWithValues(p, noun, verb); err != nil {
				fmt.Printf("Error running program with values %v, %v: %v\n", noun, verb, err)
			} else if result == part2Search {
				found = true
				break outer
			}
		}
	}

	if found {
		fmt.Printf("Part 2: 100 * %v + %v = %v\n", noun, verb, 100*noun+verb)
	} else {
		fmt.Printf("Unable to find noun and verb for %v!\n", part2Search)
	}
}
