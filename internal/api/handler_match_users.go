package api

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/neWbie-saby/leaderboard/internal/database"
	"github.com/neWbie-saby/leaderboard/internal/ds"
	"github.com/neWbie-saby/leaderboard/internal/models"
)

func (apiCfg *ApiConfig) HandlerPushMatchScores(f *fiber.Ctx) error {
	var input models.ScoresInput
	var processingErrors []string

	userId, ok := f.Locals("user_id").(int)
	if !ok {
		return RespondWithError(f, fiber.StatusInternalServerError, fmt.Sprint("User ID missing from context or invalid type"))
	}

	matchIDStr := f.Params("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 32)
	if err != nil {
		return RespondWithError(f, fiber.StatusBadRequest, fmt.Sprintf("Error converting string to int32: %v", err)) // code - 400
	}

	if err := f.BodyParser(&input); err != nil {
		return RespondWithError(f, fiber.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err)) // code - 400
	}

	if len(input.Scores) > 0 {

		filteredUserIDs := ds.NewHashMap[int, int]()

		for _, score := range input.Scores {
			val, ok := filteredUserIDs.Get(score.UserID)
			if !ok {
				val = 0
			}
			filteredUserIDs.Set(score.UserID, val+score.Score)
		}

		for userKey, scoreVal := range filteredUserIDs.Items() {
			present, err := apiCfg.DB.MatchUserExists(f.Context(), database.MatchUserExistsParams{
				UserID:  int32(userKey),
				MatchID: int32(matchID),
			})
			if err != nil {
				processingErrors = append(processingErrors, fmt.Sprintf("Failed to check user %v of match %v: %v", userKey, scoreVal, err))
				continue
			}

			if present {
				rowsAffected, err := apiCfg.DB.UpdateUserScoreOnMatch(f.Context(), database.UpdateUserScoreOnMatchParams{
					MatchID: int32(matchID),
					UserID:  int32(userKey),
					Score:   int32(scoreVal),
					ID:      int32(userId),
				})
				if err != nil {
					processingErrors = append(processingErrors, fmt.Sprintf("Match User updation failed of user_id %v: %v", userKey, err))
					continue
				}
				if rowsAffected == 0 {
					processingErrors = append(processingErrors, fmt.Sprint("Match User updator must be 'scorer' or 'official'"))
					continue
				}
			} else {
				rowsCreated, err := apiCfg.DB.AddUserToMatch(f.Context(), database.AddUserToMatchParams{
					MatchID: int32(matchID),
					UserID:  int32(userKey),
					Score:   int32(scoreVal),
					ID:      int32(userId),
				})
				if err != nil {
					processingErrors = append(processingErrors, fmt.Sprintf("Match User creation failed of user_id %v: %v", userKey, err))
					continue
				}
				if rowsCreated == 0 {
					processingErrors = append(processingErrors, fmt.Sprint("Match User creator must be 'scorer' or 'official'"))
				}
			}
		}

		if len(processingErrors) > 0 {
			if len(processingErrors) == filteredUserIDs.Len() {
				return RespondWithJSON(f, fiber.StatusInternalServerError, fiber.Map{
					"errors": processingErrors,
				})
			}
			return RespondWithJSON(f, fiber.StatusMultiStatus, fiber.Map{
				"errors": processingErrors,
			})
		}
	}

	ctx := context.WithValue(context.Background(), "user_id", userId)

	go func(ctx context.Context) {
		if err := apiCfg.CalculateAndUpdateWinner(int(matchID), ctx); err != nil {
			log.Printf("Error calculating and updating winner: %v", err)
		}
	}(ctx)

	return RespondWithJSON(f, fiber.StatusOK, fiber.Map{ // code - 200
		"success": true,
	})
}
