package scheduler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/neWbie-saby/leaderboard/internal/database"
	"github.com/neWbie-saby/leaderboard/internal/grpcclient/analyticsclient"
)

type MatchAnalysisProcessor struct {
	ctx            context.Context
	db             *database.Queries
	analysisClient *analyticsclient.AnalyticsClient
	concurrency    int32
	interval       time.Duration
	lockedBy       string
}

func NewMatchAnalysisProcessor(
	parentCtx context.Context,
	db *database.Queries,
	concurrency int,
	interval time.Duration,
	grpcPort string,
	appId string,
) (*MatchAnalysisProcessor, error) {

	grpcServerAddr := fmt.Sprintf("127.0.0.1:%s", grpcPort)
	lockedBy := fmt.Sprintf("%s_MatchAnalysisProcessor_%s", appId, uuid.New().String())

	maxRetries := 5
	initialRetryGap := 2 * time.Second
	maxRetryGap := 30 * time.Second
	connectTimeOut := 15 * time.Second

	var analyticsClient *analyticsclient.AnalyticsClient
	var connectErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("MatchAnalysisProcessor: Attempt %d of %d to establish gRPC connection to %s (timeout: %v)",
			attempt,
			maxRetries,
			grpcServerAddr,
			connectTimeOut,
		)

		grpcCtx, grpcCtxCancel := context.WithTimeout(parentCtx, 15*time.Second)

		analyticsClient, connectErr = analyticsclient.NewAnalyticsClient(grpcCtx, grpcServerAddr)
		grpcCtxCancel()

		if connectErr == nil {
			log.Println("MatchAnalysisProcessor: Successfully established gRPC connection.")
			break
		}

		log.Printf("Failed to establish initial gRPC connection to analytics service at %s: %v", grpcServerAddr, connectErr)

		if attempt == maxRetries {
			return nil, fmt.Errorf("MatchAnalysisProcessor: Exceeded maximum retries for gRPC connection attempts (%d) to %s: %w", maxRetries, grpcServerAddr, connectErr)
		}

		retryGapDuration := initialRetryGap * time.Duration(1<<(attempt-1)) // 2 to the power (attempt - 1)
		if retryGapDuration > maxRetryGap {
			retryGapDuration = maxRetryGap
		}
		log.Printf("MatchAnalysisProcessor: Retrying gRPC connection in %v...", retryGapDuration)

		select {
		case <-time.After(retryGapDuration):
			// continue
		case <-parentCtx.Done():
			return nil, fmt.Errorf("MatchAnalysisProcessor: Aborting gRPC connection retries during backoff due to context cancellation (%v): %w", parentCtx.Err(), connectErr)
		}
	}

	return &MatchAnalysisProcessor{
		ctx:            parentCtx,
		db:             db,
		analysisClient: analyticsClient,
		concurrency:    int32(concurrency),
		interval:       interval,
		lockedBy:       lockedBy,
	}, nil
}

func (p *MatchAnalysisProcessor) Close() error {
	if p.analysisClient != nil {
		log.Println("Closing Match Analysis Processor's gRPC connection")
		return p.analysisClient.Close()
	}
	return nil
}

func (p *MatchAnalysisProcessor) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		if err := p.Close(); err != nil {
			log.Printf("Error closing MatchAnalysisProcessor: %v", err)
		}
	}()

	log.Printf("Starting match analysis processor worker: %s, interval %s", p.lockedBy, p.interval)

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			log.Println("MatchAnalysisProcessor: Received shutdown signal. Exiting.")
			return

		case <-ticker.C:
			p.processBatch()
		}
	}
}

func (p *MatchAnalysisProcessor) processBatch() {
	matchIds, err := p.db.GetMatchIdsForAnalysis(p.ctx, p.concurrency)
	if err != nil {
		log.Println("Error fetching match IDs:", err)
		return
	}

	if len(matchIds) == 0 {
		log.Println("No new match IDs to analyze. Skipping batch.")
		return
	}

	requestWg := &sync.WaitGroup{}

	for _, id := range matchIds {
		requestWg.Add(1)

		go p.handleMatchAnalysis(requestWg, id)
	}
	requestWg.Wait()
	log.Println("Match Analysis Processor: Batch of gRPC requests (and lock-unlock attempts) complete")
}

func (p *MatchAnalysisProcessor) handleMatchAnalysis(wg *sync.WaitGroup, matchId int32) {
	defer wg.Done()

	rowsAffected, err := p.db.AcquireLock(p.ctx, database.AcquireLockParams{
		ResourceType: "match_analysis",
		ResourceID:   strconv.Itoa(int(matchId)),
		LockedBy:     p.lockedBy,
	})
	if err != nil {
		log.Printf("Match Analysis failed for match ID %d: %v", matchId, err)
		return
	}

	if rowsAffected > 0 {
		log.Printf("Successfully acquired lock for the match ID: %v", matchId)

		err = p.analysisClient.TriggerAnalysis(matchId)
		if err != nil {
			log.Printf("Failed to trigger match analysis for match ID %d: %v", matchId, err)
		}

		err = p.db.ReleaseLock(p.ctx, database.ReleaseLockParams{
			ResourceType: "match_analysis",
			ResourceID:   strconv.Itoa(int(matchId)),
			LockedBy:     p.lockedBy,
		})
		if err != nil {
			log.Printf("Failed to release lock for the match ID %d: %v", matchId, err)
		}
		log.Printf("Successfully released lock for the match ID: %v", matchId)

	} else {
		log.Printf("Failed to acquire lock for the match ID %d: %v", matchId, err)
	}
}

func StartMatchAnalysis(ctx context.Context, wg *sync.WaitGroup, db *database.Queries, concurrency int, interval time.Duration, gRPCPortStr string, appID string) {
	defer wg.Done()

	lockedBy := fmt.Sprintf("%s_MatchAnalysisProcessor_%s", appID, uuid.New().String())

	log.Printf("Starting match analysis processor worker: %s, interval %s", lockedBy, interval)

	grpcAddr := fmt.Sprintf("127.0.0.1:%s", gRPCPortStr)

	grpcConnCtx, grpcConncancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer grpcConncancel()

	analysisClient, err := analyticsclient.NewAnalyticsClient(grpcConnCtx, grpcAddr)
	if err != nil {
		log.Printf("Failed to establish initial gRPC connection to analytics service at %s: %v", grpcAddr, err)
		return
	}
	defer analysisClient.Close()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Match Analysis processor received shutdown signal. Exiting.")
			return

		case <-ticker.C:
			log.Println("Match Analysis processor triggered. Initiating batch of gRPC requests...")

			matchIds, err := db.GetMatchIdsForAnalysis(ctx, int32(concurrency))
			if err != nil {
				log.Println("Error fetching match IDs:", err)
				continue
			}

			requestWg := &sync.WaitGroup{}

			for _, id := range matchIds {
				requestWg.Add(1)

				go func(matchId int32) {
					defer requestWg.Done()

					rowsAffected, err := db.AcquireLock(ctx, database.AcquireLockParams{
						ResourceType: "match_analysis",
						ResourceID:   strconv.Itoa(int(matchId)),
						LockedBy:     lockedBy,
					})
					if err != nil {
						log.Printf("Match Analysis failed for match ID %d: %v", matchId, err)
						return
					}
					if rowsAffected > 0 {
						log.Printf("Successfully acquired lock for the match ID: %v", matchId)

						err = analysisClient.TriggerAnalysis(matchId)
						if err != nil {
							log.Printf("Failed to trigger match analysis for match ID %d: %v", matchId, err)
						}

						err = db.ReleaseLock(ctx, database.ReleaseLockParams{
							ResourceType: "match_analysis",
							ResourceID:   strconv.Itoa(int(matchId)),
							LockedBy:     lockedBy,
						})
						if err != nil {
							log.Printf("Failed to release lock for the match ID %d: %v", matchId, err)
						}
						log.Printf("Successfully released lock for the match ID: %v", matchId)

					} else {
						log.Printf("Failed to acquire lock for the match ID %d: %v", matchId, err)
					}
				}(id)
			}
			requestWg.Wait()
			log.Println("Match Analysis Processor: Batch of gRPC requests (and lock-unlock attempts) complete")
		}
	}
}
