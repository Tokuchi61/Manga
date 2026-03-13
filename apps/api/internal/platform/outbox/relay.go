package outbox

import (
	"context"
	"fmt"
	"time"
)

// Relay retries and publishes pending outbox events.
type Relay struct {
	store         Store
	publisher     Publisher
	batchSize     int
	leaseDuration time.Duration
	retryBase     time.Duration
	maxAttempts   int
	now           func() time.Time
}

func NewRelay(store Store, publisher Publisher) *Relay {
	return &Relay{
		store:         store,
		publisher:     publisher,
		batchSize:     20,
		leaseDuration: 30 * time.Second,
		retryBase:     2 * time.Second,
		maxAttempts:   8,
		now:           time.Now,
	}
}

func (r *Relay) Run(ctx context.Context, interval time.Duration) {
	if r == nil || r.store == nil || r.publisher == nil {
		return
	}
	if interval <= 0 {
		interval = 2 * time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	r.tick(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.tick(ctx)
		}
	}
}

func (r *Relay) tick(ctx context.Context) {
	events, err := r.store.LeaseBatch(ctx, r.batchSize, r.leaseDuration)
	if err != nil || len(events) == 0 {
		return
	}

	for _, event := range events {
		publishErr := r.publisher.Publish(ctx, event)
		if publishErr == nil {
			_ = r.store.MarkSucceeded(ctx, event.EventID)
			continue
		}

		deadLetter := event.Attempts >= r.maxAttempts
		nextAttempt := r.nextAttemptAt(event.Attempts)
		_ = r.store.MarkFailed(ctx, event.EventID, nextAttempt, fmt.Sprintf("publish_failed: %v", publishErr), deadLetter)
	}
}

func (r *Relay) nextAttemptAt(attempts int) time.Time {
	if r.now == nil {
		r.now = time.Now
	}
	if attempts < 1 {
		attempts = 1
	}
	if r.retryBase <= 0 {
		r.retryBase = 2 * time.Second
	}

	delay := r.retryBase
	for i := 1; i < attempts; i++ {
		delay *= 2
		if delay > 5*time.Minute {
			delay = 5 * time.Minute
			break
		}
	}
	return r.now().UTC().Add(delay)
}