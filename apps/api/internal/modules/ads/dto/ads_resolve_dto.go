package dto

import "time"

// ResolvePlacementsRequest resolves placement/campaign payload for ads surface.
type ResolvePlacementsRequest struct {
	Surface    string `json:"-" validate:"required,oneof=home listing manga chapter"`
	TargetType string `json:"-" validate:"omitempty,oneof=none manga chapter"`
	TargetID   string `json:"-" validate:"omitempty,max=64"`
	SessionID  string `json:"-" validate:"omitempty,max=128"`
	NoAds      bool   `json:"-"`
	Limit      int    `json:"-" validate:"omitempty,min=1,max=20"`
}

// CampaignPreviewResponse returns campaign preview payload.
type CampaignPreviewResponse struct {
	CampaignID  string     `json:"campaign_id"`
	Name        string     `json:"name"`
	State       string     `json:"state"`
	CreativeURL string     `json:"creative_url"`
	ClickURL    string     `json:"click_url"`
	Weight      int        `json:"weight"`
	StartsAt    *time.Time `json:"starts_at,omitempty"`
	EndsAt      *time.Time `json:"ends_at,omitempty"`
}

// PlacementResolveResponse returns resolved placement payload.
type PlacementResolveResponse struct {
	PlacementID  string                    `json:"placement_id"`
	Surface      string                    `json:"surface"`
	Priority     int                       `json:"priority"`
	FrequencyCap int                       `json:"frequency_cap"`
	Campaigns    []CampaignPreviewResponse `json:"campaigns"`
}

// ResolvePlacementsResponse wraps resolved placement list.
type ResolvePlacementsResponse struct {
	Items []PlacementResolveResponse `json:"items"`
	Count int                        `json:"count"`
}
