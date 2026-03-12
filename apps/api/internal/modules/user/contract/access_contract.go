package contract

import "time"

// AccessSignal exposes user status and visibility facts for access decisions.
type AccessSignal struct {
	UserID                      string     `json:"user_id"`
	AccountState                string     `json:"account_state"`
	ProfileVisibility           string     `json:"profile_visibility"`
	HistoryVisibilityPreference string     `json:"history_visibility_preference"`
	VIPActive                   bool       `json:"vip_active"`
	VIPFrozen                   bool       `json:"vip_frozen"`
	VIPEndsAt                   *time.Time `json:"vip_ends_at,omitempty"`
	UpdatedAt                   time.Time  `json:"updated_at"`
}
