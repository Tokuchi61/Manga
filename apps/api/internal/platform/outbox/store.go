package outbox

import (
	"context"
	"time"
)

// Store defines durable outbox persistence.
type Store interface {
	EnsureSchema(ctx context.Context) error
	Enqueue(ctx context.Context, event Event) error
	LeaseBatch(ctx context.Context, limit int, leaseDuration time.Duration) ([]Event, error)
	MarkSucceeded(ctx context.Context, eventID string) error
	MarkFailed(ctx context.Context, eventID string, nextAttemptAt time.Time, lastError string, deadLetter bool) error
}