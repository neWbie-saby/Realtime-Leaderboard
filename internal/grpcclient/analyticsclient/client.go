package analyticsclient

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/neWbie-saby/leaderboard/proto/analytics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type AnalyticsClient struct {
	conn   *grpc.ClientConn
	client pb.AnalyticsServiceClient
}

func NewAnalyticsClient(ctx context.Context, addr string) (*AnalyticsClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	conn.Connect()

	for {
		currentState := conn.GetState()
		if currentState == connectivity.Ready {
			log.Printf("gRPC client connected to %s", addr)
			break
		}
		if !conn.WaitForStateChange(ctx, currentState) {
			return nil, fmt.Errorf("gRPC client connection to %s failed to become ready within timeout: %w", addr, ctx.Err())
		}
	}

	client := pb.NewAnalyticsServiceClient(conn)

	return &AnalyticsClient{
		conn:   conn,
		client: client,
	}, nil
}

func (a *AnalyticsClient) TriggerAnalysis(matchID int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := a.client.TriggerMatchAnalysis(ctx, &pb.TriggerRequest{MatchId: matchID})
	if err != nil {
		log.Println("Failed to trigger analysis")
		return err
	}

	log.Println("Triggered match analysis for match_id -", matchID)
	return nil
}

func (a *AnalyticsClient) Close() error {
	return a.conn.Close()
}
