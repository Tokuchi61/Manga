package dto

import "time"

type CreatePermissionRequest struct {
	Name         string `json:"name" validate:"required,min=3,max=128"`
	Module       string `json:"module" validate:"required,min=2,max=64"`
	Surface      string `json:"surface" validate:"required,min=2,max=64"`
	Action       string `json:"action" validate:"required,min=2,max=64"`
	AudienceKind string `json:"audience_kind" validate:"required"`
}

type CreatePermissionResponse struct {
	PermissionID string    `json:"permission_id"`
	Name         string    `json:"name"`
	Module       string    `json:"module"`
	Surface      string    `json:"surface"`
	Action       string    `json:"action"`
	AudienceKind string    `json:"audience_kind"`
	CreatedAt    time.Time `json:"created_at"`
}

type ListCanonicalPermissionsResponse struct {
	Permissions []string `json:"permissions"`
}
