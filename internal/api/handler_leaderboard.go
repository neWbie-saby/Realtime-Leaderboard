package api

import (
	"container/heap"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/neWbie-saby/leaderboard/internal/database"
	"github.com/neWbie-saby/leaderboard/internal/ds"
	"github.com/neWbie-saby/leaderboard/internal/models"
)

func (apiCfg *ApiConfig) HandlerGetMatchLeaderboard(f *fiber.Ctx) error {
	matchID, err := strconv.ParseInt(f.Params("match_id"), 10, 32)
	if err != nil {
		return RespondWithError(f, fiber.StatusBadRequest, fmt.Sprintf("Error converting string to int32: %v", err)) // code - 400
	}

	userId, ok := f.Locals("user_id").(int)
	if !ok {
		return RespondWithError(f, fiber.StatusInternalServerError, fmt.Sprint("User ID missing from context or invalid type"))
	}

	limit, err := strconv.Atoi(f.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	scores, err := apiCfg.DB.GetMatchUserScoresAndUserNames(f.Context(), database.GetMatchUserScoresAndUserNamesParams{
		MatchID: int32(matchID),
		ID:      int32(userId),
	})
	if err != nil {
		return RespondWithError(f, fiber.StatusInternalServerError, fmt.Sprintf("Failed to get scores of match %v: %v", matchID, err)) // code - 500
	}
	if len(scores) == 0 {
		return RespondWithError(f, fiber.StatusNotFound, fmt.Sprintf("Failed to get scores of match %v since current player user does not belong to this match, or user is not 'scorer' or 'player'", matchID)) // code - 404
	}

	pq := ds.NewPriorityQueue()

	for _, score := range scores {
		heap.Push(pq, &models.UserScoreAndUserNameOutput{
			UserID:   score.ID,
			Score:    score.Score,
			Username: score.Username,
		})
	}

	var top []models.UserScoreAndUserNameOutput
	for i := 0; i < limit && pq.Len() > 0; i++ {
		item := heap.Pop(pq).(*models.UserScoreAndUserNameOutput)
		top = append(top, *item)
	}

	return RespondWithJSON(f, fiber.StatusOK, fiber.Map{ // code - 200
		"match_details": fiber.Map{
			"match_id": matchID,
			"top":      top,
		},
	})
}
