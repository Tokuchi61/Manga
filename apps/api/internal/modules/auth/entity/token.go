package entity

import (
	"time"

	"github.com/google/uuid"
)

// TokenType defines auth token categories.
type TokenType string

const (
	TokenTypeRefresh           TokenType = "refresh"
	TokenTypePasswordReset     TokenType = "password_reset"
	TokenTypeEmailVerification TokenType = "email_verification"
)

// Token stores hashed token material and lifecycle state.
type Token struct {
	ID           uuid.UUID
	CredentialID uuid.UUID
	SessionID    *uuid.UUID
	Type         TokenType
	TokenHash    string
	ExpiresAt    time.Time
	CreatedAt    time.Time
	ConsumedAt   *time.Time
}

func (t Token) IsConsumed() bool {
	return t.ConsumedAt != nil
}

func (t Token) IsExpired(now time.Time) bool {
	return now.After(t.ExpiresAt)
}
