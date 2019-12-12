package main

import (
	"fmt"
	"io"
)

func readMoons(input io.Reader) ([]Moon, error) {
	var moons []Moon
	for {
		moon := Moon{}
		n, _ := fmt.Fscanf(input, "<x=%d, y=%d, z=%d>\n", &moon.Pos.X, &moon.Pos.Y, &moon.Pos.Z)
		if n != 3 {
			return moons, nil
		}
		moons = append(moons, moon)
	}
	return moons, nil
}
