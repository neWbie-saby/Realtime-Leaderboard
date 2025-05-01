package models_test

import (
	"testing"

	"github.com/neWbie-saby/leaderboard/internal/ds"
)

func TestGraph(t *testing.T) {
	g := ds.NewGraph()

	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)

	edges1 := g.GetEdges(1)
	if len(edges1) != 2 || edges1[0] != 2 || edges1[1] != 3 {
		t.Errorf("Expected [2 3], got %v", edges1)
	}

	g.RemoveEdge(1, 2)
	edges1 = g.GetEdges(1)
	if len(edges1) != 1 || edges1[0] != 3 {
		t.Errorf("Expected [3], got %v", edges1)
	}
}
