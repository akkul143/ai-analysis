package run

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/akkul143/ai-analysis/internal/store/postgres"
)

type Service struct {
	store  *postgres.Store
	logger *slog.Logger
}

type Summary struct {
	ID          string
	Mode        string
	Status      string
	StartedAt   time.Time
	CompletedAt time.Time
}

func NewService(store *postgres.Store, logger *slog.Logger) *Service {
	return &Service{store: store, logger: logger}
}

func (service *Service) Execute(
	ctx context.Context,
	mode string,
	metadata map[string]any,
	work func(context.Context, string) error,
) (Summary, error) {
	id, err := newID()
	if err != nil {
		return Summary{}, err
	}
	startedAt := time.Now().UTC()
	summary := Summary{ID: id, Mode: mode, Status: "running", StartedAt: startedAt}
	if err := service.store.StartRun(ctx, postgres.RunRecord{
		ID:        id,
		Mode:      mode,
		Status:    "running",
		StartedAt: startedAt,
		Metadata:  metadata,
	}); err != nil {
		return Summary{}, err
	}

	service.logger.InfoContext(ctx, "run started", "run_id", id, "mode", mode)
	workError := work(ctx, id)
	finishedAt := time.Now().UTC()
	status := "completed"
	if workError != nil {
		status = "failed"
	}

	finishContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := service.store.FinishRun(finishContext, id, status, finishedAt, workError); err != nil {
		return summary, errors.Join(workError, err)
	}

	summary.Status = status
	summary.CompletedAt = finishedAt
	if workError != nil {
		service.logger.ErrorContext(ctx, "run failed", "run_id", id, "mode", mode, "error", workError)
		return summary, fmt.Errorf("run %s failed: %w", id, workError)
	}
	service.logger.InfoContext(ctx, "run completed", "run_id", id, "mode", mode, "duration_ms", finishedAt.Sub(startedAt).Milliseconds())
	return summary, nil
}
