package main

import (
	"fmt"
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
	"sync"
	"time"
)

/*
	Solution implementation
*/

type VertexInfo struct {
	label string
}

func NewVertexInfo(label string) VertexInfo {
	return VertexInfo{
		label: label,
	}
}

func (vi *VertexInfo) IsEnd() bool {
	return vi.label == "end"
}

func (vi *VertexInfo) IsBig() bool {
	return vi.label == strings.ToUpper(vi.label)
}

type AdjacencyList struct {
	adjacent map[VertexInfo][]VertexInfo
}

func NewAdjacencyList(puzzleInput []string) AdjacencyList {
	adjacencyList := AdjacencyList{
		adjacent: make(map[VertexInfo][]VertexInfo),
	}

	for _, line := range puzzleInput {
		parts := strings.Split(line, "-")
		from := NewVertexInfo(parts[0])
		to := NewVertexInfo(parts[1])
		adjacencyList.adjacent[from] = append(adjacencyList.adjacent[from], to)
		adjacencyList.adjacent[to] = append(adjacencyList.adjacent[to], from)
	}

	return adjacencyList
}

func (al *AdjacencyList) Traverse(handlePathFound func(path []VertexInfo)) int {
	var traverseGraph func(vertex VertexInfo, path []VertexInfo, visited map[VertexInfo]int, pathCount int) int
	traverseGraph = func(vertex VertexInfo, path []VertexInfo, visited map[VertexInfo]int, pathCount int) int {
		path = append(path, vertex)
		for _, destination := range al.adjacent[vertex] {
			if destination.IsEnd() {
				pathCount++
				handlePathFound(path)
				continue
			}

			visited[vertex]++

			if count, _ := visited[destination]; count == 0 || destination.IsBig() {
				pathCount = traverseGraph(destination, path, visited, pathCount)
			}

			visited[vertex]--
		}

		return pathCount
	}

	return traverseGraph(VertexInfo{label: "start"}, []VertexInfo{}, make(map[VertexInfo]int), 0)
}

func PrintPaths(paths [][]VertexInfo) {
	for _, path := range paths {
		vertexCount := len(path) - 1
		for index, vertex := range path {
			fmt.Printf("%v", vertex.label)
			if index < vertexCount {
				fmt.Print(" -> ")
			} else {
				fmt.Println()
			}
		}
	}
}

func FindSolutionForInput(filename string) int {
	puzzleInput := loadPuzzleInput(filename)
	adjacencyList := NewAdjacencyList(puzzleInput)
	var paths [][]VertexInfo
	pathCount := 0
	solution := adjacencyList.Traverse(func(path []VertexInfo) {
		pathCount++
		paths = append(paths, path)
	})

	//PrintPaths(paths)

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
		answer:   0, // FindSolutionForInput("example-input.dat"),
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
		answer:   0, // FindSolutionForInput("puzzle-input.dat"),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput(filename string) []string {
	return support.ReadFileIntoLines(filename)
}
