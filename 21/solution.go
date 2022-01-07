package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

/*
	Solution implementation
*/

type DeterministicDieRoller struct {
	rolls int
}

func NewDie(seed int) DeterministicDieRoller {
	return DeterministicDieRoller{rolls: seed}
}

func (d *DeterministicDieRoller) Roll() int {
	d.rolls++
	return (d.rolls-1)%100 + 1
}

func (d *DeterministicDieRoller) RollThree() (int, int, int) {
	return d.Roll(), d.Roll(), d.Roll()
}

func (d *DeterministicDieRoller) SummedThreeRolls() int {
	first, second, third := d.RollThree()
	return first + second + third
}

type Player struct {
	name  string
	space int
	score int
}

func NewPlayer(name string, start int) Player {
	return Player{
		name:  name,
		space: start,
	}
}

func (p *Player) CalculateScore(move int) {
	p.space += move
	p.space = (p.space-1)%10 + 1
	p.score += p.space
}

func (p *Player) Won() bool {
	return p.score >= 1000
}

func FindSolutionForInput(playerOneStart int, playerTwoStart int) int {
	noWinner := true
	playerIndex := 0
	players := []Player{
		NewPlayer("Player One", playerOneStart),
		NewPlayer("Player Two", playerTwoStart),
	}
	roller := NewDie(0)

	for noWinner {
		players[playerIndex].CalculateScore(roller.SummedThreeRolls())
		noWinner = !players[playerIndex].Won()
		playerIndex ^= 1
	}

	return roller.rolls * players[playerIndex].score
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
		answer:   FindSolutionForInput(4, 8),
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
		answer:   FindSolutionForInput(2, 10),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   0, //  FindSolutionForInput("puzzle-input.dat"),
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}
