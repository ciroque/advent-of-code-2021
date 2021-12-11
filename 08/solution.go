package main

import (
	"advent-of-code-2021/utility"
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math"
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

	digitLengthMap   map[int]int
	patternLengthMap map[int][]string
	patternMap       map[string]int
}

func (dl *DisplayLine) MapDigitLengths() *DisplayLine {
	for _, digit := range dl.digits {
		dl.digitLengthMap[len(digit)]++
	}

	return dl
}

func (dl *DisplayLine) MapPatterns() *DisplayLine {
	lastRemainingFor := func(length int) string {
		patterns := dl.patternLengthMap[length]
		for _, pattern := range patterns {
			if len(pattern) > 0 {
				return pattern
			}
		}

		return ""
	}

	findBy := func(length int, contains string) string {
		patterns := dl.patternLengthMap[length]
		for index, pattern := range patterns {
			if utility.ContainsAllCharacters(pattern, contains) {
				dl.patternLengthMap[length][index] = ""
				return pattern
			}
		}
		return ""
	}

	removeChars := func(str, toRemove string) string {
		for _, char := range toRemove {
			str = strings.ReplaceAll(str, string(char), "")
		}
		return str
	}

	one := utility.SortString(dl.patternLengthMap[One][0])
	four := utility.SortString(dl.patternLengthMap[Four][0])
	seven := utility.SortString(dl.patternLengthMap[Seven][0])
	eight := utility.SortString(dl.patternLengthMap[Eight][0])
	three := utility.SortString(findBy(5, one))
	nine := utility.SortString(findBy(6, three))
	zero := utility.SortString(findBy(6, one))
	six := utility.SortString(lastRemainingFor(6))
	five := utility.SortString(findBy(5, removeChars(nine, one)))
	two := utility.SortString(lastRemainingFor(5))

	dl.patternMap[zero] = 0
	dl.patternMap[one] = 1
	dl.patternMap[two] = 2
	dl.patternMap[three] = 3
	dl.patternMap[four] = 4
	dl.patternMap[five] = 5
	dl.patternMap[six] = 6
	dl.patternMap[seven] = 7
	dl.patternMap[eight] = 8
	dl.patternMap[nine] = 9

	return dl
}

func (dl *DisplayLine) MapPatternLengths() *DisplayLine {
	for _, pattern := range dl.patterns {
		dl.patternLengthMap[len(pattern)] = append(dl.patternLengthMap[len(pattern)], pattern)
	}

	return dl
}

func (dl *DisplayLine) CalculateUniqueSegmentCount() int {
	return dl.digitLengthMap[One] +
		dl.digitLengthMap[Four] +
		dl.digitLengthMap[Seven] +
		dl.digitLengthMap[Eight]
}

func (dl *DisplayLine) CalculateOutputSum() int {
	accumulator := 0
	pow := 3
	for _, digit := range dl.digits {
		digit = utility.SortString(digit)
		d := dl.patternMap[digit]
		intermediate := d * int(math.Pow10(pow))
		accumulator += intermediate
		pow--
	}
	return accumulator
}

func BuildDisplayLine(data string) DisplayLine {
	parts := strings.Split(data, "|")
	patterns := strings.Fields(parts[0])
	digits := strings.Fields(parts[1])

	return NewDisplayLine(digits, patterns)
}

func NewDisplayLine(digits []string, patterns []string) DisplayLine {
	return DisplayLine{
		digits:           digits,
		patterns:         patterns,
		digitLengthMap:   make(map[int]int),
		patternLengthMap: make(map[int][]string),
		patternMap:       make(map[string]int),
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

func FindOutputValuesSum(puzzleInput []string) int {
	accumulator := 0
	for _, line := range puzzleInput {
		displayLine := BuildDisplayLine(line)
		accumulator += displayLine.MapPatternLengths().MapPatterns().CalculateOutputSum()
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
		answer:   FindSolutionForInput("example-input.dat", FindOutputValuesSum),
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
		answer:   FindSolutionForInput("puzzle-input.dat", FindOutputValuesSum),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput(filename string) []string {
	return support.ReadFileIntoLines(filename)
}
