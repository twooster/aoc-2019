package main

type Canvas struct {
	MinX         int
	MinY         int
	MaxX         int
	MaxY         int
	Painted      int
	DefaultColor Word
	rows         map[int]map[int]Word
}

func NewCanvas() *Canvas {
	return &Canvas{
		rows: make(map[int]map[int]Word),
	}
}

func (c *Canvas) GetColor(x, y int) Word {
	if row, ok := c.rows[y]; !ok {
		return c.DefaultColor
	} else if color, ok := row[x]; !ok {
		return c.DefaultColor
	} else {
		return color
	}
}

func (c *Canvas) SetColor(x, y int, color Word) {
	if x < c.MinX {
		c.MinX = x
	} else if x > c.MaxX {
		c.MaxX = x
	}

	if y < c.MinY {
		c.MinY = y
	} else if y > c.MaxY {
		c.MaxY = y
	}

	row, ok := c.rows[y]
	if !ok {
		row = make(map[int]Word)
		c.rows[y] = row
	}
	if _, exists := row[x]; !exists {
		c.Painted += 1
	}
	row[x] = color
}
