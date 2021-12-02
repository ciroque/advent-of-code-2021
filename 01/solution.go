package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strconv"
	"sync"
	"time"
)

type Result struct {
	answer   int
	duration int64
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	_ = loadPuzzleInput()

	waitCount := 3
	var waitGroup sync.WaitGroup
	waitGroup.Add(waitCount)

	partOneChannel := make(chan Result)
	partTwoChannel := make(chan Result)

	go doExamples(&waitGroup)
	go doPartOne(partOneChannel, &waitGroup)
	go doPartTwo(partTwoChannel, &waitGroup)

	partOneResult := <-partOneChannel
	partTwoResult := <-partTwoChannel

	waitGroup.Wait()

	log.
		Info().
		Int("part-one-answer", partOneResult.answer).
		Int64("part-one-duration", partOneResult.duration).
		Int("part-two-answer", partTwoResult.answer).
		Int64("part-two-duration", partTwoResult.duration).
		Msg("day 1")
}

func doExamples(waitGroup *sync.WaitGroup) {

	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	depthMeasurements := loadPuzzleInput()

	depthMeasurementIncreaseCount := countDepthIncreases(depthMeasurements)

	channel <- Result{
		answer:   depthMeasurementIncreaseCount,
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	WindowSize := 3
	start := time.Now()

	depthMeasurements := loadPuzzleInput()

	windowedDepthMeasurements := groupIntoWindows(depthMeasurements, WindowSize)

	summedWindows := sumWindowedDepthMeasurements(windowedDepthMeasurements)

	depthMeasurementIncreaseCount := countDepthIncreases(summedWindows)

	channel <- Result{
		answer: depthMeasurementIncreaseCount,
		duration: time.
			Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func countDepthIncreases(depthMeasurements []int) int {
	depthMeasurementIncreaseCount := 0
	depthMeasurementCount := len(depthMeasurements)
	for i := 1; i < depthMeasurementCount; i++ {
		if depthMeasurements[i-1] < depthMeasurements[i] {
			depthMeasurementIncreaseCount = depthMeasurementIncreaseCount + 1
		}
	}

	return depthMeasurementIncreaseCount
}

func sumWindowedDepthMeasurements(measurements [][]int) []int {
	var sums []int
	for _, group := range measurements {
		sum := 0

		for _, measurement := range group {
			sum = sum + measurement
		}

		sums = append(sums, sum)
	}

	return sums
}

func groupIntoWindows(measurements []int, windowSize int) [][]int {
	var windowed [][]int
	measurementCount := len(measurements)
	for i := 0; i < measurementCount-windowSize+1; i++ {
		lastIndex := i + windowSize
		windowed = append(windowed, measurements[i:lastIndex])
	}

	return windowed
}

func loadPuzzleInput() []int {
	filename := "puzzle-input.dat"
	strings := support.ReadFileIntoLines(filename)
	var numbers []int
	for _, value := range strings {
		number, _ := strconv.Atoi(value)
		numbers = append(numbers, number)
	}

	return numbers
}
