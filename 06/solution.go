package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
	Solution implementation
*/

//func FindSolutionForInput(filename string, targetDays int) int {
//	solution := 0
//	const NewFishAge = 8
//	ages := loadPuzzleInput(filename)
//
//	for days := 0; days < targetDays; days++ {
//		for index, age := range ages {
//			nextAge := age - 1
//			if nextAge < 0 {
//				ages[index] = 6
//				ages = append(ages, NewFishAge)
//			} else {
//				ages[index] = nextAge
//			}
//		}
//	}
//
//	solution = len(ages)
//
//	return solution
//}

func FindSolutionFastForInput(filename string, targetDays int) int {
	solution := 0

	ages := loadPuzzleInput(filename)

	var ageCounter = make(map[int]int)

	for _, age := range ages {
		ageCounter[age]++
	}

	for days := 0; days < targetDays; days++ {
		currentCount := ageCounter[0]

		for index := 0; index < 8; index++ {
			ageCounter[index] = ageCounter[index+1]
		}

		ageCounter[8] = currentCount
		ageCounter[6] += currentCount
	}

	for _, age := range ageCounter {
		solution += age
	}

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

	waitCount := 3
	var waitGroup sync.WaitGroup
	waitGroup.Add(waitCount)

	exampleChannel := make(chan Result)
	partOneChannel := make(chan Result)
	partTwoChannel := make(chan Result)

	go doExamples(exampleChannel, &waitGroup)
	go doPartOne(partOneChannel, &waitGroup)
	go doPartTwo(partTwoChannel, &waitGroup)

	exampleResult := <-exampleChannel
	partOneResult := <-partOneChannel
	partTwoResult := <-partTwoChannel

	waitGroup.Wait()

	log.
		Info().
		Int("example-answer", exampleResult.answer).
		Int64("example-duration", exampleResult.duration).
		Int("part-one-answer", partOneResult.answer).
		Int64("part-one-duration", partOneResult.duration).
		Int("part-two-answer", partTwoResult.answer).
		Int64("part-two-duration", partTwoResult.duration).
		Msg("day 04")
}

/*
	Executors
*/

func doExamples(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionFastForInput("example-input.dat", 80),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionFastForInput("puzzle-input.dat", 80),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionFastForInput("puzzle-input.dat", 256),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput(filename string) []int {
	var ages []int
	lines := strings.Split(support.ReadFile(filename), ",")
	for _, value := range lines {
		age, _ := strconv.Atoi(value)
		ages = append(ages, age)
	}

	return ages
}
