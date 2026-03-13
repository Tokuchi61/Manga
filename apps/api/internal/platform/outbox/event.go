package outbox

import "time"

// Event represents transactional outbox payload.
type Event struct {
	EventID       string
	Topic         string
	Payload       []byte
	RequestID     string
	CorrelationID string
	CausationID   string
	Attempts      int
	CreatedAt     time.Time
}