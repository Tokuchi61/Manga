package dto

import "time"

type CreateRoleRequest struct {
	Name         string `json:"name" validate:"required,min=3,max=64"`
	Priority     int    `json:"priority" validate:"gte=0,lte=1000"`
	IsDefault    bool   `json:"is_default"`
	IsSuperAdmin bool   `json:"is_super_admin"`
}

type CreateRoleResponse struct {
	RoleID       string    `json:"role_id"`
	Name         string    `json:"name"`
	Priority     int       `json:"priority"`
	IsDefault    bool      `json:"is_default"`
	IsSuperAdmin bool      `json:"is_super_admin"`
	CreatedAt    time.Time `json:"created_at"`
}
