package main

import (
	"fmt"
)

type MemMode int

const (
	Positional MemMode = iota + 0
	Immediate
)

type Memory []int

func NewMemory(orig []int) Memory {
	m := make(Memory, len(orig))
	copy(m, orig)
	return m
}

func (m Memory) Copy() Memory {
	return NewMemory(m)
}

func (m Memory) CheckBounds(pos int) error {
	if pos < 0 {
		return fmt.Errorf("Address < 0: %v", pos)
	} else if pos >= len(m) {
		return fmt.Errorf("Address > max len: %v > %v", pos, len(m))
	}
	return nil
}

func (m Memory) Read(mode MemMode, pos int) (int, error) {
	switch mode {
	case Immediate:
		return m.ReadImmediate(pos)
	case Positional:
		return m.ReadPositional(pos)
	default:
		return 0, fmt.Errorf("Unknown memory mode: %v", mode)
	}
}

func (m Memory) Write(mode MemMode, pos, val int) error {
	switch mode {
	case Immediate:
		return m.WriteImmediate(pos, val)
	case Positional:
		return m.WritePositional(pos, val)
	default:
		return fmt.Errorf("Unknown memory mode: %v", mode)
	}
}

func (m Memory) ReadImmediate(pos int) (int, error) {
	if err := m.CheckBounds(pos); err != nil {
		return 0, err
	}
	return m[pos], nil
}

func (m Memory) WriteImmediate(pos, val int) error {
	if err := m.CheckBounds(pos); err != nil {
		return err
	}
	m[pos] = val
	return nil
}

func (m Memory) ReadPositional(pos int) (int, error) {
	addr, err := m.ReadImmediate(pos)
	if err != nil {
		return 0, err
	}
	return m.ReadImmediate(addr)
}

func (m Memory) WritePositional(pos, val int) error {
	addr, err := m.ReadImmediate(pos)
	if err != nil {
		return err
	}
	return m.WriteImmediate(addr, val)
}
