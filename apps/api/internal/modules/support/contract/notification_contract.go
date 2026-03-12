package contract

import "time"

const (
	EventSupportCreated  = "support.created"
	EventSupportReplied  = "support.replied"
	EventSupportResolved = "support.resolved"
)

// NotificationSignal is support -> notification payload.
type NotificationSignal struct {
	Event           string
	SupportID       string
	RequesterUserID string
	Status          string
	OccurredAt      time.Time
	RequestID       string
	CorrelationID   string
}
