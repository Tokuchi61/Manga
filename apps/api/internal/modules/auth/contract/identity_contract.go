package contract

import "time"

// VerifiedIdentity is the auth-owned contract passed to access module.
type VerifiedIdentity struct {
	CredentialID    string
	SessionID       string
	EmailVerified   bool
	Suspended       bool
	Banned          bool
	AuthenticatedAt time.Time
}

// SecuritySignal represents auth security outcomes for downstream policy consumers.
type SecuritySignal struct {
	CredentialID  string
	Signal        string
	OccurredAt    time.Time
	RequestID     string
	CorrelationID string
}
