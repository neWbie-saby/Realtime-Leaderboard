package models_test

import (
	"testing"

	"github.com/neWbie-saby/leaderboard/internal/ds"
)

func TestStack(t *testing.T) {
	stack := ds.NewStack[int]()

	stack.Push(10)
	stack.Push(20)

	val, ok := stack.Peek()
	if !ok || val != 20 {
		t.Errorf("Expected Peek = 20, got %v", val)
	}

	val, ok = stack.Pop()
	if !ok || val != 20 {
		t.Errorf("Expected Pop = 20, got %v", val)
	}

	val, ok = stack.Pop()
	if !ok || val != 10 {
		t.Errorf("Expected Pop = 10, got %v", val)
	}

	_, ok = stack.Pop()
	if ok {
		t.Error("Expected Pop on empty stack to fail")
	}
}
