package api

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/neWbie-saby/leaderboard/internal/models"
)

func (apiCfg *ApiConfig) HandlerGetUserByUsername(f *fiber.Ctx) error {
	username := f.Params("username")

	user, err := apiCfg.DB.GetUserByUserName(f.Context(), username)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result") {
			return RespondWithError(f, fiber.StatusNotFound, fmt.Sprintf("User not found with username: %v", username)) // code - 404
		}
		return RespondWithError(f, fiber.StatusInternalServerError, fmt.Sprintf("Failed to get user of username %v: %v", username, err)) // code - 500
	}

	return RespondWithJSON(f, fiber.StatusOK, fiber.Map{ // code - 200
		"user_details": models.DatabaseUserToUserFetchedByUserName(user),
	})
}
