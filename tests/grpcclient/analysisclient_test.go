package grpcclient_test

import (
	"context"
	"testing"
	"time"

	client "github.com/neWbie-saby/leaderboard/internal/grpcclient/analyticsclient"
)

func TestTriggerAnalysis(t *testing.T) {
	addr := "127.0.0.1:9090"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := client.NewAnalyticsClient(ctx, addr)
	if err != nil {
		t.Fatalf("failed to create analystics client %v", err)
	}
	defer c.Close()

	err = c.TriggerAnalysis(699999608)
	if err != nil {
		t.Errorf("TriggerAnalysis should not return error - %v", err)
	}
}
