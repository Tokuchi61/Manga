package dto

import "time"

// GetPublicProfileRequest fetches public profile surface.
type GetPublicProfileRequest struct {
	UserID string `validate:"required,uuid4"`
}

// PublicProfileResponse exposes non-sensitive profile surface.
type PublicProfileResponse struct {
	UserID            string `json:"user_id"`
	Username          string `json:"username"`
	DisplayName       string `json:"display_name"`
	Bio               string `json:"bio"`
	AvatarURL         string `json:"avatar_url"`
	BannerURL         string `json:"banner_url"`
	ProfileVisibility string `json:"profile_visibility"`
	VIPBadgeVisible   bool   `json:"vip_badge_visible"`
}

// GetOwnProfileRequest fetches private owner surface.
type GetOwnProfileRequest struct {
	UserID   string `validate:"required,uuid4"`
	ViewerID string `validate:"required,uuid4"`
}

// OwnProfileResponse exposes full owner profile and membership state.
type OwnProfileResponse struct {
	UserID                      string     `json:"user_id"`
	CredentialID                string     `json:"credential_id"`
	Username                    string     `json:"username"`
	DisplayName                 string     `json:"display_name"`
	Bio                         string     `json:"bio"`
	AvatarURL                   string     `json:"avatar_url"`
	BannerURL                   string     `json:"banner_url"`
	ProfileVisibility           string     `json:"profile_visibility"`
	HistoryVisibilityPreference string     `json:"history_visibility_preference"`
	AccountState                string     `json:"account_state"`
	VIPActive                   bool       `json:"vip_active"`
	VIPFrozen                   bool       `json:"vip_frozen"`
	VIPStartedAt                *time.Time `json:"vip_started_at,omitempty"`
	VIPEndsAt                   *time.Time `json:"vip_ends_at,omitempty"`
	VIPFrozenAt                 *time.Time `json:"vip_frozen_at,omitempty"`
	VIPFreezeReason             string     `json:"vip_freeze_reason,omitempty"`
}

// UpdateProfileRequest updates owner profile fields.
type UpdateProfileRequest struct {
	UserID      string  `json:"-" validate:"required,uuid4"`
	ViewerID    string  `json:"-" validate:"required,uuid4"`
	DisplayName *string `json:"display_name,omitempty" validate:"omitempty,max=64"`
	Bio         *string `json:"bio,omitempty" validate:"omitempty,max=280"`
	AvatarURL   *string `json:"avatar_url,omitempty" validate:"omitempty,url,max=255"`
	BannerURL   *string `json:"banner_url,omitempty" validate:"omitempty,url,max=255"`
}
