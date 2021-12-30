package main

import (
	"advent-of-code-2021/utility/geometry"
	"fmt"
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
	Solution implementation
*/

type Fold struct {
	axis  geometry.Axis
	index int
}

type Puzzle struct {
	coordinates []map[geometry.Coordinate]int
	folds       []Fold
}

func NewPuzzle(data []string) Puzzle {
	inFoldDefs := false
	initialCoordinates := make(map[geometry.Coordinate]int)
	var folds []Fold
	for _, line := range data {
		if len(line) == 0 {
			inFoldDefs = true
			continue
		}
		if inFoldDefs {
			components := strings.Fields(line)
			parts := strings.Split(components[2], "=")
			var axis geometry.Axis
			if parts[0] == "x" {
				axis = geometry.Horizontal
			} else {
				axis = geometry.Vertical
			}
			index, _ := strconv.Atoi(parts[1])
			folds = append(folds, Fold{axis: axis, index: index})
		} else {
			points := strings.Split(line, ",")
			abscissa, _ := strconv.Atoi(points[0])
			ordinate, _ := strconv.Atoi(points[1])
			initialCoordinates[geometry.NewCoordinate(abscissa, ordinate)]++
		}
	}

	var coordinates []map[geometry.Coordinate]int
	coordinates = append(coordinates, initialCoordinates)

	return Puzzle{
		coordinates: coordinates,
		folds:       folds,
	}
}

func (p *Puzzle) FoldAt(axis geometry.Axis, index int) int {
	nextCoordinates := make(map[geometry.Coordinate]int)
	for coordinate := range p.coordinates[p.LastFold()] {
		if axis == geometry.Vertical {
			if coordinate.Y > index {
				updatedCoordinate := geometry.NewCoordinate(coordinate.X, (index*2)-coordinate.Y)
				nextCoordinates[updatedCoordinate]++
			} else {
				nextCoordinates[coordinate]++
			}
		} else {
			if coordinate.X > index {
				updatedCoordinate := geometry.NewCoordinate((index*2)-coordinate.X, coordinate.Y)
				nextCoordinates[updatedCoordinate]++
			} else {
				nextCoordinates[coordinate]++
			}
		}
	}

	p.coordinates = append(p.coordinates, nextCoordinates)

	return len(nextCoordinates)
}

func (p *Puzzle) LastFold() int {
	return len(p.coordinates) - 1
}

func (p *Puzzle) Height() int {
	index := p.LastFold()
	height := 0
	for coordinate := range p.coordinates[index] {
		if coordinate.Y > height {
			height = coordinate.Y
		}
	}
	return height
}

func (p *Puzzle) Width() int {
	index := p.LastFold()
	width := 0
	for coordinate := range p.coordinates[index] {
		if coordinate.X > width {
			width = coordinate.X
		}
	}
	return width
}

func (p *Puzzle) Print() {
	width := p.Width()
	height := p.Height()
	index := p.LastFold()

	for y := 0; y <= height; y++ {
		for x := 0; x <= width; x++ {
			if _, found := p.coordinates[index][geometry.NewCoordinate(x, y)]; found {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func partOne(puzzle Puzzle) int {
	fold := puzzle.folds[0]
	return puzzle.FoldAt(fold.axis, fold.index)
}

//#..#...##.###..#..#.####.#..#.###...##.
//#.#.....#.#..#.#.#..#....#..#.#..#.#..#
//##......#.###..##...###..#..#.###..#...
//#.#.....#.#..#.#.#..#....#..#.#..#.#.##
//#.#..#..#.#..#.#.#..#....#..#.#..#.#..#
//#..#..##..###..#..#.####..##..###...###
func partTwo(puzzle Puzzle) int {
	solution := 0
	for _, fold := range puzzle.folds {
		solution = puzzle.FoldAt(fold.axis, fold.index)
	}
	puzzle.Print()
	return solution
}

func FindSolutionForInput(filename string, operation func(puzzle Puzzle) int) int {
	puzzleInput := loadPuzzleInput(filename)
	puzzle := NewPuzzle(puzzleInput)
	return operation(puzzle)
}

/*
	Main
*/

type Result struct {
	answer   int
	duration int64
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	waitCount := 4
	var waitGroup sync.WaitGroup
	waitGroup.Add(waitCount)

	exampleChannelOne := make(chan Result)
	exampleChannelTwo := make(chan Result)
	partOneChannel := make(chan Result)
	partTwoChannel := make(chan Result)

	go doExampleOne(exampleChannelOne, &waitGroup)
	go doExampleTwo(exampleChannelTwo, &waitGroup)
	go doPartOne(partOneChannel, &waitGroup)
	go doPartTwo(partTwoChannel, &waitGroup)

	exampleResultOne := <-exampleChannelOne
	exampleResultTwo := <-exampleChannelTwo
	partOneResult := <-partOneChannel
	partTwoResult := <-partTwoChannel

	waitGroup.Wait()

	log.
		Info().
		Int("example-one-answer", exampleResultOne.answer).
		Int64("example-one-duration", exampleResultOne.duration).
		Int("example-two-answer", exampleResultTwo.answer).
		Int64("example-two-duration", exampleResultTwo.duration).
		Int("part-one-answer", partOneResult.answer).
		Int64("part-one-duration", partOneResult.duration).
		Int("part-two-answer", partTwoResult.answer).
		Int64("part-two-duration", partTwoResult.duration).
		Msg("day 04")
}

/*
	Executors
*/

func doExampleOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("example-input.dat", partOne),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doExampleTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("example-input.dat", partTwo),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("puzzle-input.dat", partOne),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("puzzle-input.dat", partTwo),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput(filename string) []string {
	return support.ReadFileIntoLines(filename)
}
