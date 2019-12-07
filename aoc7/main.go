package main

import (
	"fmt"
	"os"
	"sync"
)

func primedChannel(ch <-chan int, a ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range a {
			out <- v
		}
		for v := range ch {
			out <- v
		}
	}()
	return out
}

func executeCodeAcrossComputers(code []int, inputs []int) int {
	maxVal := 0

	for perm := range Permutate(inputs) {
		wg := sync.WaitGroup{}
		errors := make([]error, len(perm))

		loop := make(chan int)
		lastOut := primedChannel(loop, 0)

		for i := range perm {
			p := Program{memory: NewMemory(code)}
			in := primedChannel(lastOut, perm[i])
			out := make(chan int)
			wg.Add(1)
			go func(i int) {
				errors[i] = p.Execute(in, out)
				close(out)
				wg.Done()
			}(i)

			lastOut = out
		}

		lastResult := 0
		go func() {
			for v := range lastOut {
				lastResult = v
				loop <- v
			}
			close(loop)
		}()
		wg.Wait()

		for i, err := range errors {
			if err != nil {
				fmt.Printf("Encountered error, program %v: %v\n", i, err)
			}
		}

		if lastResult > maxVal {
			maxVal = lastResult
		}
	}

	return maxVal
}

func main() {
	code, err := readAsCSVInts(os.Stdin)
	if err != nil {
		fmt.Printf("Could not parse input: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1: %v\n", executeCodeAcrossComputers(code, []int{0, 1, 2, 3, 4}))
	fmt.Printf("Part 2: %v\n", executeCodeAcrossComputers(code, []int{5, 6, 7, 8, 9}))
}
