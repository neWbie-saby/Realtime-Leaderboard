package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"log"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/neWbie-saby/leaderboard/internal/api"
	"github.com/neWbie-saby/leaderboard/internal/database"
	"github.com/neWbie-saby/leaderboard/internal/httpserver"
	"github.com/neWbie-saby/leaderboard/internal/scheduler"
)

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("PORT not found in the environment")
	}

	grpcServerPortString := os.Getenv("GRPC_SERVER_PORT")

	if grpcServerPortString == "" {
		log.Fatal("GRPC_SERVER_PORT not found in the environment")
	}

	dbString := os.Getenv("DB_URL")

	if dbString == "" {
		log.Fatal("DB_URL not found in the environment")
	}

	conn, err := sql.Open("postgres", dbString)

	if err != nil {
		log.Fatal("Can't connect to database")
	}

	defer func() {
		log.Println("Closing database conection...")
		if err := conn.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed.")
		}
	}()

	db := database.New(conn)

	ctx, cancel := context.WithCancel(context.Background())
	var appWg sync.WaitGroup

	apiCfg := api.ApiConfig{
		DB: db,
	}

	// Channel to listen for OS signals (Ctrl+C, `kill` command)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// Alternatively, signal.NotifyContext may also be used

	httpServer := httpserver.New(apiCfg, portString)

	appWg.Add(1)
	go httpServer.Run(ctx, &appWg)

	// ==== Next Lone Go-routines to be added here ====
	go func() {
		matchAnalyzer, err := scheduler.NewMatchAnalysisProcessor(ctx, db, 10, time.Minute, grpcServerPortString, uuid.New().String())
		if err != nil {
			log.Printf("Failed to create Match Analyzer: %v", err)
		} else {
			appWg.Add(1)
			go matchAnalyzer.Start(&appWg)
		}
	}()

	// Go-routine to clear un-released and expired processing_locks
	appWg.Add(1)
	go scheduler.RunExpiryProcessor(ctx, &appWg, db, 2*time.Minute)

	// === Blocking main GoRoutine Channel ===
	// Block main GoRoutine until a signal is received
	sig := <-sigChan
	log.Printf("Received OS signal: %s. Initiating graceful shutdown...", sig)

	cancel()

	log.Println("Waiting for all background services and HTTP server to shut down...")
	appWg.Wait()

	log.Println("Application shutdown complete.")
}
