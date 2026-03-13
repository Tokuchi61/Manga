package outbox

import "context"

// Publisher dispatches leased outbox events to downstream transport.
type Publisher interface {
	Publish(ctx context.Context, event Event) error
}