package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"sync"
	"syscall"

	///Required for stub gRPC server

	// "context"
	"log"
	// "net"

	// pb "github.com/neWbie-saby/leaderboard/proto/analytics"
	// "google.golang.org/grpc"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/neWbie-saby/leaderboard/internal/api"
	"github.com/neWbie-saby/leaderboard/internal/database"
	"github.com/neWbie-saby/leaderboard/internal/httpserver"
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

	httpServer := httpserver.New(apiCfg, portString)

	appWg.Add(1)
	go httpServer.Run(ctx, &appWg)

	// ==== Next Lone Go-routines to be added here ====

	// === Main GoRoutine Channel ===
	// Channel to listen for OS signals (e.g., Ctrl+C, `kill` command)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block main GoRoutine until a signal is received
	sig := <-sigChan
	log.Printf("Received OS signal: %s. Initiating graceful shutdown...", sig)

	cancel()

	log.Println("Waiting for all background services and HTTP server to shut down...")
	appWg.Wait()

	log.Println("Application shutdown complete.")
}

// type stubAnalyticsServer struct {
// 	pb.UnimplementedAnalyticsServiceServer
// }

// func (s *stubAnalyticsServer) TriggerMatchAnalysis(ctx context.Context, req *pb.TriggerRequest) (*pb.TriggerResponse, error) {
// 	log.Printf("Stub received for TriggerMatchAnalysis for match_id %d", req.MatchId)
// 	return &pb.TriggerResponse{}, nil
// }

// func main() {
// 	lis, err := net.Listen("tcp", ":50051")
// 	if err != nil {
// 		log.Fatalf("Failed to listen: %v", err)
// 	}

// 	s := grpc.NewServer()
// 	pb.RegisterAnalyticsServiceServer(s, &stubAnalyticsServer{})

// 	log.Println("Stub gRPC listening on :50051")
// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("Failed to serve: %v", err)
// 	}
// }
