package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sort"
	"strings"
	"sync"
	"time"
)

/*
	Solution implementation
*/

type PolymerFormulator struct {
	elementCounts  []map[string]int
	insertionRules map[string]string
	template       string
	polymers       []string
}

func NewPolymerFormulator(data []string) PolymerFormulator {
	polymerFormulator := PolymerFormulator{
		elementCounts:  []map[string]int{},
		insertionRules: map[string]string{},
		polymers:       []string{},
	}

	inInsertionRules := false

	for _, line := range data {
		if len(line) == 0 {
			inInsertionRules = true
			continue
		}
		if inInsertionRules {
			parts := strings.Split(line, " -> ")
			polymerFormulator.insertionRules[parts[0]] = parts[1]
		} else {
			polymerFormulator.template = line
		}
	}

	return polymerFormulator
}

func (pf *PolymerFormulator) GetPairs(polymer string) []string {
	var windowed []string
	for index := 0; index < len(polymer)-1; index++ {
		windowed = append(windowed, polymer[index:index+2])
	}
	return windowed
}

func (pf *PolymerFormulator) RunFirstInsertion() {
	pf.RunInsertion(pf.template)
}

func (pf *PolymerFormulator) RunInsertion(polymer string) {
	var newPolymer []string
	pairs := pf.GetPairs(polymer)

	for _, pair := range pairs {
		if insertion, found := pf.insertionRules[pair]; found {
			newPolymer = append(newPolymer, string(pair[0]), insertion)
		} else {
			newPolymer = append(newPolymer, string(pair[0]), string(pair[1]))
		}
	}

	newPolymer = append(newPolymer, string(pairs[len(pairs)-1][1]))

	pf.polymers = append(pf.polymers, strings.Join(newPolymer, ""))
}

func (pf *PolymerFormulator) CalculateSolution() int {
	counts := []int{}
	index := len(pf.elementCounts) - 1
	for _, value := range pf.elementCounts[index] {
		counts = append(counts, value)
	}

	sort.Ints(counts)

	return counts[len(counts)-1] - counts[0]
}

func (pf *PolymerFormulator) CalculateElementFrequencies() {
	index := len(pf.polymers) - 1
	elementCounts := map[string]int{}
	for _, element := range pf.polymers[index] {
		elementCounts[string(element)]++
	}
	pf.elementCounts = append(pf.elementCounts, elementCounts)
}

func FindSolutionForInput(filename string) int {
	polymerFormulator := NewPolymerFormulator(loadPuzzleInput(filename))

	polymerFormulator.RunFirstInsertion()
	polymerFormulator.CalculateElementFrequencies()

	for index := 0; index < 9; index++ {
		polymerFormulator.RunInsertion(polymerFormulator.polymers[index])
		polymerFormulator.CalculateElementFrequencies()
	}

	return polymerFormulator.CalculateSolution()
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
		answer:   0, //  FindSolutionForInput("example-input.dat"),
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
		answer:   0, //  FindSolutionForInput("puzzle-input.dat"),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput(filename string) []string {
	return support.ReadFileIntoLines(filename)
}
