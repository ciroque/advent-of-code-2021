package geometry

type Coordinate struct {
	x int
	y int
}

func NewCoordinate(x, y int) Coordinate {
	return Coordinate{x: x, y: y}
}
