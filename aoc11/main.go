package main

import (
	"fmt"
	"os"
	"sync"
)

func runOneComputer(code []Word, input <-chan Word, output chan<- Word) error {
	var err error

	wg := sync.WaitGroup{}

	p := NewProgram(code)
	wg.Add(1)
	go func() {
		err = p.Execute(input, output)
		close(output)
		wg.Done()
	}()

	wg.Wait()

	if err != nil {
		return err
	}
	return nil
}

const (
	Up    = 0
	Left  = 1
	Down  = 2
	Right = 3
)

const initGridSize = 120

type Painted struct {
	painted bool
	color   Word
}

func doTheRobot(initColor Word, input <-chan Word, output chan<- Word) (int, [][]Painted) {
	rows := make([][]Painted, initGridSize)
	for y, _ := range rows {
		rows[y] = make([]Painted, initGridSize)
	}

	painted := 0
	x := initGridSize / 2
	y := x
	facing := Up

	steps := 0
	cur := &rows[y][x]
	cur.color = initColor
	output <- cur.color
	for color := range input {
		direction := <-input
		steps += 1
		cur.color = color
		if !cur.painted {
			painted += 1
			cur.painted = true
		}
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
		cur = &rows[y][x]
		output <- cur.color
	}
	close(output)

	return painted, rows
}

func main() {
	code, err := readAsCommanSeparatedWords(os.Stdin)
	if err != nil {
		fmt.Printf("Could not parse input: %v", err)
		os.Exit(1)
	}

	robotToProgram := make(chan Word, 1)
	programToRobot := make(chan Word, 1)

	var painted int

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		painted, _ = doTheRobot(0, programToRobot, robotToProgram)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := runOneComputer(code, robotToProgram, programToRobot)
		if err != nil {
			fmt.Printf("err %v\n", err)
		}
		wg.Done()
	}()

	wg.Wait()

	robotToProgram = make(chan Word, 1)
	programToRobot = make(chan Word, 1)

	var rows [][]Painted
	wg.Add(1)
	go func() {
		_, rows = doTheRobot(1, programToRobot, robotToProgram)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := runOneComputer(code, robotToProgram, programToRobot)
		if err != nil {
			fmt.Printf("err %v\n", err)
		}
		wg.Done()
	}()

	wg.Wait()

	for _, row := range rows {
		for _, sq := range row {
			if sq.color == 0 {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Print("\n")
	}
}
