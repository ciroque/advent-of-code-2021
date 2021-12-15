package main

import (
	"advent-of-code-2021/utility/collections"
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sort"
	"sync"
	"time"
)

/*
	Solution implementation
*/

func InitializeInvalidScores() map[int32]int {
	return map[int32]int{
		')': 3,
		']': 57,
		'}': 1197,
		'>': 25137,
	}
}

func InitializeIncompleteScores() map[int32]int {
	return map[int32]int{
		')': 1,
		']': 2,
		'}': 3,
		'>': 4,
	}
}

func InitializeComplements() map[int32]int32 {
	return map[rune]rune{
		'}': '{',
		')': '(',
		']': '[',
		'>': '<',
	}
}

func CalculateAutocompleteScore(puzzleInput []string) int {
	complements := map[rune]rune{
		'{': '}',
		'(': ')',
		'[': ']',
		'<': '>',
	}

	incompleteScores := InitializeIncompleteScores()
	var scores []int

	for _, line := range puzzleInput {
		score := 0
		stack := collections.NewStack()
		valid := true
		for _, char := range line {
			if _, found := complements[char]; found {
				stack.Push(complements[char])
			} else if c, _ := stack.Peek(); char == c.(int32) {
				_, _ = stack.Pop()
			} else {
				valid = false
				break
			}
		}

		if valid {
			for !stack.IsEmpty() {
				char, _ := stack.Pop()
				score *= 5
				if v, found := incompleteScores[char.(int32)]; found {
					score += v
				}
			}
			scores = append(scores, score)
		}
	}

	sort.Ints(scores)

	return scores[len(scores)/2]
}

func CalculateTotalSyntaxErrorScore(puzzleInput []string) int {
	solution := 0
	complements := InitializeComplements()
	scores := InitializeInvalidScores()

	for _, line := range puzzleInput {
		score := 0
		stack := collections.NewStack()
		for _, char := range line {
			if complement, found := complements[char]; found {
				item, _ := stack.Pop()
				if item.(int32) != complement {
					score += scores[char]
					break
				}
			} else {
				stack.Push(char)
			}
		}
		solution += score
	}

	return solution
}

func FindSolutionForInput(filename string, calculateSolution func([]string) int) int {
	return calculateSolution(loadPuzzleInput(filename))
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
		answer:   FindSolutionForInput("example-input.dat", CalculateTotalSyntaxErrorScore),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doExampleTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("example-input.dat", CalculateAutocompleteScore),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("puzzle-input.dat", CalculateTotalSyntaxErrorScore),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("puzzle-input.dat", CalculateAutocompleteScore),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput(filename string) []string {
	return support.ReadFileIntoLines(filename)
}
