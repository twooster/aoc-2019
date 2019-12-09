package main

import (
	"errors"
	"fmt"
)

type Instruction int64
type Operation int
type MemMode int

const (
	Positional MemMode = iota + 0
	Immediate
	Relative
)

const (
	OpAdd         Operation = 1
	OpMul                   = 2
	OpInput                 = 3
	OpOutput                = 4
	OpJumpNonZero           = 5
	OpJumpZero              = 6
	OpLessThan              = 7
	OpEquals                = 8
	OpSetRelBase            = 9
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

type binaryWriteFn func(a, b int64) int64
type unaryJumpFn func(a int64) bool

func add(a, b int64) int64 {
	return a + b
}

func mul(a, b int64) int64 {
	return a * b
}

func lessThan(a, b int64) int64 {
	if a < b {
		return 1
	}
	return 0
}

func equals(a, b int64) int64 {
	if a == b {
		return 1
	}
	return 0
}

func jumpZero(a int64) bool {
	return a == 0
}

func jumpNonZero(a int64) bool {
	return a != 0
}

type Program struct {
	memory []int64
	ip     int64
	rb     int64
}

func NewProgram(mem []int64) *Program {
	c := make([]int64, len(mem)+1000)
	copy(c, mem)
	return &Program{
		memory: c,
	}
}

const MaxMemorySize = 10000

func (p *Program) BinaryWrite(inst Instruction, fn binaryWriteFn) error {
	a, err := p.Read(inst.MemMode(1), p.ip+1)
	if err != nil {
		return err
	}

	b, err := p.Read(inst.MemMode(2), p.ip+2)
	if err != nil {
		return err
	}

	if err := p.Write(inst.MemMode(3), p.ip+3, fn(a, b)); err != nil {
		return err
	}

	p.ip += 4
	return nil
}

func (p *Program) UnaryJump(inst Instruction, fn unaryJumpFn) error {
	a, err := p.Read(inst.MemMode(1), p.ip+1)
	if err != nil {
		return err
	}

	if fn(a) {
		p.ip, err = p.Read(inst.MemMode(2), p.ip+2)
		if err != nil {
			return err
		}
	} else {
		p.ip += 3
	}
	return nil
}

func (p *Program) Step(input <-chan int64, output chan<- int64) (done bool, err error) {
	var v int64
	var inst Instruction

	if v, err = p.ReadImmediate(p.ip); err != nil {
		return
	} else {
		inst = Instruction(v)
	}

	switch inst.Operation() {
	case OpAdd:
		err = p.BinaryWrite(inst, add)
	case OpMul:
		err = p.BinaryWrite(inst, mul)
	case OpInput:
		if v, ok := <-input; ok {
			err = p.Write(inst.MemMode(1), p.ip+1, v)
			if err == nil {
				p.ip += 2
			}
		} else {
			err = errors.New("Input EOF")
		}
	case OpOutput:
		if v, err = p.Read(inst.MemMode(1), p.ip+1); err == nil {
			output <- v
			p.ip += 2
		}
	case OpJumpNonZero:
		err = p.UnaryJump(inst, jumpNonZero)
	case OpJumpZero:
		err = p.UnaryJump(inst, jumpZero)
	case OpLessThan:
		err = p.BinaryWrite(inst, lessThan)
	case OpEquals:
		err = p.BinaryWrite(inst, equals)
	case OpSetRelBase:
		if offset, err := p.Read(inst.MemMode(1), p.ip+1); err != nil {
			return false, err
		} else {
			p.rb += offset
			p.ip += 2
		}
	case OpHalt:
		done = true
	default:
		err = fmt.Errorf("Unknown operation: %v at %v", inst, p.ip)
	}
	return done, err
}

func (p *Program) Execute(input <-chan int64, output chan<- int64) (err error) {
	var done bool
	for !done && err == nil {
		done, err = p.Step(input, output)
	}
	return err
}

func (p *Program) EnsureBounds(pos int64) error {
	if pos < 0 {
		return fmt.Errorf("Address < 0: %v", pos)
	} else if pos >= MaxMemorySize {
		return fmt.Errorf("Address >= MAX_MEM: %v", pos)
	} else if cur := int64(len(p.memory)); pos >= cur {
		for cur <= pos {
			cur *= 2
		}
		c := make([]int64, cur)
		copy(c, p.memory)
		p.memory = c
	}
	return nil
}

func (p *Program) Read(mode MemMode, pos int64) (int64, error) {
	switch mode {
	case Immediate:
		return p.ReadImmediate(pos)
	case Positional:
		return p.ReadPositional(0, pos)
	case Relative:
		return p.ReadPositional(p.rb, pos)
	default:
		return 0, fmt.Errorf("Unknown memory mode: %v", mode)
	}
}

func (p *Program) Write(mode MemMode, pos, val int64) error {
	switch mode {
	case Immediate:
		return p.WriteImmediate(pos, val)
	case Positional:
		return p.WritePositional(0, pos, val)
	case Relative:
		return p.WritePositional(p.rb, pos, val)
	default:
		return fmt.Errorf("Unknown memory mode: %v", mode)
	}
}

func (p *Program) ReadImmediate(pos int64) (int64, error) {
	if err := p.EnsureBounds(pos); err != nil {
		return 0, err
	}
	val := p.memory[pos]
	return val, nil
}

func (p *Program) WriteImmediate(pos, val int64) error {
	if err := p.EnsureBounds(pos); err != nil {
		return err
	}
	p.memory[pos] = val
	return nil
}

func (p *Program) ReadPositional(rb, pos int64) (int64, error) {
	if err := p.EnsureBounds(pos); err != nil {
		return 0, err
	}
	addr := p.memory[pos] + rb
	if err := p.EnsureBounds(addr); err != nil {
		return 0, err
	}
	value := p.memory[addr]
	return value, nil
}

func (p *Program) WritePositional(rb, pos, val int64) error {
	if err := p.EnsureBounds(pos); err != nil {
		return err
	}
	addr := p.memory[pos] + rb
	if err := p.EnsureBounds(addr); err != nil {
		return err
	}
	p.memory[addr] = val
	return nil
}
