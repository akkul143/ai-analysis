package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const runLockKey int64 = 0x4A4F4246494E4445

type RunLock struct {
	connection *pgxpool.Conn
	acquired   bool
}

func (store *Store) AcquireRunLock(ctx context.Context) (*RunLock, bool, error) {
	connection, err := store.pool.Acquire(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("acquire lock connection: %w", err)
	}
	var acquired bool
	if err := connection.QueryRow(ctx, "SELECT pg_try_advisory_lock($1)", runLockKey).Scan(&acquired); err != nil {
		connection.Release()
		return nil, false, fmt.Errorf("acquire advisory run lock: %w", err)
	}
	if !acquired {
		connection.Release()
		return nil, false, nil
	}
	return &RunLock{connection: connection, acquired: true}, true, nil
}

func (lock *RunLock) Release(ctx context.Context) error {
	if lock == nil || !lock.acquired {
		return nil
	}
	defer lock.connection.Release()
	var released bool
	if err := lock.connection.QueryRow(ctx, "SELECT pg_advisory_unlock($1)", runLockKey).Scan(&released); err != nil {
		return fmt.Errorf("release advisory run lock: %w", err)
	}
	lock.acquired = false
	if !released {
		return fmt.Errorf("release advisory run lock: PostgreSQL reported that the lock was not held")
	}
	return nil
}
