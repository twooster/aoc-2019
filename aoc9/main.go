package main

import (
	"fmt"
	"os"
	"sync"
)

func inputChannel(a []int64) <-chan int64 {
	out := make(chan int64)
	go func() {
		for _, v := range a {
			out <- v
		}
	}()
	return out
}

func runOneComputer(code []int64, input []int64) []int64 {
	var output []int64
	var err error

	wg := sync.WaitGroup{}

	p := NewProgram(code)
	in := inputChannel(input)
	out := make(chan int64)
	wg.Add(1)
	go func() {
		err = p.Execute(in, out)
		close(out)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for v := range out {
			output = append(output, v)
		}
		wg.Done()
	}()

	wg.Wait()

	if err != nil {
		fmt.Printf("Encountered error, program: %v\n", err)
	}

	return output
}

func main() {
	code, err := readAsCSVInts(os.Stdin)
	if err != nil {
		fmt.Printf("Could not parse input: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1: %v\n", runOneComputer(code, []int64{1}))
	fmt.Printf("Part 2: %v\n", runOneComputer(code, []int64{2}))
}
