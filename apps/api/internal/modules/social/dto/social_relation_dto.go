package dto

import "time"

// UpdateBlockRequest updates block relation state.
type UpdateBlockRequest struct {
	ActorUserID  string `json:"-" validate:"required,uuid4"`
	TargetUserID string `json:"-" validate:"required,uuid4"`
	Enabled      bool   `json:"enabled"`
}

// UpdateMuteRequest updates mute relation state.
type UpdateMuteRequest struct {
	ActorUserID  string `json:"-" validate:"required,uuid4"`
	TargetUserID string `json:"-" validate:"required,uuid4"`
	Enabled      bool   `json:"enabled"`
}

// UpdateRestrictRequest updates restrict relation state.
type UpdateRestrictRequest struct {
	ActorUserID  string `json:"-" validate:"required,uuid4"`
	TargetUserID string `json:"-" validate:"required,uuid4"`
	Enabled      bool   `json:"enabled"`
}

// ListRelationsRequest resolves own relation list by type.
type ListRelationsRequest struct {
	ActorUserID  string `json:"-" validate:"required,uuid4"`
	RelationType string `json:"-" validate:"required,oneof=blocked muted restricted"`
}

// RelationUpdateResponse returns relation update result.
type RelationUpdateResponse struct {
	ActorUserID  string    `json:"actor_user_id"`
	TargetUserID string    `json:"target_user_id"`
	RelationType string    `json:"relation_type"`
	Enabled      bool      `json:"enabled"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RelationItemResponse is relation list payload.
type RelationItemResponse struct {
	TargetUserID string    `json:"target_user_id"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ListRelationsResponse wraps relation list payload.
type ListRelationsResponse struct {
	Items []RelationItemResponse `json:"items"`
	Count int                    `json:"count"`
}
