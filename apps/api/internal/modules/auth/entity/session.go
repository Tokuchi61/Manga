package entity

import (
	"time"

	"github.com/google/uuid"
)

// Session stores auth-owned session lifecycle records.
type Session struct {
	ID           uuid.UUID
	CredentialID uuid.UUID
	Device       string
	IP           string
	CreatedAt    time.Time
	LastSeenAt   time.Time
	RevokedAt    *time.Time
}

func (s Session) IsRevoked() bool {
	return s.RevokedAt != nil
}
