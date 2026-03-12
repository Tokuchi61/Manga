package entity

import (
	"time"

	"github.com/google/uuid"
)

type RolePermission struct {
	RoleID       uuid.UUID
	PermissionID uuid.UUID
	CreatedAt    time.Time
}

type UserRole struct {
	UserID    string
	RoleID    uuid.UUID
	ExpiresAt *time.Time
	CreatedAt time.Time
}

func (u UserRole) IsActive(at time.Time) bool {
	if u.ExpiresAt == nil {
		return true
	}
	return u.ExpiresAt.After(at)
}

type TemporaryGrant struct {
	ID           uuid.UUID
	UserID       string
	PermissionID uuid.UUID
	Reason       string
	ExpiresAt    time.Time
	RevokedAt    *time.Time
	CreatedAt    time.Time
}

func (g TemporaryGrant) IsActive(at time.Time) bool {
	if g.RevokedAt != nil {
		return false
	}
	return g.ExpiresAt.After(at)
}
