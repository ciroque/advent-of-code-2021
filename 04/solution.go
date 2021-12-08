package main

import (
	"fmt"
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

type Cell struct {
	value  int
	marked bool
}

func NewCell(value int) Cell {
	return Cell{
		value:  value,
		marked: false,
	}
}

func (c *Cell) mark() {
	c.marked = true
}

type Puzzle struct {
	size  int
	cells []Cell
}

func NewPuzzle() Puzzle {
	return Puzzle{
		size:  5,
		cells: []Cell{},
	}
}

func (p *Puzzle) appendNewCell(value int) {
	p.cells = append(p.cells, NewCell(value))
}

func (p *Puzzle) markNumber(value int) bool {
	length := len(p.cells)
	for index := 0; index < length; index++ {
		if p.cells[index].value == value {
			p.cells[index].mark()
			return true
		}
	}

	return false
}

type IndexCalculator func(colIndex int, rowIndex int, size int) int

func columnIndexCalculator(colIndex int, rowIndex int, size int) int {
	return (colIndex % size) + (rowIndex * size)
}

func rowIndexCalculator(colIndex int, rowIndex int, size int) int {
	return (colIndex * size) + (rowIndex % size)
}

func (p *Puzzle) checkDimension(calculateIndex IndexCalculator) bool {
	size := p.size

	for cellIndex := 0; cellIndex < size; cellIndex++ {
		foundCount := 0
		for rowIndex := 0; rowIndex < size; rowIndex++ {
			index := calculateIndex(cellIndex, rowIndex, size)
			if p.cells[index].marked == true {
				foundCount = foundCount + 1
			}
		}

		if foundCount == size {
			return true
		}
	}

	return false
}

func (p *Puzzle) checkForWin() bool {
	return p.checkDimension(columnIndexCalculator) || p.checkDimension(rowIndexCalculator)
}

func (p *Puzzle) print() string {
	puzzle := ""
	indicator := " "
	for index, cell := range p.cells {
		if index%p.size == 0 {
			puzzle = puzzle + "\n"
		}
		if cell.marked {
			indicator = "+"
		} else {
			indicator = "-"
		}
		puzzle = puzzle + fmt.Sprintf("%v%2v ", indicator, cell.value)
	}

	puzzle = puzzle + "\n"
	return puzzle
}

func (p *Puzzle) sumUnmarkedCells() int {
	sum := 0

	for _, cell := range p.cells {
		if cell.marked == false {
			sum = sum + cell.value
		}
	}

	return sum
}

func buildPuzzles(rawPuzzles [][]string) []Puzzle {
	var puzzles []Puzzle

	for _, rawPuzzle := range rawPuzzles {
		puzzle := NewPuzzle()

		for _, puzzleRow := range rawPuzzle {
			numbers := strings.Fields(puzzleRow)

			for _, number := range numbers {
				value, _ := strconv.Atoi(number)
				puzzle.appendNewCell(value)
			}
		}

		puzzles = append(puzzles, puzzle)
	}

	return puzzles
}

func loadPuzzles(input string) []Puzzle {
	lines := strings.Split(input, "\n")
	rawPuzzles := separatePuzzleInput(lines)
	return buildPuzzles(rawPuzzles)
}

func runGame(numbers []int, puzzles []Puzzle) (bool, int, *Puzzle) {
	for _, number := range numbers {
		for _, puzzle := range puzzles {
			if puzzle.markNumber(number) {
				if puzzle.checkForWin() {
					return true, number, &puzzle
				}
			}
		}
	}

	return false, -1, nil
}

func runGame2(numbers []int, puzzles []Puzzle) (bool, int, *Puzzle) {
	mostCallsToWinNumber := 0
	mostCallsToWinIndex := 0
	var mostCallsToWinPuzzle Puzzle

	for _, puzzle := range puzzles {
		for nidx, number := range numbers {
			if puzzle.markNumber(number) {
				if puzzle.checkForWin() {
					if mostCallsToWinIndex < nidx {
						mostCallsToWinIndex = nidx
						mostCallsToWinPuzzle = puzzle
						mostCallsToWinNumber = number
					}
					break
				}
			}
		}
	}

	return true, mostCallsToWinNumber, &mostCallsToWinPuzzle
}

func separatePuzzleInput(lines []string) [][]string {
	var output [][]string

	startIndex := 0
	endIndex := 0

	for _, line := range lines {
		if line == "" {
			output = append(output, lines[startIndex:endIndex])
			startIndex = endIndex + 1
		}
		endIndex += 1
	}

	output = append(output, lines[startIndex:endIndex])

	return output
}

func stringToIntList(input string) []int {
	var integers []int

	splitString := strings.Split(input, ",")

	for _, s := range splitString {
		n, _ := strconv.Atoi(s)
		integers = append(integers, n)
	}

	return integers
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

	solution := 0

	inputNumbers := loadPuzzleInput("example-numbers.dat")
	drawnNumbers := stringToIntList(inputNumbers)

	inputPuzzles := loadPuzzleInput("example-input.dat")
	puzzles := loadPuzzles(inputPuzzles)

	winner, number, puzzle := runGame(drawnNumbers, puzzles)
	if winner {
		fmt.Printf(puzzle.print())
		puzzleScore := puzzle.sumUnmarkedCells()
		solution = number * puzzleScore
	}

	channel <- Result{
		answer:   solution,
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	solution := 0

	inputNumbers := loadPuzzleInput("puzzle-numbers.dat")
	drawnNumbers := stringToIntList(inputNumbers)

	inputPuzzles := loadPuzzleInput("puzzle-input.dat")
	puzzles := loadPuzzles(inputPuzzles)

	winner, number, puzzle := runGame(drawnNumbers, puzzles)
	if winner {
		fmt.Printf(puzzle.print())
		puzzleScore := puzzle.sumUnmarkedCells()
		solution = number * puzzleScore
	}

	channel <- Result{
		answer:   solution,
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	solution := 0

	inputNumbers := loadPuzzleInput("puzzle-numbers.dat")
	inputPuzzles := loadPuzzleInput("puzzle-input.dat")

	drawnNumbers := stringToIntList(inputNumbers)
	puzzles := loadPuzzles(inputPuzzles)

	winner, number, puzzle := runGame2(drawnNumbers, puzzles)
	if winner {
		fmt.Printf(puzzle.print())
		puzzleScore := puzzle.sumUnmarkedCells()
		solution = number * puzzleScore
	}

	channel <- Result{
		answer:   solution,
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput(filename string) string {
	return support.ReadFile(filename)
}
