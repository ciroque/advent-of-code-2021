package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
)

/*
	Solution implementation
*/

type PolymerFormulator struct {
	insertionRules map[string]string
	template       string
}

func NewPolymerFormulator(data []string) PolymerFormulator {
	polymerFormulator := PolymerFormulator{
		insertionRules: map[string]string{},
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

func (pf *PolymerFormulator) RunSubstitutions(count int) int {
	pairs := pf.GetPairs()

	for i := 0; i < count; i++ {
		nextPairs := map[string]int{}

		for pair := range pairs {
			insertionMonomer := pf.insertionRules[pair]
			firstMonomer := string(pair[0]) + insertionMonomer
			secondMonomer := insertionMonomer + string(pair[1])

			nextPairs[firstMonomer] += pairs[pair]
			nextPairs[secondMonomer] += pairs[pair]
		}

		pairs = nextPairs
	}

	firstCounts := map[string]int{}
	secondCounts := map[string]int{}

	for pair := range pairs {
		firstMonomer := string(pair[0])
		secondMonomer := string(pair[1])

		firstCounts[firstMonomer] += pairs[pair]
		secondCounts[secondMonomer] += pairs[pair]
	}

	monomerCounts := map[string]int{}

	for monomer := range firstCounts {
		monomerCounts[monomer] = int(math.Max(float64(firstCounts[monomer]), float64(secondCounts[monomer])))
	}

	var counts []int

	for _, value := range monomerCounts {
		counts = append(counts, value)
	}

	sort.Ints(counts)

	return counts[len(counts)-1] - counts[0]
}

func (pf *PolymerFormulator) GetPairs() map[string]int {
	windowed := map[string]int{}
	for index := 0; index < len(pf.template)-1; index++ {
		windowed[pf.template[index:index+2]]++
	}
	return windowed
}

func FindSolutionForInput(filename string, count int) int {
	polymerFormulator := NewPolymerFormulator(loadPuzzleInput(filename))

	return polymerFormulator.RunSubstitutions(count)
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
		Msg("Solved!")
}

/*
	Executors
*/

func doExampleOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("example-input.dat", 10),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doExampleTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("example-input.dat", 40),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("puzzle-input.dat", 10),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("puzzle-input.dat", 40),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput(filename string) []string {
	return support.ReadFileIntoLines(filename)
}
