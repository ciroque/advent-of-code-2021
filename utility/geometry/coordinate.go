package geometry

import "fmt"

type Coordinate struct {
	X, Y int
}

func NewCoordinate(x, y int) Coordinate {
	return Coordinate{X: x, Y: y}
}

func (c *Coordinate) Adjacent() []Coordinate {
	north := Coordinate{X: c.X, Y: c.Y - 1}
	south := Coordinate{X: c.X, Y: c.Y + 1}
	east := Coordinate{X: c.X + 1, Y: c.Y}
	west := Coordinate{X: c.X - 1, Y: c.Y}

	return []Coordinate{north, south, east, west}
}

func (c *Coordinate) String() string {
	return fmt.Sprintf("{ %v, %v }", c.X, c.Y)
}
