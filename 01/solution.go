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
		Msg("day ...")
}

func doExamples(waitGroup *sync.WaitGroup) {

	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	depthMeasurements := loadPuzzleInput()

	depthMeasurementCount := len(depthMeasurements)

	depthMeasurementIncreaseCount := 0

	for i := 1; i < depthMeasurementCount; i++ {
		if depthMeasurements[i-1] < depthMeasurements[i] {
			depthMeasurementIncreaseCount = depthMeasurementIncreaseCount + 1
		}
	}

	channel <- Result{
		answer:   depthMeasurementIncreaseCount,
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   1,
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
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
