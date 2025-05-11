package api

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/neWbie-saby/leaderboard/internal/database"
	"github.com/neWbie-saby/leaderboard/internal/ds"
	"github.com/neWbie-saby/leaderboard/internal/models"
)

func (apiCfg *ApiConfig) HandlerPushMatchScores(f *fiber.Ctx) error {
	matchIDStr := f.Params("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 32)
	if err != nil {
		log.Fatalf("Error converting string to int32: %v", err)
	}

	var input models.ScoresInput

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
				return RespondWithError(f, fiber.StatusInternalServerError, fmt.Sprintf("Failed to check user %v of match %v: %v", userKey, scoreVal, err)) // code - 500
			}

			if present {
				err := apiCfg.DB.UpdateUserScoreOnMatch(f.Context(), database.UpdateUserScoreOnMatchParams{
					MatchID: int32(matchID),
					UserID:  int32(userKey),
					Score:   int32(scoreVal),
				})
				if err != nil {
					return RespondWithError(f, fiber.StatusConflict, fmt.Sprintf("Match User updation failed: %v", err)) // code - 409
				}
			} else {
				err := apiCfg.DB.AddUserToMatch(f.Context(), database.AddUserToMatchParams{
					MatchID: int32(matchID),
					UserID:  int32(userKey),
					Score:   int32(scoreVal),
				})
				if err != nil {
					return RespondWithError(f, fiber.StatusConflict, fmt.Sprintf("Match User creation failed: %v", err)) // code - 409
				}
			}
		}
	}

	return RespondWithJSON(f, fiber.StatusOK, fiber.Map{ // code - 200
		"success": true,
	})
}
