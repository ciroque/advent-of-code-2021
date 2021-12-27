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

func (vi *VertexInfo) IsBig() bool {
	return vi.label == strings.ToUpper(vi.label)
}

func (vi *VertexInfo) IsEnd() bool {
	return vi.label == "end"
}

func (vi *VertexInfo) IsSmall() bool {
	return vi.label == strings.ToLower(vi.label)
}

func (vi *VertexInfo) IsStart() bool {
	return vi.label == "start"
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

func (al *AdjacencyList) Traverse(
	handlePathFound func(path []VertexInfo),
	navigateNext func(destination VertexInfo, visited map[VertexInfo]int, exceededSmallCaveVisits bool) bool) int {

	multipleVisits := false

	var traverseGraph func(current VertexInfo, path []VertexInfo, visited map[VertexInfo]int, pathCount int, exceededSmallCaveVisits bool) int
	traverseGraph = func(current VertexInfo, path []VertexInfo, visited map[VertexInfo]int, pathCount int, exceededSmallCaveVisits bool) int {
		path = append(path, current)
		for _, destination := range al.adjacent[current] {
			if destination.IsEnd() {
				//path = append(path, destination)
				pathCount++
				handlePathFound(path)
				continue
			}

			multipleVisits = multipleVisits || (current.IsSmall() && visited[current] >= 1)

			visited[current]++

			if navigateNext(destination, visited, exceededSmallCaveVisits) {
				pathCount = traverseGraph(destination, path, visited, pathCount, multipleVisits)
			}

			visited[current]--
		}

		return pathCount
	}

	startingPathCount := 0
	return traverseGraph(VertexInfo{label: "start"}, []VertexInfo{}, make(map[VertexInfo]int), startingPathCount, multipleVisits)
}

func PrintPaths(paths [][]VertexInfo) {
	for _, path := range paths {
		PrintPath(path)
	}
	fmt.Println()
}

func PrintPath(path []VertexInfo) {
	vertexCount := len(path) - 1
	for index, vertex := range path {
		fmt.Printf("%v", vertex.label)
		if index < vertexCount {
			fmt.Print(",")
		} else {
			fmt.Println()
		}
	}
}

func VisitSmallCavesOnlyOnce(destination VertexInfo, visited map[VertexInfo]int, _ bool) bool {
	count := visited[destination]
	return count == 0 || destination.IsBig()
}

func ExtendedSearch(destination VertexInfo, visited map[VertexInfo]int, multipleVisits bool) bool {
	return (destination.IsBig() ||
		(visited[destination] < 2 && !multipleVisits) ||
		visited[destination] == 0) && !destination.IsStart()
}

func FindSolutionForInput(filename string, navigateNext func(destination VertexInfo, visited map[VertexInfo]int, multipleVisitsToSmall bool) bool) int {
	var paths [][]VertexInfo
	trackPaths := func(path []VertexInfo) { paths = append(paths, path) }
	puzzleInput := loadPuzzleInput(filename)
	adjacencyList := NewAdjacencyList(puzzleInput)

	solution := adjacencyList.Traverse(trackPaths, navigateNext)
	//PrintPaths(paths)

	//toString := func(infos []VertexInfo) string {
	//	var labels []string
	//	for _, info := range infos {
	//		labels = append(labels, info.label)
	//	}
	//	labels = append(labels, "end")
	//	return strings.Join(labels, ",")
	//}
	//
	//expectedPaths := ExpectedPaths()
	//for _, path := range paths {
	//	expectedPaths[toString(path)]++
	//}
	//
	//fmt.Println("Missing paths:")
	//for path, count := range expectedPaths {
	//	if count == 0 {
	//		fmt.Println(path)
	//	}
	//}

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
		answer:   FindSolutionForInput("example-input.dat", VisitSmallCavesOnlyOnce),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doExampleTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("example-input.dat", ExtendedSearch),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("puzzle-input.dat", VisitSmallCavesOnlyOnce),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   FindSolutionForInput("puzzle-input.dat", ExtendedSearch),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput(filename string) []string {
	return support.ReadFileIntoLines(filename)
}

func ExpectedPaths() map[string]int {
	return map[string]int{
		"start,A,b,A,b,A,c,A,end": 0,
		"start,A,b,A,b,A,end":     0,
		"start,A,b,A,b,end":       0,
		"start,A,b,A,c,A,b,A,end": 0,
		"start,A,b,A,c,A,b,end":   0,
		"start,A,b,A,c,A,c,A,end": 0,
		"start,A,b,A,c,A,end":     0,
		"start,A,b,A,end":         0,
		"start,A,b,d,b,A,c,A,end": 0,
		"start,A,b,d,b,A,end":     0,
		"start,A,b,d,b,end":       0,
		"start,A,b,end":           0,
		"start,A,c,A,b,A,b,A,end": 0,
		"start,A,c,A,b,A,b,end":   0,
		"start,A,c,A,b,A,c,A,end": 0,
		"start,A,c,A,b,A,end":     0,
		"start,A,c,A,b,d,b,A,end": 0,
		"start,A,c,A,b,d,b,end":   0,
		"start,A,c,A,b,end":       0,
		"start,A,c,A,c,A,b,A,end": 0,
		"start,A,c,A,c,A,b,end":   0,
		"start,A,c,A,c,A,end":     0,
		"start,A,c,A,end":         0,
		"start,A,end":             0,
		"start,b,A,b,A,c,A,end":   0,
		"start,b,A,b,A,end":       0,
	}
}
