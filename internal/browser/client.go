package browser

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/akkul143/ai-analysis/internal/config"
)

type Client struct {
	config config.BrowserConfig
}

type request struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}

type Response struct {
	ID      string         `json:"id"`
	Status  string         `json:"status"`
	Result  map[string]any `json:"result,omitempty"`
	Error   string         `json:"error,omitempty"`
}

func NewClient(cfg config.BrowserConfig) *Client {
	return &Client{config: cfg}
}

func (client *Client) Enabled() bool {
	return client.config.Enabled
}

func (client *Client) Health(ctx context.Context) (Response, error) {
	if !client.config.Enabled {
		return Response{Status: "disabled"}, nil
	}

	timeout := time.Duration(client.config.HealthTimeoutSeconds) * time.Second
	workerContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	payload, err := json.Marshal(request{ID: "phase1-health", Action: "health"})
	if err != nil {
		return Response{}, fmt.Errorf("encode browser worker request: %w", err)
	}

	command := exec.CommandContext(workerContext, client.config.Command, client.config.Args...)
	command.Stdin = bytes.NewReader(append(payload, '\n'))
	output, err := command.Output()
	if err != nil {
		return Response{}, fmt.Errorf("execute browser worker health check: %w", err)
	}

	line := strings.TrimSpace(string(output))
	var response Response
	if err := json.Unmarshal([]byte(line), &response); err != nil {
		return Response{}, fmt.Errorf("decode browser worker response %q: %w", line, err)
	}
	if response.Status != "ok" {
		return response, fmt.Errorf("browser worker health check returned %q: %s", response.Status, response.Error)
	}
	return response, nil
}
