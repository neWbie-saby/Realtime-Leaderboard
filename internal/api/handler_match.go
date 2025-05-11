package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/neWbie-saby/leaderboard/internal/models"
)

func (apiCfg *ApiConfig) HandlerCreateMatch(f *fiber.Ctx) error {
	var input models.CreateMatchInput

	if err := f.BodyParser(&input); err != nil {
		return RespondWithError(f, fiber.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err)) // code - 400
	}

	match, err := apiCfg.DB.CreateMatch(f.Context(), input.MatchType)
	if err != nil {
		return RespondWithError(f, fiber.StatusConflict, fmt.Sprintf("Match creation failed: %v", err)) // code - 409
	}

	return RespondWithJSON(f, fiber.StatusOK, fiber.Map{ // code - 200
		"match_details": models.DatabaseCreatedMatchToMatch(match),
	})
}
