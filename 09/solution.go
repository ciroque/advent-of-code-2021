package main

import (
	"advent-of-code-2021/utility/geometry"
	"fmt"
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
	Solution implementation
*/

const highestPoint = "9"

type FloorMap struct {
	dimX int
	dimY int

	heights      [][]int
	lowestPoints map[geometry.Coordinate]int
	basins       []int
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

func (fm *FloorMap) CalculateBasinSizeProduct() int {
	accumulator := 1
	sort.Sort(sort.Reverse(sort.IntSlice(fm.basins)))
	topThree := fm.basins[0:3]
	for _, size := range topThree {
		accumulator *= size
	}
	return accumulator
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
			coordinate := geometry.NewCoordinate(xi, yi)
			if fm.IsLowestPoint(coordinate) {
				fm.lowestPoints[coordinate] = 1 + fm.heights[yi][xi]
			}
		}
	}

	return fm
}

func (fm *FloorMap) HeightAt(coordinate geometry.Coordinate) int {
	return fm.heights[coordinate.Y][coordinate.X]
}

func (fm *FloorMap) IsLowestPoint(coordinate geometry.Coordinate) bool {
	adjacentCoordinates := coordinate.Adjacent()

	accumulator := 0

	for _, adjacent := range adjacentCoordinates {
		if fm.HeightAt(coordinate) < fm.HeightAt(adjacent) {
			accumulator++
		}
	}

	return accumulator == len(adjacentCoordinates)
}

func (fm *FloorMap) MapBasins() *FloorMap {
	HighestPoint, _ := strconv.Atoi(highestPoint)

	mapBasin := func(coordinate geometry.Coordinate) int {
		// choosing the size of the buffered channel is ham-fisted at this point
		coordinatesChannel := make(chan geometry.Coordinate, (fm.dimX*fm.dimY)/2)
		accumulator := 0
		var visited = make(map[geometry.Coordinate]bool)
		coordinatesChannel <- coordinate

		for {
			select {
			case coordinate := <-coordinatesChannel:
				if _, alreadyVisited := visited[coordinate]; alreadyVisited {
					continue
				}

				visited[coordinate] = true

				accumulator++

				for _, adjacent := range coordinate.Adjacent() {
					if fm.HeightAt(adjacent) != HighestPoint {
						coordinatesChannel <- adjacent
					}
				}
			default:
				return accumulator
			}
		}
	}

	for lowestPoint := range fm.lowestPoints {
		fm.basins = append(fm.basins, mapBasin(lowestPoint))
	}

	return fm
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
		lowestPoints: make(map[geometry.Coordinate]int),
		basins:       []int{},
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

func CalculatePartTwoSolution(floorMap FloorMap) int {
	solution := 0
	solution = floorMap.FindLowestPoints().MapBasins().CalculateBasinSizeProduct()
	return solution
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
		answer:   FindSolutionForInput("example-input.dat", CalculatePartTwoSolution),
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
		answer:   FindSolutionForInput("puzzle-input.dat", CalculatePartTwoSolution),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput(filename string) []string {
	return support.ReadFileIntoLines(filename)
}
