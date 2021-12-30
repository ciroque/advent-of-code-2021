package main

import (
	"advent-of-code-2021/utility/geometry"
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

func (p *Puzzle) FoldAt(axis geometry.Axis, index int) []geometry.Coordinate {
	var foldedCoordinates []geometry.Coordinate

	return foldedCoordinates
}

func FindSolutionForInput(filename string) int {
	solution := 0

	puzzleInput := loadPuzzleInput(filename)
	puzzle := NewPuzzle(puzzleInput)

	return solution + len(puzzle.coordinates)
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
		answer:   FindSolutionForInput("example-input.dat"),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doExampleTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   0, // FindSolutionForInput("example-input.dat"),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   0, //  FindSolutionForInput("puzzle-input.dat"),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   0, //  FindSolutionForInput("puzzle-input.dat"),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput(filename string) []string {
	return support.ReadFileIntoLines(filename)
}
