package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
	"sync"
	"time"
)

/*
	Solution implementation
*/

const (
	One   int = 2
	Four      = 4
	Seven     = 3
	Eight     = 7
)

type DisplayLine struct {
	digits   []string
	patterns []string

	digitLengthMap map[int]int
}

func (dl *DisplayLine) MapDigitLengths() *DisplayLine {
	for _, digit := range dl.digits {
		dl.digitLengthMap[len(digit)]++
	}

	return dl
}

func (dl *DisplayLine) CalculateUniqueSegmentCount() int {
	return dl.digitLengthMap[One] +
		dl.digitLengthMap[Four] +
		dl.digitLengthMap[Seven] +
		dl.digitLengthMap[Eight]
}

func BuildDisplayLine(data string) DisplayLine {
	parts := strings.Split(data, "|")
	patterns := strings.Fields(parts[0])
	digits := strings.Fields(parts[1])

	return NewDisplayLine(digits, patterns)
}

func NewDisplayLine(digits []string, patterns []string) DisplayLine {
	return DisplayLine{
		digits:         digits,
		patterns:       patterns,
		digitLengthMap: make(map[int]int),
	}
}

func FindSolutionForInput(filename string, solutionCalculator func([]string) int) int {
	puzzleInput := loadPuzzleInput(filename)
	return solutionCalculator(puzzleInput)
}

func FindUniqueSegmentCount(puzzleInput []string) int {
	accumulator := 0
	for _, line := range puzzleInput {
		displayLine := BuildDisplayLine(line)
		accumulator += displayLine.MapDigitLengths().CalculateUniqueSegmentCount()
	}

	return accumulator
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
		answer:   FindSolutionForInput("example-input.dat", FindUniqueSegmentCount),
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
		answer:   FindSolutionForInput("puzzle-input.dat", FindUniqueSegmentCount),
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
