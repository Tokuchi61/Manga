package dto

import "time"

// GetPreferenceRequest resolves own preference state.
type GetPreferenceRequest struct {
	UserID string `json:"-" validate:"required,uuid4"`
}

// UpdatePreferenceRequest updates own preference state.
type UpdatePreferenceRequest struct {
	UserID            string   `json:"-" validate:"required,uuid4"`
	MutedCategories   []string `json:"muted_categories,omitempty" validate:"omitempty,max=20,dive,max=64"`
	QuietHoursEnabled bool     `json:"quiet_hours_enabled"`
	QuietHoursStart   int      `json:"quiet_hours_start" validate:"min=0,max=23"`
	QuietHoursEnd     int      `json:"quiet_hours_end" validate:"min=0,max=23"`
	InAppEnabled      bool     `json:"in_app_enabled"`
	EmailEnabled      bool     `json:"email_enabled"`
	PushEnabled       bool     `json:"push_enabled"`
	DigestEnabled     bool     `json:"digest_enabled"`
}

// PreferenceResponse is own preference payload.
type PreferenceResponse struct {
	UserID            string    `json:"user_id"`
	MutedCategories   []string  `json:"muted_categories"`
	QuietHoursEnabled bool      `json:"quiet_hours_enabled"`
	QuietHoursStart   int       `json:"quiet_hours_start"`
	QuietHoursEnd     int       `json:"quiet_hours_end"`
	InAppEnabled      bool      `json:"in_app_enabled"`
	EmailEnabled      bool      `json:"email_enabled"`
	PushEnabled       bool      `json:"push_enabled"`
	DigestEnabled     bool      `json:"digest_enabled"`
	UpdatedAt         time.Time `json:"updated_at"`
}
