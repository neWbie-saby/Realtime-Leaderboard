package api

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/neWbie-saby/leaderboard/internal/database"
	"github.com/neWbie-saby/leaderboard/internal/models"
)

func (apiCfg *ApiConfig) CalculateAndUpdateWinner(matchID int, ctx context.Context) error {
	log.Printf("Calculating winner for match %v", matchID)

	var processingErrors []string

	scores, err := apiCfg.DB.GetMatchUserScores(ctx, int32(matchID))
	if err != nil {
		return fmt.Errorf("error fetching scores for match %v: %w", matchID, err)
	}

	winnerScores := determineWinner(models.DatabaseUserScoresToUserScores(scores))

	if len(winnerScores) > 0 {

		err := apiCfg.DB.DeleteMatchWinners(ctx, int32(matchID))
		if err != nil {
			return fmt.Errorf("error deleting scores for match %v: %w", matchID, err)
		}

		for _, winnerScore := range winnerScores {
			err := apiCfg.DB.AddWinnerOfMatch(ctx, database.AddWinnerOfMatchParams{
				MatchID: int32(matchID),
				UserID:  int32(winnerScore.UserID),
			})
			if err != nil {
				processingErrors = append(processingErrors, fmt.Sprintf("Match Winner creation failed of user_id %v: %v", winnerScore.UserID, err))
				continue
			}
		}

		if len(processingErrors) > 0 {
			return fmt.Errorf("errors encountered while adding match winners:\n%s", strings.Join(processingErrors, "\n"))
		}
	}

	log.Printf("Successfully calculated winner for match %v", matchID)
	return nil
}

func determineWinner(scores []models.UserScore) []models.UserScore {
	if len(scores) == 0 {
		return nil
	}

	maxScore := math.MinInt

	for _, score := range scores {
		if score.Score > maxScore {
			maxScore = score.Score
		}
	}

	var maxUserScores []models.UserScore

	for _, score := range scores {
		if maxScore == score.Score {
			maxUserScores = append(maxUserScores, models.UserScore{
				Score:  score.Score,
				UserID: score.UserID,
			})
		}
	}

	return maxUserScores
}

func (apiCfg *ApiConfig) HandlerCalculateWinner(f *fiber.Ctx) error {
	matchIDStr := f.Params("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 32)
	if err != nil {
		return RespondWithError(f, fiber.StatusBadRequest, fmt.Sprintf("Error converting string to int32: %v", err)) // code - 400
	}

	err = apiCfg.CalculateAndUpdateWinner(int(matchID), f.Context())
	if err != nil {
		return RespondWithError(f, fiber.StatusInternalServerError, fmt.Sprintf("Failed to calculate winner: %v", err)) // code - 500
	}

	return RespondWithJSON(f, fiber.StatusOK, fiber.Map{ // code - 200
		"message": "Winner calculation initiated",
	})
}
