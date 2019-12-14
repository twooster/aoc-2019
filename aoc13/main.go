package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

func runProgram(code []Word, input <-chan Word, output chan<- Word) error {
	defer close(output)
	p := NewProgram(code)
	return p.Execute(input, output)
}

func doTheScreen(c *Canvas, input <-chan Word, output chan<- Word) {
	for {
		x, ok := <-input
		if !ok {
			break
		}
		y := <-input
		tile := <-input
		c.SetColor(int(x), int(y), tile)
	}
	close(output)
}

func part1(code []Word) {
	screenToProgram := make(chan Word, 1)
	programToScreen := make(chan Word, 1)
	canvas := NewCanvas()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		doTheScreen(canvas, programToScreen, screenToProgram)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := runProgram(code, screenToProgram, programToScreen)
		if err != nil {
			fmt.Printf("err %v\n", err)
		}
		wg.Done()
	}()

	wg.Wait()

	blockTiles := 0
	for y := canvas.MinY; y <= canvas.MaxY; y += 1 {
		var sb strings.Builder
		for x := canvas.MinX; x <= canvas.MaxX; x += 1 {
			color := canvas.GetColor(x, y)
			switch color {
			case 0:
				sb.WriteRune(' ')
			case 1:
				sb.WriteRune('#')
			case 2:
				sb.WriteRune('x')
				blockTiles += 1
			case 3:
				sb.WriteRune('-')
			case 4:
				sb.WriteRune('o')
			default:
				fmt.Printf("color %v\n")
			}
		}
		sb.WriteRune('\n')
		fmt.Print(sb.String())
	}

	fmt.Printf("Part 1: %v\n", blockTiles)
}

func calcJoystickMove(canvas *Canvas) Word {
	ballX := -1
	paddleX := -1
outer:
	for y := canvas.MinY; y <= canvas.MaxY; y += 1 {
		for x := canvas.MinX; x <= canvas.MaxX; x += 1 {
			color := canvas.GetColor(x, y)
			if color == 4 {
				ballX = x
				if paddleX != -1 {
					break outer
				}
			} else if color == 3 {
				paddleX = x
				if ballX != -1 {
					break outer
				}
			}
		}
	}
	if ballX < paddleX {
		return -1
	} else if ballX > paddleX {
		return 1
	}
	return 0
}

const ballLoc = 388
const paddleLoc = 392

func part2Cheating(orig []Word) {
	code := make([]Word, len(orig))
	copy(code, orig)
	code[0] = 2

	screenToProgram := make(chan Word, 1)
	programToScreen := make(chan Word, 1)
	canvas := NewCanvas()

	var score Word

	move := Word(1)

	potentialBallLoc := make(map[int]int)
	potentialPaddleLoc := make(map[int]int)

	program := NewProgram(code)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		var err error
		var done bool
		for !done && err == nil {
			done, err = program.Step(screenToProgram, programToScreen)
			b, e := program.ReadImmediate(Word(ballLoc))
			if e != nil {
				err = e
				break
			}
			err = program.WriteImmediate(Word(paddleLoc), b)
		}
		if err != nil {
			fmt.Printf("err %v\n", err)
		}
		close(programToScreen)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for {
			select {
			case screenToProgram <- move:
				//noop
			default:
				//noop
			}

			x, ok := <-programToScreen
			if !ok {
				break
			}
			y := <-programToScreen
			tile := <-programToScreen
			if x == -1 {
				score = tile
			} else {
				if tile == 4 {
					//fmt.Printf("ball   y, x: %v %v\n", y, x)
					for i, v := range program.memory {
						if v == x {
							_, ok := potentialBallLoc[i]
							if !ok {
								potentialBallLoc[i] = 1
							} else {
								potentialBallLoc[i] += 1
							}
						}
					}
				} else if tile == 3 {
					//fmt.Printf("paddle y, x: %v %v\n", y, x)
					for i, v := range program.memory {
						if v == x {
							_, ok := potentialPaddleLoc[i]
							if !ok {
								potentialPaddleLoc[i] = 1
							} else {
								potentialPaddleLoc[i] += 1
							}
						}
					}
				}
				canvas.SetColor(int(x), int(y), tile)
			}
		}
		close(screenToProgram)
		wg.Done()
	}()

	wg.Wait()

	//fmt.Printf("Potential ball locations: %v\n", potentialBallLoc)
	//fmt.Printf("Potential padl locations: %v\n", potentialPaddleLoc)
	fmt.Printf("Part 2: %v\n", score)
}

func main() {
	code, err := readAsCommaSeparatedWords(os.Stdin)
	if err != nil {
		fmt.Printf("Could not parse input: %v", err)
		os.Exit(1)
	}

	part1(code)
	part2Cheating(code)
}
