package collections

import "testing"

func TestStackIsEmpty(t *testing.T) {
	stack := NewStack()

	if !stack.IsEmpty() {
		t.Log("New Stack should be empty")
		t.Fail()
	}
}

func TestStack_Push(t *testing.T) {
	stack := NewStack()
	stack.Push(0)

	if stack.IsEmpty() {
		t.Log("Stack should not be empty after Push")
		t.Fail()
	}
}

func TestStack_Pop(t *testing.T) {
	stack := NewStack()

	stack.Push(0)
	value, _ := stack.Pop()
	if value.(int) != 0 {
		t.Logf("Expected 0, got %v", value)
		t.Fail()
	}
}

func TestStack_Peek(t *testing.T) {
	stack := NewStack()
	stack.Push(0)
	if value, _ := stack.Peek(); value.(int) != 0 {
		t.Logf("Expected 0, got %v", value)
		t.Fail()
	}
}

func TestStack_MultiplePushesAndPops(t *testing.T) {
	stack := NewStack()

	count := 4

	for index := 0; index < count; index++ {
		stack.Push(index)
	}

	for index := count - 1; index >= 0; index-- {
		if value, _ := stack.Pop(); value.(int) != index {
			t.Logf("Expected %v got %v", index, value)
			t.Fail()
		}
	}
}

func TestStack_PopWhenEmpty(t *testing.T) {
	stack := NewStack()
	_, err := stack.Pop()
	if err == nil {
		t.Log("Expected error")
		t.Fail()
	}
}

func TestStack_PeekWhenEmpty(t *testing.T) {
	stack := NewStack()
	_, err := stack.Peek()
	if err == nil {
		t.Log("Expected error")
		t.Fail()
	}
}

func TestStack_PeekDoesNotChangeLength(t *testing.T) {
	stack := NewStack()
	stack.Push(0)
	_, _ = stack.Peek()
	if stack.IsEmpty() {
		t.Log("Stack should not be empty after Peek")
		t.Fail()
	}
}
