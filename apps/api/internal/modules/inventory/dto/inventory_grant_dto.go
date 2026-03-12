package dto

import "time"

// AdminGrantItemRequest grants item to target inventory.
type AdminGrantItemRequest struct {
	TargetUserID  string     `json:"target_user_id" validate:"required,uuid4"`
	ItemID        string     `json:"item_id" validate:"required,max=128"`
	Quantity      int        `json:"quantity" validate:"required,min=1,max=10000"`
	SourceType    string     `json:"source_type" validate:"required,max=64"`
	SourceRef     string     `json:"source_ref,omitempty" validate:"omitempty,max=128"`
	RequestID     string     `json:"request_id,omitempty" validate:"omitempty,max=128"`
	CorrelationID string     `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
}

// AdminGrantItemResponse returns admin grant result.
type AdminGrantItemResponse struct {
	Status    string                     `json:"status"`
	Created   bool                       `json:"created"`
	Inventory InventoryEntryItemResponse `json:"inventory"`
}

// AdminRevokeItemRequest revokes quantity from target inventory.
type AdminRevokeItemRequest struct {
	TargetUserID string `json:"target_user_id" validate:"required,uuid4"`
	ItemID       string `json:"item_id" validate:"required,max=128"`
	Quantity     int    `json:"quantity" validate:"required,min=1,max=10000"`
}

// AdminRevokeItemResponse returns revoke operation result.
type AdminRevokeItemResponse struct {
	Status            string    `json:"status"`
	ItemID            string    `json:"item_id"`
	RemainingQuantity int       `json:"remaining_quantity"`
	RevokedAt         time.Time `json:"revoked_at"`
}
