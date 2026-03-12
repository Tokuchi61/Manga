package dto

import "time"

type AssignPermissionToRoleRequest struct {
	PermissionName string `json:"permission_name" validate:"required,min=3,max=128"`
}

type AssignPermissionToRoleResponse struct {
	Status         string `json:"status"`
	RoleID         string `json:"role_id"`
	PermissionName string `json:"permission_name"`
}

type AssignRoleToUserRequest struct {
	RoleName  string     `json:"role_name" validate:"required,min=3,max=64"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type AssignRoleToUserResponse struct {
	Status    string     `json:"status"`
	UserID    string     `json:"user_id"`
	RoleName  string     `json:"role_name"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type CreateTemporaryGrantRequest struct {
	PermissionName string    `json:"permission_name" validate:"required,min=3,max=128"`
	Reason         string    `json:"reason" validate:"required,min=3,max=256"`
	ExpiresAt      time.Time `json:"expires_at" validate:"required"`
}

type CreateTemporaryGrantResponse struct {
	Status         string    `json:"status"`
	GrantID        string    `json:"grant_id"`
	UserID         string    `json:"user_id"`
	PermissionName string    `json:"permission_name"`
	ExpiresAt      time.Time `json:"expires_at"`
}
