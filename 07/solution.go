package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
	Solution implementation
*/

func CalculateLinearFuelConsumption(count int) int {
	return count
}

func CalculateTriangularFuelConsumption(count int) int {
	return ((count + 1) * count) / 2
}

func FindSolutionForInput(filename string, fuelConsumptionCalculation func(int) int) int {
	puzzleInput := loadPuzzleInput(filename)
	sort.Ints(puzzleInput)

	min := puzzleInput[0]
	max := puzzleInput[len(puzzleInput)-1]

	differences := []int{}
	for position := min; position < max; position++ {
		accumulator := 0
		for _, value := range puzzleInput {
			accumulator += fuelConsumptionCalculation(int(math.Abs(float64(value - position))))
		}
		differences = append(differences, accumulator)
	}

	sort.Ints(differences)

	return differences[0]
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
		answer:   FindSolutionForInput("example-input.dat", CalculateLinearFuelConsumption),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doExampleTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("example-input.dat", CalculateTriangularFuelConsumption),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("puzzle-input.dat", CalculateLinearFuelConsumption),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("puzzle-input.dat", CalculateTriangularFuelConsumption),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput(filename string) []int {
	input := support.ReadFile(filename)
	split := strings.Split(input, ",")
	var numbers []int

	for _, num := range split {
		number, _ := strconv.Atoi(num)
		numbers = append(numbers, number)
	}

	return numbers
}
