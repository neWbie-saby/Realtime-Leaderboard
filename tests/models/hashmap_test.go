package models_test

import (
	"testing"

	"github.com/neWbie-saby/leaderboard/internal/ds"
)

func TestHashMap(t *testing.T) {
	hm := ds.NewHashMap[string, int]()

	hm.Set("player1", 100)
	hm.Set("player2", 200)

	val, ok := hm.Get("player1")
	if !ok || val != 100 {
		t.Errorf("Expected player1 = 100, got %v", val)
	}

	hm.Delete("player1")
	_, ok = hm.Get("player1")
	if ok {
		t.Error("Expected player1 to be deleted")
	}
}
