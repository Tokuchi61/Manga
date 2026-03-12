package entity

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID           uuid.UUID
	Name         string
	Module       string
	Surface      string
	Action       string
	AudienceKind string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
