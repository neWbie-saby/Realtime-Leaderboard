package httpserver

import (
	"context"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/neWbie-saby/leaderboard/internal/api"
	"github.com/neWbie-saby/leaderboard/internal/middlewares"
)

type HTTPServer struct {
	App        *fiber.App
	ListenAddr string
}

func New(apiCfg api.ApiConfig, port string) *HTTPServer {
	router := fiber.New(fiber.Config{
		ErrorHandler: middlewares.ErrorHandler,
	})

	router.Use(middlewares.CORS())
	// router.Use(middlewares.RateLimiter())
	router.Use(middlewares.Logger)

	v1 := router.Group("/v1")

	api.RegisterRoutes(v1, apiCfg)

	return &HTTPServer{
		App:        router,
		ListenAddr: port,
	}
}

func (s *HTTPServer) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		<-ctx.Done()
		log.Println("HTTP Server: Initiating graceful shutdown....")

		if err := s.App.Shutdown(); err != nil {
			log.Printf("HTTP Server: Shutdown error: %v", err)
		} else {
			log.Println("HTTP Server: Shutdown complete.")
		}
	}()

	log.Printf("Fiber HTTP Server starting on port %v", s.ListenAddr)
	if err := s.App.Listen(":" + s.ListenAddr); err != nil {
		log.Fatalf("HTTP Server: Listen error: %v", err)
	}
	log.Println("HTTP server: Goroutine exited.")
}
