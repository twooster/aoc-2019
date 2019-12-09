package main

import (
	"fmt"
	"os"
	"sync"
)

func inputChannel(a []Word) <-chan Word {
	out := make(chan Word)
	go func() {
		for _, v := range a {
			out <- v
		}
	}()
	return out
}

func runOneComputer(code []Word, input []Word) []Word {
	var output []Word
	var err error

	wg := sync.WaitGroup{}

	p := NewProgram(code)
	in := inputChannel(input)
	out := make(chan Word)
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
	code, err := readAsCommanSeparatedWords(os.Stdin)
	if err != nil {
		fmt.Printf("Could not parse input: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1: %v\n", runOneComputer(code, []Word{1}))
	fmt.Printf("Part 2: %v\n", runOneComputer(code, []Word{2}))
}
