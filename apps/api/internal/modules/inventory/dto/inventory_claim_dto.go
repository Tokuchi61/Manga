package dto

import "time"

// ClaimInventoryItemRequest claims reward into actor inventory.
type ClaimInventoryItemRequest struct {
	ActorUserID   string     `json:"-" validate:"required,uuid4"`
	ItemID        string     `json:"item_id" validate:"required,max=128"`
	Quantity      int        `json:"quantity" validate:"required,min=1,max=10000"`
	SourceType    string     `json:"source_type" validate:"required,max=64"`
	SourceRef     string     `json:"source_ref,omitempty" validate:"omitempty,max=128"`
	RequestID     string     `json:"request_id,omitempty" validate:"omitempty,max=128"`
	CorrelationID string     `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
}

// ClaimInventoryItemResponse returns claim operation result.
type ClaimInventoryItemResponse struct {
	Status    string                     `json:"status"`
	Created   bool                       `json:"created"`
	Inventory InventoryEntryItemResponse `json:"inventory"`
}

// ConsumeInventoryItemRequest consumes quantity from actor inventory entry.
type ConsumeInventoryItemRequest struct {
	ActorUserID   string `json:"-" validate:"required,uuid4"`
	ItemID        string `json:"-" validate:"required,max=128"`
	Quantity      int    `json:"quantity" validate:"required,min=1,max=10000"`
	RequestID     string `json:"request_id,omitempty" validate:"omitempty,max=128"`
	CorrelationID string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
}

// ConsumeInventoryItemResponse returns consume operation result.
type ConsumeInventoryItemResponse struct {
	Status            string    `json:"status"`
	ItemID            string    `json:"item_id"`
	RemainingQuantity int       `json:"remaining_quantity"`
	ConsumedAt        time.Time `json:"consumed_at"`
}

// UpdateEquipStateRequest toggles actor equipped state for an item.
type UpdateEquipStateRequest struct {
	ActorUserID string `json:"-" validate:"required,uuid4"`
	ItemID      string `json:"-" validate:"required,max=128"`
	Enabled     bool   `json:"enabled"`
}

// UpdateEquipStateResponse returns equip operation result.
type UpdateEquipStateResponse struct {
	Status    string    `json:"status"`
	ItemID    string    `json:"item_id"`
	Equipped  bool      `json:"equipped"`
	UpdatedAt time.Time `json:"updated_at"`
}
