package dto

import "time"

// UpdateMaintenanceStateRequest updates site maintenance runtime state.
type UpdateMaintenanceStateRequest struct {
	RequestID         string `json:"request_id" validate:"required,max=128"`
	CorrelationID     string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
	Enabled           bool   `json:"enabled"`
	Reason            string `json:"reason" validate:"required,max=512"`
	RiskLevel         string `json:"risk_level" validate:"required,oneof=low medium high critical"`
	DoubleConfirmed   bool   `json:"double_confirmed"`
	ConfirmationToken string `json:"confirmation_token,omitempty" validate:"omitempty,max=256"`
}

// RuntimeConfigResponse is stage-21 runtime control payload.
type RuntimeConfigResponse struct {
	MaintenanceEnabled bool      `json:"maintenance_enabled"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// SettingChangeResponse returns maintenance update result payload.
type SettingChangeResponse struct {
	Status    string    `json:"status"`
	ActionID  string    `json:"action_id"`
	Key       string    `json:"key"`
	Enabled   bool      `json:"enabled"`
	Event     string    `json:"event"`
	UpdatedAt time.Time `json:"updated_at"`
}
