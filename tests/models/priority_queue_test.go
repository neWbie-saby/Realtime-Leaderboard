package models_test

import (
	"container/heap"
	"testing"

	"github.com/neWbie-saby/leaderboard/internal/ds"
	"github.com/neWbie-saby/leaderboard/internal/models"
)

func TestPriorityQueue(t *testing.T) {
	pq := ds.NewPriorityQueue()

	heap.Push(pq, &models.UserScore{UserID: 1, Score: 50})
	heap.Push(pq, &models.UserScore{UserID: 2, Score: 90})
	heap.Push(pq, &models.UserScore{UserID: 3, Score: 75})

	item := heap.Pop(pq).(*models.UserScore)
	if item.UserID != 2 {
		t.Errorf("Expected UserID 2 (Score 90), got %v", item.UserID)
	}

	item = heap.Pop(pq).(*models.UserScore)
	if item.UserID != 3 {
		t.Errorf("Expected UserID 3 (Score 75), got %v", item.UserID)
	}

	item = heap.Pop(pq).(*models.UserScore)
	if item.UserID != 1 {
		t.Errorf("Expected UserID 1 (Score 50), got %v", item.UserID)
	}
}
