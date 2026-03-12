package entity

import (
	"time"

	"github.com/google/uuid"
)

// SecurityEvent captures auth security and audit-relevant events.
type SecurityEvent struct {
	ID            uuid.UUID
	CredentialID  *uuid.UUID
	ActorID       *uuid.UUID
	TargetID      *uuid.UUID
	Action        string
	Result        string
	Reason        string
	RequestID     string
	CorrelationID string
	Device        string
	IP            string
	CreatedAt     time.Time
}
