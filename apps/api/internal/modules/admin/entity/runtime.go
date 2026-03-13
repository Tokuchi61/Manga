package entity

import "time"

// RuntimeConfig stores stage-21 runtime controls.
type RuntimeConfig struct {
	MaintenanceEnabled bool
	UpdatedAt          time.Time
}
