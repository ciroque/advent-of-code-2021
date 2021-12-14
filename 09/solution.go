package main

import (
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

const highestPoint = "9"

//type Coordinate struct {
//	x, y int
//}

type FloorMap struct {
	dimX int
	dimY int

	heights      [][]int
	lowestPoints map[Coordinate]int
}

func (fm *FloorMap) Append(heights string) *FloorMap {
	yi := len(fm.heights)
	fm.heights = append(fm.heights, []int{})

	for _, height := range heights {
		value, _ := strconv.Atoi(string(height))
		fm.heights[yi] = append(fm.heights[yi], value)
	}

	return fm
}

func (fm *FloorMap) CalculateRiskSum() int {
	accumulator := 0
	for _, risk := range fm.lowestPoints {
		accumulator += risk
	}
	return accumulator
}

func (fm *FloorMap) FindLowestPoints() *FloorMap {
	for yi := 1; yi <= fm.dimY; yi++ {
		for xi := 1; xi <= fm.dimX; xi++ {
			coordinate := Coordinate{x: xi, y: yi}
			if fm.IsLowestPoint(coordinate) {
				fm.lowestPoints[coordinate] = 1 + fm.heights[yi][xi]
			}
		}
	}

	return fm
}

func (fm *FloorMap) HeightAt(coordinate Coordinate) int {
	return fm.heights[coordinate.y][coordinate.x]
}

func (fm *FloorMap) IsLowestPoint(coordinate Coordinate) bool {
	north := Coordinate{x: coordinate.x, y: coordinate.y - 1}
	south := Coordinate{x: coordinate.x, y: coordinate.y + 1}
	east := Coordinate{x: coordinate.x + 1, y: coordinate.y}
	west := Coordinate{x: coordinate.x - 1, y: coordinate.y}

	adjacentCoordinates := []Coordinate{north, south, east, west}

	accumulator := 0

	for _, adjacent := range adjacentCoordinates {
		if fm.HeightAt(coordinate) < fm.HeightAt(adjacent) {
			accumulator++
		}
	}

	return accumulator == len(adjacentCoordinates)
}

func (fm *FloorMap) Print() {
	for yi := 0; yi <= fm.dimY+1; yi++ {
		for xi := 0; xi <= fm.dimX+1; xi++ {
			fmt.Printf("%v ", fm.heights[yi][xi])
		}
		fmt.Print("\n")
	}
}

// NewFloorMap
// This creates a 'border' around the input.
// That border is filled with 9s so the introduced values will not affect the determination of the lowest point.
// What it WILL do is relieve the burden of bounds checking. Yay.
func NewFloorMap(puzzleInput []string) FloorMap {
	floorMap := FloorMap{
		dimY: len(puzzleInput),
		dimX: len(puzzleInput[0]),

		heights:      make([][]int, 0),
		lowestPoints: make(map[Coordinate]int),
	}

	rowOfNines := strings.Repeat(highestPoint, floorMap.dimX+2)

	floorMap.Append(rowOfNines)

	for _, line := range puzzleInput {
		borderedLine := fmt.Sprintf("%v%v%v", highestPoint, line, highestPoint)
		floorMap.Append(borderedLine)
	}

	floorMap.Append(rowOfNines)

	return floorMap
}

func CalculatePartOneSolution(floorMap FloorMap) int {
	// Yeah, I know, this reads nicely, but it imposes an order on the methods that could be confusing
	return floorMap.FindLowestPoints().CalculateRiskSum()
}

func FindSolutionForInput(filename string, calculateSolution func(floorMap FloorMap) int) int {
	solution := 0

	puzzleInput := loadPuzzleInput(filename)
	floorMap := NewFloorMap(puzzleInput)

	solution = calculateSolution(floorMap)

	return solution
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
		answer:   FindSolutionForInput("example-input.dat", CalculatePartOneSolution),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doExampleTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   0, //  FindSolutionForInput("example-input.dat"),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("puzzle-input.dat", CalculatePartOneSolution),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   0, // FindSolutionForInput("puzzle-input.dat"),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput(filename string) []string {
	return support.ReadFileIntoLines(filename)
}
