package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math"
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
		Msg("day 02")
}

func findMostCommonValueFor(input []string, index int) byte {
	var occurrences = make(map[byte]int)

	occurrences['0'] = 0
	occurrences['1'] = 0

	for _, number := range input {
		occurrences[number[index]]++
	}

	if occurrences['0'] < occurrences['1'] {
		return '1'
	} else {
		return '0'
	}
}

func findMostCommonValueFor2(input []string, index int) byte {
	var occurrences = make(map[byte]int)

	occurrences['0'] = 0
	occurrences['1'] = 0

	for _, number := range input {
		occurrences[number[index]]++
	}

	if occurrences['0'] == occurrences['1'] {
		return '1'
	} else if occurrences['0'] < occurrences['1'] {
		return '1'
	} else {
		return '0'
	}
}

func findLeastCommonValueFor(input []string, index int) byte {
	var occurrences = make(map[byte]int)

	occurrences['0'] = 0
	occurrences['1'] = 0

	for _, number := range input {
		occurrences[number[index]]++
	}

	if occurrences['0'] == occurrences['1'] {
		return '0'
	} else if occurrences['0'] < occurrences['1'] {
		return '0'
	} else {
		return '1'
	}
}

func filterByValueAtIndex(input []string, index int, value byte) []string {
	var filtered []string

	for _, number := range input {
		if number[index] == value {
			filtered = append(filtered, number)
		}
	}

	return filtered
}

func findOxygenGeneratorRating(input []string) string {
	width := len(input[0])
	for index := 0; index < width; index++ {
		mostCommonValue := findMostCommonValueFor2(input, index)
		input = filterByValueAtIndex(input, index, mostCommonValue)
		if len(input) == 1 {
			return input[0]
		}
	}

	return ""
}

func findC02ScrubberRating(input []string) string {
	width := len(input[0])
	for index := 0; index < width; index++ {
		leastCommonValue := findLeastCommonValueFor(input, index)
		input = filterByValueAtIndex(input, index, leastCommonValue)
		if len(input) == 1 {
			return input[0]
		}
	}

	return ""
}

func calculateGamma(input []string) int {
	width := len(input[0])
	power := width
	gamma := 0

	for index := 0; index < width; index++ {
		power--

		var occurrences = make(map[byte]int)

		occurrences['0'] = 0
		occurrences['1'] = 0

		for _, number := range input {
			occurrences[number[index]]++
		}

		if findMostCommonValueFor(input, index) == '1' {
			additive := int(math.Pow(2, float64(power)))
			gamma = gamma + additive
		}
	}

	return gamma
}

func calculateEpsilon(input []string) int {
	width := len(input[0])
	power := width
	epsilon := 0

	for index := 0; index < width; index++ {
		power--

		var occurrences = make(map[byte]int)

		occurrences['0'] = 0
		occurrences['1'] = 0

		for _, number := range input {
			occurrences[number[index]]++
		}

		if findMostCommonValueFor(input, index) == '0' {
			additive := int(math.Pow(2, float64(power)))
			epsilon = epsilon + additive
		}
	}

	return epsilon
}

func doExamples(waitGroup *sync.WaitGroup) {
	exampleData := support.ReadFileIntoLines("example-input.dat")

	gamma := calculateGamma(exampleData)
	epsilon := calculateEpsilon(exampleData)

	oxygenGeneratorRating := findOxygenGeneratorRating(exampleData)
	c02ScrubberRating := findC02ScrubberRating(exampleData)

	oxy := toDecimal(oxygenGeneratorRating)
	c02 := toDecimal(c02ScrubberRating)

	log.
		Info().
		Int("gamma", gamma).
		Int("epsilon", epsilon).
		Int("part-one-answer", gamma*epsilon).
		Str("oxygenGeneratorRating", oxygenGeneratorRating).
		Str("c02ScrubberRating", c02ScrubberRating).
		Int("part-two-answer", oxy*c02).
		Msg("Example Data")

	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	diagnosticReport := loadPuzzleInput()
	gamma := calculateGamma(diagnosticReport)
	epsilon := calculateEpsilon(diagnosticReport)

	channel <- Result{
		answer:   gamma * epsilon,
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func toDecimal(binary string) int {
	power := len(binary)
	decimal := 0

	for _, value := range binary {
		power--
		if value == '1' {
			decimal += int(math.Pow(2, float64(power)))
		}
	}

	return decimal
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	diagnosticReport := loadPuzzleInput()
	oxygenGeneratorRating := findOxygenGeneratorRating(diagnosticReport)
	c02ScrubberRating := findC02ScrubberRating(diagnosticReport)

	oxygen := toDecimal(oxygenGeneratorRating)
	co2 := toDecimal(c02ScrubberRating)

	channel <- Result{
		answer:   oxygen * co2,
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput() []string {
	filename := "puzzle-input.dat"
	return support.ReadFileIntoLines(filename)
}
