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

func (l *Line) MovementSteps() (int, int, int) {
	deltaX := l.end.x - l.start.x
	deltaY := l.end.y - l.start.y
	absDeltaX := int(math.Abs(float64(deltaX)))
	absDeltaY := int(math.Abs(float64(deltaY)))

	xm, ym, count := 0, 0, 0

	if deltaX > deltaY {
		xm = int(deltaX / absDeltaX)
		ym = int(deltaY / absDeltaX)
		count = absDeltaX
	} else {
		xm = int(deltaX / absDeltaY)
		ym = int(deltaY / absDeltaY)
		count = absDeltaY
	}

	return xm, ym, count
}

type Line struct {
	start              Point
	end                Point
	intermediatePoints []Point
	orientation        string
}

func (l *Line) CalculatePoints() {
	if l.IsHorizontal() {
		l.orientation = "Horizontal"
		for x2 := l.start.x + 1; x2 < l.end.x; x2++ {
			l.intermediatePoints = append(l.intermediatePoints, Point{x: x2, y: l.start.y})
		}

	} else if l.IsVertical() {
		l.orientation = "Vertical"
		for y2 := l.start.y + 1; y2 < l.end.y; y2++ {
			l.intermediatePoints = append(l.intermediatePoints, Point{x: l.start.x, y: y2})
		}

	} else if l.IsDiagonal() {
		l.orientation = "Diagonal"

		xd, yd, count := l.MovementSteps()

		x := l.start.x
		y := l.start.y

		for n := 1; n < count; n++ {
			x += xd
			y += yd
			l.intermediatePoints = append(l.intermediatePoints, Point{x: x, y: y})
		}
	}
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

func (l *Line) NormalizePoints() {
	if l.IsHorizontal() {
		if l.start.x > l.end.x {
			l.Swap()
		}
	} else if l.IsVertical() {
		if l.start.y > l.end.y {
			l.Swap()
		}
	} else if l.IsDiagonal() {
		//d := int(math.Pow(float64(l.start.x + l.end.x), 2) + math.Pow(float64(l.start.y + l.end.y), 2))
		//if d > 0 {
		//	l.Swap()
		//}
	}
}

func (l *Line) Swap() {
	swap := l.end
	l.end = l.start
	l.start = swap
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
		line.NormalizePoints()
		line.CalculatePoints()
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
