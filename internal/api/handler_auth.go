package api

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/neWbie-saby/leaderboard/internal/database"
	"github.com/neWbie-saby/leaderboard/internal/models"
	"github.com/neWbie-saby/leaderboard/internal/utils"
)

func (apiCfg *ApiConfig) HandlerRegister(f *fiber.Ctx) error {
	var input models.RegisterInput

	if err := f.BodyParser(&input); err != nil {
		return RespondWithError(f, fiber.StatusBadRequest, fmt.Sprintf("Error parsing register JSON: %v", err)) // code - 400
	}

	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		return RespondWithError(f, fiber.StatusInternalServerError, fmt.Sprintf("Password hashing failed: %v", err)) // code - 500
	}

	user, err := apiCfg.DB.CreateUser(f.Context(), database.CreateUserParams{
		Username:     input.UserName,
		Email:        input.Email,
		PasswordHash: hashed,
	})
	if err != nil {
		return RespondWithError(f, fiber.StatusConflict, fmt.Sprintf("User creation failed: %v", err)) // code - 409
	}

	token, err := utils.GenerateJWT(int(user.ID))
	if err != nil {
		return RespondWithError(f, fiber.StatusInternalServerError, fmt.Sprintf("Token generation failed: %v", err)) // code - 500
	}

	return RespondWithJSON(f, fiber.StatusOK, fiber.Map{ // code - 200
		"user_details": models.DatabaseUserToUser(user),
		"token":        token,
	})
}

func (apiCfg *ApiConfig) HandlerLogin(f *fiber.Ctx) error {
	var input models.LoginInput

	if err := f.BodyParser(&input); err != nil {
		return RespondWithError(f, fiber.StatusBadRequest, fmt.Sprintf("Error parsing register JSON: %v", err)) // code - 400
	}

	var user database.User
	var err error

	user, err = apiCfg.DB.GetUserByUserName(f.Context(), input.Identifier)
	if err != nil {
		user, err = apiCfg.DB.GetUserByEmail(f.Context(), input.Identifier)
		if err != nil {
			return RespondWithError(f, fiber.StatusUnauthorized, fmt.Sprintf("Invalid credentials: %v", err)) // code - 401
		}
	}

	if err := utils.CheckPasswordHash(input.Password, user.PasswordHash); err != nil {
		return RespondWithError(f, fiber.StatusUnauthorized, fmt.Sprintf("Invalid credentials: %v", err)) // code - 401
	}

	token, err := utils.GenerateJWT(int(user.ID))
	if err != nil {
		return RespondWithError(f, fiber.StatusInternalServerError, fmt.Sprintf("Token generation failed: %v", err)) // code - 500
	}

	return RespondWithJSON(f, fiber.StatusOK, fiber.Map{ // code - 200
		"user_details": models.DatabaseUserToUser(user),
		"token":        token,
	})
}

func (apiCfg *ApiConfig) AuthenticateToken(f *fiber.Ctx) error {
	authHeader := f.Get("Authorization")
	if authHeader == "" {
		return RespondWithError(f, fiber.StatusUnauthorized, "No authentication info found")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return RespondWithError(f, fiber.StatusUnauthorized, "Invalid Authorization format. Expected 'Bearer <token>'")
	}

	tokenStr := tokenParts[1]
	userId, err := utils.ParseJWT(tokenStr)
	if err != nil {
		return RespondWithError(f, fiber.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
	}

	f.Locals("user_id", userId)

	return f.Next()
}
