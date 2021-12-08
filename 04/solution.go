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

type Board struct {
	size  int
	cells []Cell
}

func NewBoard() Board {
	return Board{
		size:  5,
		cells: []Cell{},
	}
}

func (p *Board) appendNewCell(value int) {
	p.cells = append(p.cells, NewCell(value))
}

func (p *Board) markNumber(value int) bool {
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

func (p *Board) checkDimension(calculateIndex IndexCalculator) bool {
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

func (p *Board) checkForWin() bool {
	return p.checkDimension(columnIndexCalculator) || p.checkDimension(rowIndexCalculator)
}

func (p *Board) print() string {
	board := ""
	indicator := " "
	for index, cell := range p.cells {
		if index%p.size == 0 {
			board = board + "\n"
		}
		if cell.marked {
			indicator = "+"
		} else {
			indicator = "-"
		}
		board = board + fmt.Sprintf("%v%2v ", indicator, cell.value)
	}

	board = board + "\n"
	return board
}

func (p *Board) sumUnmarkedCells() int {
	sum := 0

	for _, cell := range p.cells {
		if cell.marked == false {
			sum = sum + cell.value
		}
	}

	return sum
}

func buildBoards(rawBoards [][]string) []Board {
	var boards []Board

	for _, rawBoard := range rawBoards {
		board := NewBoard()

		for _, row := range rawBoard {
			numbers := strings.Fields(row)

			for _, number := range numbers {
				value, _ := strconv.Atoi(number)
				board.appendNewCell(value)
			}
		}

		boards = append(boards, board)
	}

	return boards
}

func loadBoards(input string) []Board {
	lines := strings.Split(input, "\n")
	rawBoards := separateBoardInput(lines)
	return buildBoards(rawBoards)
}

func runGame(numbers []int, boards []Board) (bool, int, *Board) {
	for _, number := range numbers {
		for _, board := range boards {
			if board.markNumber(number) {
				if board.checkForWin() {
					return true, number, &board
				}
			}
		}
	}

	return false, -1, nil
}

func runGame2(numbers []int, boards []Board) (bool, int, *Board) {
	mostCallsToWinNumber := 0
	mostCallsToWinIndex := 0
	var mostCallsToWinBoard Board

	for _, board := range boards {
		for nidx, number := range numbers {
			if board.markNumber(number) {
				if board.checkForWin() {
					if mostCallsToWinIndex < nidx {
						mostCallsToWinIndex = nidx
						mostCallsToWinBoard = board
						mostCallsToWinNumber = number
					}
					break
				}
			}
		}
	}

	return true, mostCallsToWinNumber, &mostCallsToWinBoard
}

func separateBoardInput(lines []string) [][]string {
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

	inputNumbers := loadBoardInput("example-numbers.dat")
	drawnNumbers := stringToIntList(inputNumbers)

	inputBoards := loadBoardInput("example-input.dat")
	boards := loadBoards(inputBoards)

	winner, number, board := runGame(drawnNumbers, boards)
	if winner {
		fmt.Printf(board.print())
		boardScore := board.sumUnmarkedCells()
		solution = number * boardScore
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

	inputNumbers := loadBoardInput("puzzle-numbers.dat")
	drawnNumbers := stringToIntList(inputNumbers)

	inputBaords := loadBoardInput("puzzle-input.dat")
	boards := loadBoards(inputBaords)

	winner, number, board := runGame(drawnNumbers, boards)
	if winner {
		fmt.Printf(board.print())
		boardScore := board.sumUnmarkedCells()
		solution = number * boardScore
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

	inputNumbers := loadBoardInput("puzzle-numbers.dat")
	inputBoards := loadBoardInput("puzzle-input.dat")

	drawnNumbers := stringToIntList(inputNumbers)
	boards := loadBoards(inputBoards)

	winner, number, board := runGame2(drawnNumbers, boards)
	if winner {
		fmt.Printf(board.print())
		score := board.sumUnmarkedCells()
		solution = number * score
	}

	channel <- Result{
		answer:   solution,
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadBoardInput(filename string) string {
	return support.ReadFile(filename)
}
