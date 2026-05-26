package api

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/neWbie-saby/leaderboard/internal/database"
	"github.com/neWbie-saby/leaderboard/internal/models"
)

func (apiCfg *ApiConfig) HandlerCreateMatch(f *fiber.Ctx) error {
	userId, ok := f.Locals("user_id").(int)
	if !ok {
		return RespondWithError(f, fiber.StatusInternalServerError, fmt.Sprint("User ID missing from context or invalid type"))
	}

	var input models.CreateMatchInput
	if err := f.BodyParser(&input); err != nil {
		return RespondWithError(f, fiber.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err)) // code - 400
	}

	match, err := apiCfg.DB.CreateMatch(f.Context(), database.CreateMatchParams{
		ID:        int32(userId),
		MatchType: input.MatchType,
	})
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result") {
			return RespondWithError(f, fiber.StatusNotFound, fmt.Sprint("Match creator must be an 'official'")) // code - 404
		}
		return RespondWithError(f, fiber.StatusConflict, fmt.Sprintf("Match creation failed: %v", err)) // code - 409
	}

	return RespondWithJSON(f, fiber.StatusOK, fiber.Map{ // code - 200
		"match_details": models.DatabaseCreatedMatchToMatch(match),
	})
}

func (apiCfg *ApiConfig) HandlerGetPlayerMatches(f *fiber.Ctx) error {
	userId, ok := f.Locals("user_id").(int)
	if !ok {
		return RespondWithError(f, fiber.StatusInternalServerError, fmt.Sprint("User ID missing from context or invalid type"))
	}

	matches, err := apiCfg.DB.GetUserPlayerMatches(f.Context(), int32(userId))
	if err != nil {
		return RespondWithError(f, fiber.StatusInternalServerError, fmt.Sprintf("Failed to get matches for user_id %d: %v", userId, err))
	}

	cleanMatches := models.DatabasePlayerMatchesToMatches(matches)

	return RespondWithJSON(f, fiber.StatusOK, fiber.Map{
		"matches": cleanMatches,
	})
}

func (apiCfg *ApiConfig) HandlerGetPlayerMatchScore(f *fiber.Ctx) error {
	userId, ok := f.Locals("user_id").(int)
	if !ok {
		return RespondWithError(f, fiber.StatusInternalServerError, fmt.Sprint("User ID missing from context or invalid type"))
	}

	matchID, err := strconv.ParseInt(f.Params("match_id"), 10, 32)
	if err != nil {
		return RespondWithError(f, fiber.StatusBadRequest, fmt.Sprintf("Error converting string to int32: %v", err)) // code - 400
	}

	playerMatchDetails, err := apiCfg.DB.GetUserPlayerMatchScore(f.Context(), database.GetUserPlayerMatchScoreParams{
		ID:   int32(userId),
		ID_2: int32(matchID),
	})
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result") {
			return RespondWithError(f, fiber.StatusNotFound, fmt.Sprint("No score found for the selected match")) // code - 404
		}
		return RespondWithError(f, fiber.StatusInternalServerError, fmt.Sprintf("Failed to get score for the selected match: %v", err)) // code - 500
	}

	return RespondWithJSON(f, fiber.StatusOK, fiber.Map{ // code - 200
		"match_details": models.DatabasePlayerMatchScoreToScore(playerMatchDetails),
	})
}
