package dto

// UpdateProfileVisibilityRequest updates global profile visibility surface.
type UpdateProfileVisibilityRequest struct {
	UserID            string `json:"-" validate:"required,uuid4"`
	ViewerID          string `json:"-" validate:"required,uuid4"`
	ProfileVisibility string `json:"profile_visibility" validate:"required,oneof=public private"`
}

// UpdateHistoryVisibilityRequest updates global history visibility preference.
type UpdateHistoryVisibilityRequest struct {
	UserID                      string `json:"-" validate:"required,uuid4"`
	ViewerID                    string `json:"-" validate:"required,uuid4"`
	HistoryVisibilityPreference string `json:"history_visibility_preference" validate:"required,oneof=public private"`
}

// VisibilityResponse returns updated visibility state.
type VisibilityResponse struct {
	Status                      string `json:"status"`
	ProfileVisibility           string `json:"profile_visibility,omitempty"`
	HistoryVisibilityPreference string `json:"history_visibility_preference,omitempty"`
}
