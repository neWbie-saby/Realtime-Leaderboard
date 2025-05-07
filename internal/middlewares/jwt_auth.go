package middlewares

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/neWbie-saby/leaderboard/internal/api"
	"github.com/neWbie-saby/leaderboard/internal/utils"
)

func AuthenticateToken(f *fiber.Ctx) error {
	authHeader := f.Get("Authorization")
	if authHeader == "" {
		return api.RespondWithError(f, fiber.StatusUnauthorized, "No authentication info found")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return api.RespondWithError(f, fiber.StatusUnauthorized, "Invalid Authorization format. Expected 'Bearer <token>'")
	}

	tokenStr := tokenParts[1]
	userId, err := utils.ParseJWT(tokenStr)
	if err != nil {
		return api.RespondWithError(f, fiber.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
	}

	f.Locals("user_id", userId)

	return f.Next()
}
