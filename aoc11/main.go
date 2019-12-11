package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

func runProgram(code []Word, input <-chan Word, output chan<- Word) error {
	defer closu(output)
	p := NewProgram(code)
	return p.Execute(input, output)
}

type Direction int

const (
	Up Direction = iota
	Left
	Down
	Right
)

func doTheRobot(c *Canvas, input <-chan Word, output chan<- Word) {
	x := 0
	y := 0
	facing := Up

	output <- c.GetColor(x, y)
	for color := range input {
		c.SetColor(x, y, color)

		direction := <-input
		switch facing {
		case Up:
			if direction == 0 {
				facing = Left
				x -= 1
			} else {
				facing = Right
				x += 1
			}
		case Left:
			if direction == 0 {
				facing = Down
				y += 1
			} else {
				facing = Up
				y -= 1
			}
		case Down:
			if direction == 0 {
				facing = Right
				x += 1
			} else {
				facing = Left
				x -= 1
			}
		case Right:
			if direction == 0 {
				facing = Up
				y -= 1
			} else {
				facing = Down
				y += 1
			}
		}
		output <- c.GetColor(x, y)
	}
	close(output)
}

func main() {
	code, err := readAsCommaSeparatedWords(os.Stdin)
	if err != nil {
		fmt.Printf("Could not parse input: %v", err)
		os.Exit(1)
	}

	robotToProgram := make(chan Word, 1)
	programToRobot := make(chan Word, 1)
	canvas := NewCanvas()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		doTheRobot(canvas, programToRobot, robotToProgram)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := runProgram(code, robotToProgram, programToRobot)
		if err != nil {
			fmt.Printf("err %v\n", err)
		}
		wg.Done()
	}()

	wg.Wait()

	fmt.Printf("Part 1: %v\n", canvas.Painted)

	robotToProgram = make(chan Word, 1)
	programToRobot = make(chan Word, 1)
	canvas = NewCanvas()
	canvas.SetColor(0, 0, 1)

	wg.Add(1)
	go func() {
		doTheRobot(canvas, programToRobot, robotToProgram)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := runProgram(code, robotToProgram, programToRobot)
		if err != nil {
			fmt.Printf("err %v\n", err)
		}
		wg.Done()
	}()

	wg.Wait()

	for y := canvas.MinY; y <= canvas.MaxY; y += 1 {
		var sb strings.Builder
		for x := canvas.MinX; x <= canvas.MaxX; x += 1 {
			color := canvas.GetColor(x, y)
			if color == 0 {
				sb.WriteRune(' ')
			} else {
				sb.WriteRune('#')
			}
		}
		sb.WriteRune('\n')
		fmt.Print(sb.String())
	}
}
