package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/neWbie-saby/leaderboard/internal/database"
)

func RunExpiryProcessor(
	parentCtx context.Context,
	wg *sync.WaitGroup,
	db *database.Queries,
	interval time.Duration) {

	defer wg.Done()

	log.Printf("Starting ExpiredLocksTTLProcessor: interval %s", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-parentCtx.Done():
			log.Println("ExpiredLocksTTLProcessor: Received shutdown signal. Exiting.")
			return

		case <-ticker.C:
			rowsAffected, err := db.DeleteExpiredLocks(parentCtx)
			if err != nil {
				log.Printf("Failed to delete expired processing_locks rows: %v", err)
			} else {
				log.Printf("Successfully deleted %d expired processing_locks rows", rowsAffected)
			}
		}
	}
}
