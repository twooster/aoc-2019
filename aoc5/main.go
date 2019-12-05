package main

import (
	"errors"
	"fmt"
	"os"
)

type Operation int
type Instruction int

const (
	OpAdd         Operation = 1
	OpMul                   = 2
	OpInput                 = 3
	OpOutput                = 4
	OpJumpIfTrue            = 5
	OpJumpIfFalse           = 6
	OpLessThan              = 7
	OpEquals                = 8
	OpHalt                  = 99
)

func (i Instruction) MemMode(pos int) MemMode {
	switch pos {
	case 1:
		return MemMode((i / 100) % 10)
	case 2:
		return MemMode((i / 1000) % 10)
	case 3:
		return MemMode((i / 10000) % 10)
	}
	return -1
}

func (i Instruction) Operation() Operation {
	return Operation(i % 100)
}

type binaryFn func(a, b int) int

func runProgram(m Memory, input <-chan int, output chan<- int) error {
	var inst Instruction
	ip := 0

	binaryWriteOp := func(fn binaryFn) error {
		a, err := m.Read(inst.MemMode(1), ip+1)
		if err != nil {
			return err
		}

		b, err := m.Read(inst.MemMode(2), ip+2)
		if err != nil {
			return err
		}

		if err := m.Write(inst.MemMode(3), ip+3, fn(a, b)); err != nil {
			return err
		}

		ip += 4
		return nil
	}

	for {
		if v, err := m.ReadImmediate(ip); err != nil {
			return err
		} else {
			inst = Instruction(v)
		}

		switch inst.Operation() {
		case OpAdd:
			if err := binaryWriteOp(func(a, b int) int { return a + b }); err != nil {
				return err
			}
		case OpMul:
			if err := binaryWriteOp(func(a, b int) int { return a * b }); err != nil {
				return err
			}
		case OpInput:
			v, ok := <-input
			if !ok {
				return errors.New("Input EOF")
			}
			if err := m.Write(inst.MemMode(1), ip+1, v); err != nil {
				return err
			}
			ip += 2
		case OpOutput:
			v, err := m.Read(inst.MemMode(1), ip+1)
			if err != nil {
				return err
			}
			output <- v
			ip += 2
		case OpJumpIfTrue:
			a, err := m.Read(inst.MemMode(1), ip+1)
			if err != nil {
				return err
			}

			if a != 0 {
				ip, err = m.Read(inst.MemMode(2), ip+2)
				if err != nil {
					return err
				}
			} else {
				ip += 3
			}
		case OpJumpIfFalse:
			a, err := m.Read(inst.MemMode(1), ip+1)
			if err != nil {
				return err
			}

			if a == 0 {
				ip, err = m.Read(inst.MemMode(2), ip+2)
				if err != nil {
					return err
				}
			} else {
				ip += 3
			}
		case OpLessThan:
			if err := binaryWriteOp(func(a, b int) int {
				if a < b {
					return 1
				}
				return 0
			}); err != nil {
				return err
			}
		case OpEquals:
			if err := binaryWriteOp(func(a, b int) int {
				if a == b {
					return 1
				}
				return 0
			}); err != nil {
				return err
			}
		case OpHalt:
			return nil
		default:
			return fmt.Errorf("Unknown operation: %v at %v", inst, ip)
		}
	}
}

func runProgramWithInput(m Memory, inputArr []int) ([]int, error) {
	var outputArr []int

	input := make(chan int)
	output := make(chan int)
	halted := make(chan struct{})
	doneCollectingOutput := make(chan struct{})

	go func() {
		for _, v := range inputArr {
			select {
			case input <- v:
				continue
			case <-halted:
				return
			}
		}
		<-halted
	}()

	go func() {
		for v := range output {
			outputArr = append(outputArr, v)
		}
		doneCollectingOutput <- struct{}{}
	}()

	err := runProgram(m, input, output)
	halted <- struct{}{}
	close(output)
	<-doneCollectingOutput

	if err != nil {
		return nil, err
	} else {
		return outputArr, nil
	}
}

func main() {
	p, err := readAsCSVInts(os.Stdin)
	if err != nil {
		fmt.Printf("Could not parse input: %v", err)
		os.Exit(1)
	}

	output, err := runProgramWithInput(NewMemory(p), []int{1})
	if err != nil {
		fmt.Printf("Part 1, error: %v\n", err)
	} else {
		fmt.Printf("Part 1: %v\n", output)
	}

	output, err = runProgramWithInput(NewMemory(p), []int{5})
	if err != nil {
		fmt.Printf("Part 2, error: %v\n", err)
	} else {
		fmt.Printf("Part 2: %v\n", output)
	}
}
