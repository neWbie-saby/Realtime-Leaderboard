package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	router := fiber.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:  "*",
		ExposeHeaders: "Link",
		MaxAge:        int(5 * time.Minute / time.Second),
	}))

	apiCfg := api.ApiConfig{
		DB: db,
	}

	v1 := router.Group("/v1")

	v1.Get("/healthz", api.HandlerReadiness)
	v1.Get("/err", api.HandlerErr)

	v1.Post("/register", apiCfg.HandlerRegister)
	v1.Post("/login", apiCfg.HandlerLogin)

	v1.Get("/users/by-username", middlewares.AuthenticateToken, apiCfg.HandlerGetUserByUsername)

	log.Printf("Fiber Server starting on port %v", portString)
	if err := router.Listen(":" + portString); err != nil {
		log.Fatal(err)
	}
}
