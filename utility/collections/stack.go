package collections

import (
	"advent-of-code-2021/utility/geometry"
	"errors"
)

// If you're looking for a concurrent Stack, this is not it.

type Stack struct {
	headIndex int
	items     []interface{}
}

func NewStack() Stack {
	return Stack{
		headIndex: -1,
		items:     []interface{}{},
	}
}

func (s *Stack) Clear() {
	s.items = []interface{}{}
}

func (s *Stack) IsEmpty() bool {
	return s.headIndex < 0
}

func (s *Stack) Peek() (interface{}, error) {
	if s.IsEmpty() {
		return nil, errors.New("stack is empty")
	}

	return s.items[s.headIndex], nil
}

func (s *Stack) Pop() (interface{}, error) {
	if s.IsEmpty() {
		return nil, errors.New("stack is empty")
	}

	item := s.items[s.headIndex]
	s.items = s.items[:s.headIndex]

	s.headIndex = len(s.items) - 1

	return item, nil
}

func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
	s.headIndex = len(s.items) - 1
}

func (s *Stack) ToSlice() []interface{} {
	var slice []interface{}
	for !s.IsEmpty() {
		c, _ := s.Pop()
		slice = append(slice, c.(geometry.Coordinate))
	}
	return slice
}
