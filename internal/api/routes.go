package api

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(v1 fiber.Router, apiCfg ApiConfig) {

	v1.Get("/healthz", HandlerReadiness)
	v1.Get("/err", HandlerErr)

	v1.Post("/register", apiCfg.HandlerRegister)
	v1.Post("/login", apiCfg.HandlerLogin)

	v1.Get("/users/:username", apiCfg.AuthenticateToken, apiCfg.HandlerGetUserByUsername)

	v1.Post("/matches", apiCfg.AuthenticateToken, apiCfg.HandlerCreateMatch)
	v1.Post("/matches/:match_id/scores", apiCfg.AuthenticateToken, apiCfg.HandlerPushMatchScores)

	v1.Post("/matches/:match_id/compute-winners", apiCfg.AuthenticateToken, apiCfg.HandlerCalculateWinner)

	v1.Get("/matches/:match_id/leaderboard/scores", apiCfg.AuthenticateToken, apiCfg.HandlerGetMatchLeaderboard)
}
