package outbox

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresStore persists outbox events in PostgreSQL.
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
CREATE TABLE IF NOT EXISTS app_outbox_events (
	event_id UUID PRIMARY KEY,
	topic TEXT NOT NULL,
	payload JSONB NOT NULL,
	request_id TEXT NOT NULL,
	correlation_id TEXT,
	causation_id TEXT,
	status TEXT NOT NULL DEFAULT 'pending',
	attempts INTEGER NOT NULL DEFAULT 0,
	next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	leased_until TIMESTAMPTZ,
	last_error TEXT,
	dead_letter BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	processed_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_app_outbox_events_pending ON app_outbox_events(status, dead_letter, next_attempt_at);
`)
	if err != nil {
		return fmt.Errorf("outbox schema init failed: %w", err)
	}
	return nil
}

func (s *PostgresStore) Enqueue(ctx context.Context, event Event) error {
	if s == nil || s.pool == nil {
		return nil
	}
	eventID := strings.TrimSpace(event.EventID)
	if eventID == "" {
		eventID = uuid.NewString()
	}
	if strings.TrimSpace(event.Topic) == "" {
		return fmt.Errorf("outbox enqueue failed: topic is required")
	}
	if strings.TrimSpace(event.RequestID) == "" {
		return fmt.Errorf("outbox enqueue failed: request_id is required")
	}
	if len(event.Payload) == 0 {
		return fmt.Errorf("outbox enqueue failed: payload is required")
	}

	_, err := s.pool.Exec(ctx, `
INSERT INTO app_outbox_events(event_id, topic, payload, request_id, correlation_id, causation_id, status, attempts, next_attempt_at, created_at)
VALUES ($1::uuid, $2, $3::jsonb, $4, $5, $6, 'pending', 0, NOW(), NOW())
`, eventID, strings.TrimSpace(event.Topic), event.Payload, strings.TrimSpace(event.RequestID), strings.TrimSpace(event.CorrelationID), strings.TrimSpace(event.CausationID))
	if err != nil {
		return fmt.Errorf("outbox enqueue failed: %w", err)
	}
	return nil
}

func (s *PostgresStore) LeaseBatch(ctx context.Context, limit int, leaseDuration time.Duration) ([]Event, error) {
	if s == nil || s.pool == nil {
		return nil, nil
	}
	if limit <= 0 {
		limit = 20
	}
	if leaseDuration <= 0 {
		leaseDuration = 30 * time.Second
	}
	leaseSeconds := int(leaseDuration.Seconds())
	if leaseSeconds <= 0 {
		leaseSeconds = 30
	}

	rows, err := s.pool.Query(ctx, `
WITH candidates AS (
	SELECT event_id
	FROM app_outbox_events
	WHERE status = 'pending'
	  AND dead_letter = FALSE
	  AND next_attempt_at <= NOW()
	  AND (leased_until IS NULL OR leased_until <= NOW())
	ORDER BY created_at ASC
	FOR UPDATE SKIP LOCKED
	LIMIT $1
)
UPDATE app_outbox_events e
SET leased_until = NOW() + ($2 * interval '1 second'), attempts = attempts + 1
FROM candidates c
WHERE e.event_id = c.event_id
RETURNING e.event_id::text, e.topic, e.payload::text, e.request_id, COALESCE(e.correlation_id, ''), COALESCE(e.causation_id, ''), e.attempts, e.created_at
`, limit, leaseSeconds)
	if err != nil {
		return nil, fmt.Errorf("outbox lease failed: %w", err)
	}
	defer rows.Close()

	events := make([]Event, 0, limit)
	for rows.Next() {
		var event Event
		var payload string
		if scanErr := rows.Scan(&event.EventID, &event.Topic, &payload, &event.RequestID, &event.CorrelationID, &event.CausationID, &event.Attempts, &event.CreatedAt); scanErr != nil {
			return nil, fmt.Errorf("outbox lease scan failed: %w", scanErr)
		}
		event.Payload = []byte(payload)
		events = append(events, event)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("outbox lease rows failed: %w", rowsErr)
	}

	return events, nil
}

func (s *PostgresStore) MarkSucceeded(ctx context.Context, eventID string) error {
	if s == nil || s.pool == nil {
		return nil
	}
	_, err := s.pool.Exec(ctx, `
UPDATE app_outbox_events
SET status = 'processed', processed_at = NOW(), leased_until = NULL, last_error = NULL
WHERE event_id = $1::uuid
`, strings.TrimSpace(eventID))
	if err != nil {
		return fmt.Errorf("outbox mark succeeded failed: %w", err)
	}
	return nil
}

func (s *PostgresStore) MarkFailed(ctx context.Context, eventID string, nextAttemptAt time.Time, lastError string, deadLetter bool) error {
	if s == nil || s.pool == nil {
		return nil
	}
	status := "pending"
	if deadLetter {
		status = "dead"
	}
	_, err := s.pool.Exec(ctx, `
UPDATE app_outbox_events
SET status = $2,
	dead_letter = $3,
	next_attempt_at = $4,
	leased_until = NULL,
	last_error = $5
WHERE event_id = $1::uuid
`, strings.TrimSpace(eventID), status, deadLetter, nextAttemptAt.UTC(), strings.TrimSpace(lastError))
	if err != nil {
		return fmt.Errorf("outbox mark failed failed: %w", err)
	}
	return nil
}