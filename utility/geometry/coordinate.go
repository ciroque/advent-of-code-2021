package geometry

type Coordinate struct {
	X, Y int
}

func NewCoordinate(x, y int) Coordinate {
	return Coordinate{X: x, Y: y}
}
