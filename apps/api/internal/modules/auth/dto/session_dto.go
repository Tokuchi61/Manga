package dto

import "time"

// ListSessionsRequest represents session list input.
type ListSessionsRequest struct {
	CredentialID string `json:"-" validate:"required,uuid4"`
}

// RevokeCurrentSessionRequest revokes the active session.
type RevokeCurrentSessionRequest struct {
	CredentialID string `json:"-" validate:"required,uuid4"`
	SessionID    string `json:"session_id" validate:"required,uuid4"`
}

// RevokeOtherSessionsRequest revokes all sessions except current.
type RevokeOtherSessionsRequest struct {
	CredentialID string `json:"-" validate:"required,uuid4"`
	SessionID    string `json:"session_id" validate:"required,uuid4"`
}

// RevokeAllSessionsRequest revokes all sessions.
type RevokeAllSessionsRequest struct {
	CredentialID string `json:"-" validate:"required,uuid4"`
}

// LogoutRequest maps to revoke current session.
type LogoutRequest struct {
	CredentialID string `json:"-" validate:"required,uuid4"`
	SessionID    string `json:"session_id" validate:"required,uuid4"`
}

// SessionInfo is a serializable session view.
type SessionInfo struct {
	SessionID  string     `json:"session_id"`
	Device     string     `json:"device"`
	IP         string     `json:"ip"`
	CreatedAt  time.Time  `json:"created_at"`
	LastSeenAt time.Time  `json:"last_seen_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
}

// ListSessionsResponse returns active and revoked sessions owned by auth.
type ListSessionsResponse struct {
	Sessions []SessionInfo `json:"sessions"`
}
