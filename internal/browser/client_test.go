package browser

import (
	"context"
	"testing"

	"github.com/akkul143/ai-analysis/internal/config"
)

func TestDisabledWorkerDoesNotLaunchProcess(t *testing.T) {
	client := NewClient(config.BrowserConfig{Enabled: false})
	response, err := client.Health(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.Status != "disabled" {
		t.Fatalf("expected disabled status, got %q", response.Status)
	}
}
