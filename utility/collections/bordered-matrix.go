package collections

import (
	"advent-of-code-2021/utility/geometry"
	"fmt"
	"strconv"
)

type BorderedIntMatrix struct {
	borderValue, height, width int
	matrix                     [][]int
}

func NewBorderedIntMatrix() BorderedIntMatrix {
	return BorderedIntMatrix{
		width:  0,
		height: 0,
		matrix: [][]int{},
	}
}

func (b *BorderedIntMatrix) ForEachAdjacentIn(coordinates []geometry.Coordinate, process func(coordinate geometry.Coordinate, value int) int) {
	for _, coordinate := range coordinates {
		for _, adjacent := range coordinate.AllAdjacent() {
			value := b.matrix[adjacent.Y][adjacent.X]
			if value == b.borderValue {
				continue
			}

			b.matrix[adjacent.Y][adjacent.X] = process(adjacent, value)
		}
	}
}

func (b *BorderedIntMatrix) ForEachIn(coordinates []geometry.Coordinate, update func(v int) int) {
	for _, coordinate := range coordinates {
		b.matrix[coordinate.Y][coordinate.X] = update(b.matrix[coordinate.Y][coordinate.X])
	}
}

func (b *BorderedIntMatrix) Populate(input []string, borderValue int) {
	b.borderValue = borderValue
	b.width = len(input[0]) + 2
	b.height = len(input) + 2

	addPadRow := func() {
		b.matrix = append(b.matrix, []int{})
		rowIndex := len(b.matrix) - 1

		for index := 0; index < b.width; index++ {
			b.matrix[rowIndex] = append(b.matrix[rowIndex], borderValue)
		}
	}

	addRow := func(line string) {
		b.matrix = append(b.matrix, []int{})
		rowIndex := len(b.matrix) - 1
		b.matrix[rowIndex] = append(b.matrix[rowIndex], borderValue)

		for _, char := range line {
			value, _ := strconv.Atoi(string(char))
			b.matrix[rowIndex] = append(b.matrix[rowIndex], value)
		}

		b.matrix[rowIndex] = append(b.matrix[rowIndex], borderValue)
	}

	addPadRow()
	for _, line := range input {
		addRow(line)
	}
	addPadRow()
}

func (b *BorderedIntMatrix) Print() {
	for yi := 0; yi < b.height; yi++ {
		for xi := 0; xi < b.width; xi++ {
			fmt.Printf("%2v ", b.ValueAt(xi, yi))
		}
		fmt.Println()
	}
	fmt.Println()
}

func (b *BorderedIntMatrix) Size() int {
	return (b.height - 2) * (b.width - 2)
}

func (b *BorderedIntMatrix) ValueAt(x, y int) int {
	return b.matrix[y][x]
}

func (b *BorderedIntMatrix) VisitEach(visit func(coordinate geometry.Coordinate, v int) int) {
	for yi := 1; yi < b.height-1; yi++ {
		for xi := 1; xi < b.width-1; xi++ {
			coordinate := geometry.NewCoordinate(xi, yi)
			b.matrix[yi][xi] = visit(coordinate, b.matrix[yi][xi])
		}
	}
}
