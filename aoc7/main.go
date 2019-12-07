package main

import (
	"fmt"
	"os"
	"sync"
)

type Executable interface {
	Execute(input <-chan int, output <-chan int) error
}

type exeuteFn func(input <-chan int, output chan<- int) error

func chainIntChans(id int, chs ...<-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for _, ch := range chs {
			for v := range ch {
				out <- v
			}
		}
		close(out)
	}()
	return out
}

func intArrayToChan(i []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range i {
			out <- v
		}
		close(out)
	}()
	return out
}

func executeCodeAcrossComputers(code []int, inputs []int) int {
	maxVal := 0

	for perm := range Permutate(inputs) {
		wg := sync.WaitGroup{}
		errors := make([]error, len(perm))

		loop := make(chan int)
		lastOut := chainIntChans(-1, intArrayToChan([]int{0}), loop)
		var firstIn <-chan int

		for i := range perm {
			p := Program{memory: NewMemory(code)}
			in := chainIntChans(i, intArrayToChan(perm[i:i+1]), lastOut)
			if firstIn == nil {
				firstIn = in
			}
			out := make(chan int)
			wg.Add(1)
			go func(i int) {
				errors[i] = p.Execute(in, out)
				close(out)
				wg.Done()
			}(i)

			lastOut = out
		}

		go func() {
			for v := range lastOut {
				loop <- v
			}
			close(loop)
		}()
		wg.Wait()

		var results []int
		wg.Add(1)
		go func() {
			for v := range firstIn {
				results = append(results, v)
			}
			wg.Done()
		}()
		wg.Wait()

		for i, err := range errors {
			if err != nil {
				fmt.Printf("Encountered error, program %v: %v\n", i, err)
			}
		}

		if results[0] > maxVal {
			maxVal = results[0]
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
