package main

import (
	"fmt"
	"math/bits"
	"os"
)

func applyGravityUnit(p1, p2, v1, v2 *int) {
	if *p1 < *p2 {
		*v1 += 1
		*v2 -= 1
	} else if *p1 > *p2 {
		*v1 -= 1
		*v2 += 1
	}
}

func applyGravity(m0, m1 *Moon) {
	applyGravityUnit(&m0.Pos.X, &m1.Pos.X, &m0.Vel.X, &m1.Vel.X)
	applyGravityUnit(&m0.Pos.Y, &m1.Pos.Y, &m0.Vel.Y, &m1.Vel.Y)
	applyGravityUnit(&m0.Pos.Z, &m1.Pos.Z, &m0.Vel.Z, &m1.Vel.Z)
}

type Vec3 struct {
	X int
	Y int
	Z int
}

type Moon struct {
	Pos Vec3
	Vel Vec3
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (v Vec3) Energy() int {
	return abs(v.X) + abs(v.Y) + abs(v.Z)
}

func (v Vec3) String() string {
	return fmt.Sprintf("<x=%d, y=%d, z=%d>", v.X, v.Y, v.Z)
}

func (m Moon) Energy() int {
	return m.Pos.Energy() * m.Vel.Energy()
}

func (m Moon) String() string {
	return fmt.Sprintf("pos=%v vel=%v", m.Pos, m.Vel)
}

func simulate(moons []Moon) {
	for i := 0; i < len(moons); i += 1 {
		m0 := &moons[i]
		for j := i + 1; j < len(moons); j += 1 {
			m1 := &moons[j]
			applyGravity(m0, m1)
		}
	}

	for i := range moons {
		moon := &moons[i]
		moon.Pos.X += moon.Vel.X
		moon.Pos.Y += moon.Vel.Y
		moon.Pos.Z += moon.Vel.Z
	}
}

func signum(x int) int {
	return int((x >> (bits.UintSize - 1)) | int(uint(-x)>>(bits.UintSize-1)))
}

func step(p *[4]int, v *[4]int) {
	d := signum(p[1] - p[0])
	v[0] += d
	v[1] -= d
	d = signum(p[2] - p[0])
	v[0] += d
	v[2] -= d
	d = signum(p[3] - p[0])
	v[0] += d
	v[3] -= d
	d = signum(p[2] - p[1])
	v[1] += d
	v[2] -= d
	d = signum(p[3] - p[1])
	v[1] += d
	v[3] -= d
	d = signum(p[3] - p[2])
	v[2] += d
	v[3] -= d
	p[0] += v[0]
	p[1] += v[1]
	p[2] += v[2]
	p[3] += v[3]
}

func findloop(o [4]int) int {
	var p [4]int
	var v [4]int

	p = o

	i := 1
	for {
		step(&p, &v)
		if p == o {
			return i + 1
		}
		i += 1
	}
	return -1

}

func absGCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	orig, err := readMoons(os.Stdin)
	if err != nil {
		fmt.Printf("Could not parse input: %v\n", err)
		os.Exit(1)
	}

	moons := make([]Moon, len(orig))
	copy(moons, orig)

	for n := 0; n < 1000; n += 1 {
		simulate(moons)
	}

	totalEnergy := 0
	for _, m := range moons {
		totalEnergy += m.Energy()
	}
	fmt.Printf("Part 1: %v\n", totalEnergy)

	x := findloop([4]int{
		orig[0].Pos.X,
		orig[1].Pos.X,
		orig[2].Pos.X,
		orig[3].Pos.X,
	})
	y := findloop([4]int{
		orig[0].Pos.Y,
		orig[1].Pos.Y,
		orig[2].Pos.Y,
		orig[3].Pos.Y,
	})
	z := findloop([4]int{
		orig[0].Pos.Z,
		orig[1].Pos.Z,
		orig[2].Pos.Z,
		orig[3].Pos.Z,
	})
	gcd := absGCD(absGCD(x, y), z)
	x0 := int64(x / gcd)
	y0 := int64(y / gcd)
	z0 := int64(z / gcd)
	lcm := x0 * y0 * z0
	fmt.Printf("Part 2: %v\n", lcm)
}
