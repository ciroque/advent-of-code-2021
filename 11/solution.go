package main

import (
	"advent-of-code-2021/utility/collections"
	"advent-of-code-2021/utility/geometry"
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

/*
	Solution implementation
*/

func FindSolutionForInput(filename string) int {
	SentinelValue := -1
	FlashPoint := 9
	flashedCount := 0

	var flashedDuringStep []geometry.Coordinate
	flashedMap := make(map[geometry.Coordinate]int)

	setToZero := func(v int) int { return 0 }
	incrementValue := func(coordinate geometry.Coordinate, v int) int { return v + 1 }
	recordFlashed := func(coordinate geometry.Coordinate, energyLevel int) int {
		_, found := flashedMap[coordinate]
		if energyLevel > FlashPoint && !found {
			flashedDuringStep = append(flashedDuringStep, coordinate)
			flashedMap[coordinate]++
			flashedCount++
		}
		return energyLevel
	}

	matrix := collections.NewBorderedIntMatrix()
	puzzleInput := loadPuzzleInput(filename)
	matrix.Populate(puzzleInput, SentinelValue)

	for j := 0; j < 100; j++ {

		// stage 1
		matrix.VisitEach(incrementValue)

		// stage 2
		for {
			matrix.VisitEach(recordFlashed)
			matrix.ForEachAdjacentIn(flashedDuringStep, incrementValue)
			if len(flashedDuringStep) == 0 {
				break
			}
			flashedDuringStep = []geometry.Coordinate{}
		}

		// stage 3
		for k := range flashedMap {
			flashedDuringStep = append(flashedDuringStep, k)
		}
		flashedMap = make(map[geometry.Coordinate]int)
		matrix.ForEachIn(flashedDuringStep, setToZero)

		flashedDuringStep = []geometry.Coordinate{}
	}

	return flashedCount
}

func FindSolutionForInput2(filename string) int {
	SentinelValue := -1
	FlashPoint := 9

	puzzleInput := loadPuzzleInput(filename)
	matrix := collections.NewBorderedIntMatrix()
	matrix.Populate(puzzleInput, SentinelValue)

	flashedCount := 0
	zeroCount := 0

	var flashedDuringStep []geometry.Coordinate
	flashedMap := make(map[geometry.Coordinate]int)

	setToZero := func(v int) int { return 0 }

	incrementValue := func(coordinate geometry.Coordinate, v int) int { return v + 1 }

	recordFlashed := func(coordinate geometry.Coordinate, energyLevel int) int {
		_, found := flashedMap[coordinate]
		if energyLevel > FlashPoint && !found {
			flashedDuringStep = append(flashedDuringStep, coordinate)
			flashedMap[coordinate]++
			flashedCount++
		}
		return energyLevel
	}

	countZeros := func(coordinate geometry.Coordinate, value int) int {
		if value == 0 {
			zeroCount++
		}

		return value
	}

	for j := 0; j < 500; j++ {

		// stage 1
		matrix.VisitEach(incrementValue)

		// stage 2
		for {
			matrix.VisitEach(recordFlashed)
			matrix.ForEachAdjacentIn(flashedDuringStep, incrementValue)
			if len(flashedDuringStep) == 0 {
				break
			}
			flashedDuringStep = []geometry.Coordinate{}
		}

		// stage 3
		for k := range flashedMap {
			flashedDuringStep = append(flashedDuringStep, k)
		}
		flashedMap = make(map[geometry.Coordinate]int)
		matrix.ForEachIn(flashedDuringStep, setToZero)

		matrix.VisitEach(countZeros)

		if zeroCount == matrix.Size() {
			return j + 1
		}

		zeroCount = 0
		flashedDuringStep = []geometry.Coordinate{}
	}

	return flashedCount
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
		answer:   FindSolutionForInput2("example-input.dat"),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("puzzle-input.dat"),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput2("puzzle-input.dat"),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput(filename string) []string {
	return support.ReadFileIntoLines(filename)
}
