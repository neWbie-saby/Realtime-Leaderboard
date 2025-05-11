package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/neWbie-saby/leaderboard/internal/api"
	"github.com/neWbie-saby/leaderboard/internal/database"
	"github.com/neWbie-saby/leaderboard/internal/middlewares"
)

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("PORT not found in the environment")
	}

	dbString := os.Getenv("DB_URL")

	if dbString == "" {
		log.Fatal("DB_URL not found in the environment")
	}

	conn, err := sql.Open("postgres", dbString)

	if err != nil {
		log.Fatal("Can't connect to database")
	}

	db := database.New(conn)

	router := fiber.New(fiber.Config{
		ErrorHandler: middlewares.ErrorHandler,
	})

	router.Use(middlewares.CORS())
	// router.Use(middlewares.RateLimiter())
	router.Use(middlewares.Logger)

	apiCfg := api.ApiConfig{
		DB: db,
	}

	v1 := router.Group("/v1")

	v1.Get("/healthz", api.HandlerReadiness)
	v1.Get("/err", api.HandlerErr)

	v1.Post("/register", apiCfg.HandlerRegister)
	v1.Post("/login", apiCfg.HandlerLogin)

	v1.Get("/users/:username", middlewares.AuthenticateToken, apiCfg.HandlerGetUserByUsername)

	v1.Post("/matches", middlewares.AuthenticateToken, apiCfg.HandlerCreateMatch)
	v1.Post("/matches/:match_id/scores", middlewares.AuthenticateToken, apiCfg.HandlerPushMatchScores)

	log.Printf("Fiber Server starting on port %v", portString)
	if err := router.Listen(":" + portString); err != nil {
		log.Fatal(err)
	}
}
