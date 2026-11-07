package run

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/akkul143/ai-analysis/internal/store/postgres"
)

var ErrRunAlreadyActive = errors.New("another job-finder run already holds the PostgreSQL advisory lock")

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
) (summary Summary, returnError error) {
	lock, acquired, err := service.store.AcquireRunLock(ctx)
	if err != nil {
		return Summary{}, err
	}
	if !acquired {
		return Summary{}, ErrRunAlreadyActive
	}
	defer func() {
		cleanupContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := lock.Release(cleanupContext); err != nil {
			returnError = errors.Join(returnError, err)
		}
	}()

	id, err := newID()
	if err != nil {
		return Summary{}, err
	}
	startedAt := time.Now().UTC()
	summary = Summary{ID: id, Mode: mode, Status: "running", StartedAt: startedAt}
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
