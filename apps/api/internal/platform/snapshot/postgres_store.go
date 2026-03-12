package snapshot

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresStore persists snapshots in PostgreSQL.
type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) EnsureSchema(ctx context.Context) error {
	if s == nil || s.pool == nil {
		return nil
	}
	_, err := s.pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS app_runtime_snapshots (
	module_name TEXT PRIMARY KEY,
	payload BYTEA NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)
`)
	if err != nil {
		return fmt.Errorf("snapshot schema init failed: %w", err)
	}
	return nil
}

func (s *PostgresStore) Load(ctx context.Context, module string) ([]byte, error) {
	if s == nil || s.pool == nil {
		return nil, nil
	}
	var payload []byte
	err := s.pool.QueryRow(ctx, `SELECT payload FROM app_runtime_snapshots WHERE module_name = $1`, module).Scan(&payload)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("snapshot load failed for %s: %w", module, err)
	}
	if len(payload) == 0 {
		return nil, nil
	}
	return payload, nil
}

func (s *PostgresStore) Save(ctx context.Context, module string, payload []byte) error {
	if s == nil || s.pool == nil {
		return nil
	}
	_, err := s.pool.Exec(ctx, `
INSERT INTO app_runtime_snapshots(module_name, payload, updated_at)
VALUES ($1, $2, NOW())
ON CONFLICT (module_name)
DO UPDATE SET payload = EXCLUDED.payload, updated_at = NOW()
`, module, payload)
	if err != nil {
		return fmt.Errorf("snapshot save failed for %s: %w", module, err)
	}
	return nil
}
