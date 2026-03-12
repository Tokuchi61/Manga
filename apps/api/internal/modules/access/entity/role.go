package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID           uuid.UUID
	Name         string
	Priority     int
	IsDefault    bool
	IsSuperAdmin bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
