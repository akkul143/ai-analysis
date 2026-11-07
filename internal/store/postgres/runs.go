package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type RunRecord struct {
	ID        string
	Mode      string
	Status    string
	StartedAt time.Time
	Metadata  map[string]any
}

func (store *Store) StartRun(ctx context.Context, record RunRecord) error {
	metadata, err := json.Marshal(record.Metadata)
	if err != nil {
		return fmt.Errorf("encode run metadata: %w", err)
	}
	_, err = store.pool.Exec(ctx, `
		INSERT INTO runs (id, mode, status, started_at, metadata)
		VALUES ($1, $2, $3, $4, $5::jsonb)
	`, record.ID, record.Mode, record.Status, record.StartedAt, string(metadata))
	if err != nil {
		return fmt.Errorf("insert run record: %w", err)
	}
	return nil
}

func (store *Store) FinishRun(ctx context.Context, id string, status string, finishedAt time.Time, runError error) error {
	var errorMessage *string
	if runError != nil {
		message := runError.Error()
		errorMessage = &message
	}
	result, err := store.pool.Exec(ctx, `
		UPDATE runs
		SET status = $2, finished_at = $3, error_message = $4
		WHERE id = $1
	`, id, status, finishedAt, errorMessage)
	if err != nil {
		return fmt.Errorf("finish run record: %w", err)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("finish run record: expected one updated row, got %d", result.RowsAffected())
	}
	return nil
}
