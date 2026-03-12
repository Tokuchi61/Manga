package entity

import (
	"time"

	"github.com/google/uuid"
)

// Credential stores auth-owned identity secrets and security state.
type Credential struct {
	ID                            uuid.UUID
	Email                         string
	PasswordHash                  string
	EmailVerified                 bool
	Suspended                     bool
	Banned                        bool
	FailedLoginAttempts           int
	FailedLoginWindowStartedAt    *time.Time
	LoginCooldownUntil            *time.Time
	VerificationResendAvailableAt *time.Time
	CreatedAt                     time.Time
	UpdatedAt                     time.Time
}
