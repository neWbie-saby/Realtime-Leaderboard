package grpcclient_test

import (
	"testing"

	client "github.com/neWbie-saby/leaderboard/internal/grpcclient/analyticsclient"
)

func TestTriggerAnalysis(t *testing.T) {
	addr := "127.0.0.1:9090"

	c, err := client.NewAnalyticsClient(addr)
	if err != nil {
		t.Fatalf("failed to create analystics client %v", err)
	}
	defer c.Close()

	err = c.TriggerAnalysis(699999608)
	if err != nil {
		t.Errorf("TriggerAnalysis should not return error - %v", err)
	}
}
