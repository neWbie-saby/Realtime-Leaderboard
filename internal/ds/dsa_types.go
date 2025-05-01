package ds

import "github.com/neWbie-saby/leaderboard/internal/models"

type HashMap[K comparable, V any] struct {
	data map[K]V
}

type Stack[T any] struct {
	data []T
}

type PriorityQueue []*models.UserScore

type Graph struct {
	adj map[int][]int
}
