package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
	Solution implementation
*/

type Point struct {
	x int
	y int
}

type Line struct {
	start              Point
	end                Point
	intermediatePoints []Point
	orientation        string
}

func (l *Line) CalculateMovementSteps() (int, int, int) {
	deltaX := l.end.x - l.start.x
	deltaY := l.end.y - l.start.y
	absDeltaX := int(math.Abs(float64(deltaX)))
	absDeltaY := int(math.Abs(float64(deltaY)))

	xm, ym, count := 0, 0, 0

	if deltaX == 0 {
		xm = deltaX
		ym = deltaY / absDeltaY
		count = absDeltaY
	} else if deltaY == 0 {
		xm = deltaX / absDeltaX
		ym = deltaY
		count = absDeltaX
	} else if deltaX > deltaY {
		xm = deltaX / absDeltaX
		ym = deltaY / absDeltaX
		count = absDeltaX
	} else {
		xm = deltaX / absDeltaY
		ym = deltaY / absDeltaY
		count = absDeltaY
	}

	return xm, ym, count
}

func (l *Line) IsDiagonal() bool {
	return !l.IsVertical() && !l.IsHorizontal()
}

func (l *Line) IsHorizontal() bool {
	return l.start.y == l.end.y
}

func (l *Line) IsVertical() bool {
	return l.start.x == l.end.x
}

func (l *Line) ParseLine(input string) {
	parsePoint := func(points string) Point {
		coordinates := strings.Split(points, ",")
		x, _ := strconv.Atoi(coordinates[0])
		y, _ := strconv.Atoi(coordinates[1])

		return Point{x: x, y: y}
	}

	points := strings.Split(input, " -> ")

	l.start = parsePoint(points[0])
	l.end = parsePoint(points[1])
}

func Parse(input []string) []Line {
	var lines []Line

	for _, inputLine := range input {
		line := Line{}
		line.ParseLine(inputLine)
		line.PopulateIntermediatePoints()
		lines = append(lines, line)
	}

	return lines
}

func (l *Line) Points() []Point {
	var points []Point
	points = append(points, l.start)
	points = append(points, l.intermediatePoints...)
	points = append(points, l.end)
	return points
}

func (l *Line) PopulateIntermediatePoints() {
	xd, yd, count := l.CalculateMovementSteps()

	x := l.start.x
	y := l.start.y

	for n := 1; n < count; n++ {
		x += xd
		y += yd
		l.intermediatePoints = append(l.intermediatePoints, Point{x: x, y: y})
	}
}

func FindSolutionForInput(filename string, includeDiagonals bool) int {
	puzzleInput := loadPuzzleInput(filename)
	lines := Parse(puzzleInput)

	var points = make(map[Point]int)

	for _, line := range lines {
		if includeDiagonals || !line.IsDiagonal() {
			for _, point := range line.Points() {
				points[point]++
			}
		}
	}

	solution := 0
	for _, point := range points {
		if point >= 2 {
			solution++
		}
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
		answer:   FindSolutionForInput("example-input.dat", true),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("puzzle-input.dat", false),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("puzzle-input.dat", true),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput(filename string) []string {
	return support.ReadFileIntoLines(filename)
}
