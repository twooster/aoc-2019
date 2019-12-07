package main

import (
	"errors"
	"fmt"
)

type Operation int
type Instruction int

const (
	OpAdd         Operation = 1
	OpMul                   = 2
	OpInput                 = 3
	OpOutput                = 4
	OpJumpNonZero           = 5
	OpJumpZero              = 6
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

type binaryWriteFn func(a, b int) int
type unaryJumpFn func(a int) bool

func add(a, b int) int {
	return a + b
}

func mul(a, b int) int {
	return a * b
}

func lessThan(a, b int) int {
	if a < b {
		return 1
	}
	return 0
}

func equals(a, b int) int {
	if a == b {
		return 1
	}
	return 0
}

func jumpZero(a int) bool {
	return a == 0
}

func jumpNonZero(a int) bool {
	return a != 0
}

type Program struct {
	memory Memory
	ip     int
}

func (c *Program) BinaryWrite(inst Instruction, fn binaryWriteFn) error {
	a, err := c.memory.Read(inst.MemMode(1), c.ip+1)
	if err != nil {
		return err
	}

	b, err := c.memory.Read(inst.MemMode(2), c.ip+2)
	if err != nil {
		return err
	}

	if err := c.memory.Write(inst.MemMode(3), c.ip+3, fn(a, b)); err != nil {
		return err
	}

	c.ip += 4
	return nil
}

func (c *Program) UnaryJump(inst Instruction, fn unaryJumpFn) error {
	a, err := c.memory.Read(inst.MemMode(1), c.ip+1)
	if err != nil {
		return err
	}

	if fn(a) {
		c.ip, err = c.memory.Read(inst.MemMode(2), c.ip+2)
		if err != nil {
			return err
		}
	} else {
		c.ip += 3
	}
	return nil
}

func (c *Program) Step(input <-chan int, output chan<- int) (done bool, err error) {
	var v int
	var inst Instruction
	if v, err = c.memory.ReadImmediate(c.ip); err != nil {
		return
	} else {
		inst = Instruction(v)
	}

	switch inst.Operation() {
	case OpAdd:
		err = c.BinaryWrite(inst, add)
	case OpMul:
		err = c.BinaryWrite(inst, mul)
	case OpInput:
		if v, ok := <-input; ok {
			err = c.memory.Write(inst.MemMode(1), c.ip+1, v)
			if err == nil {
				c.ip += 2
			}
		} else {
			err = errors.New("Input EOF")
		}
	case OpOutput:
		if v, err = c.memory.Read(inst.MemMode(1), c.ip+1); err == nil {
			output <- v
			c.ip += 2
		}
	case OpJumpNonZero:
		err = c.UnaryJump(inst, jumpNonZero)
	case OpJumpZero:
		err = c.UnaryJump(inst, jumpZero)
	case OpLessThan:
		err = c.BinaryWrite(inst, lessThan)
	case OpEquals:
		err = c.BinaryWrite(inst, equals)
	case OpHalt:
		done = true
	default:
		err = fmt.Errorf("Unknown operation: %v at %v", inst, c.ip)
	}
	return done, err
}

func (c *Program) Execute(input <-chan int, output chan<- int) (err error) {
	var done bool
	for !done && err == nil {
		done, err = c.Step(input, output)
	}
	return err
}
