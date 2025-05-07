package api

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/neWbie-saby/leaderboard/internal/models"
)

func (apiCfg *ApiConfig) HandlerGetUserByUsername(f *fiber.Ctx) error {
	var input models.UserNameInput

	if err := f.BodyParser(&input); err != nil {
		return RespondWithError(f, fiber.StatusBadRequest, fmt.Sprintf("Error parsing register JSON: %v", err)) // code - 400
	}

	user, err := apiCfg.DB.GetUserByUserName(f.Context(), input.UserName)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result") {
			return RespondWithError(f, fiber.StatusNotFound, fmt.Sprintf("User not found with username: %v", input.UserName)) // code - 404
		}
		return RespondWithError(f, fiber.StatusInternalServerError, fmt.Sprintf("Failed to get user of username %v: %v", input.UserName, err)) // code - 500
	}

	return RespondWithJSON(f, fiber.StatusOK, fiber.Map{ // code - 200
		"user_details": models.DatabaseUserToUserFetchedByUserName(user),
	})
}
