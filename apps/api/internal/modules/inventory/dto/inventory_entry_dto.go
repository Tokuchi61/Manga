package dto

import "time"

// ListInventoryEntriesRequest resolves actor inventory list.
type ListInventoryEntriesRequest struct {
	ActorUserID  string `json:"-" validate:"required,uuid4"`
	ItemType     string `json:"-" validate:"omitempty,max=64"`
	EquippedOnly bool   `json:"-"`
	SortBy       string `json:"-" validate:"omitempty,oneof=newest oldest"`
	Limit        int    `json:"-" validate:"omitempty,min=1,max=100"`
	Offset       int    `json:"-" validate:"omitempty,min=0,max=10000"`
}

// GetInventoryItemDetailRequest resolves own item detail.
type GetInventoryItemDetailRequest struct {
	ActorUserID string `json:"-" validate:"required,uuid4"`
	ItemID      string `json:"-" validate:"required,max=128"`
}

// InventoryEntryItemResponse is inventory entry payload.
type InventoryEntryItemResponse struct {
	ItemID         string     `json:"item_id"`
	ItemType       string     `json:"item_type"`
	Quantity       int        `json:"quantity"`
	Equipped       bool       `json:"equipped"`
	Stackable      bool       `json:"stackable"`
	Equipable      bool       `json:"equipable"`
	Consumable     bool       `json:"consumable"`
	MaxStack       int        `json:"max_stack"`
	LastSourceType string     `json:"last_source_type,omitempty"`
	LastSourceRef  string     `json:"last_source_ref,omitempty"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ListInventoryEntriesResponse wraps actor inventory listing.
type ListInventoryEntriesResponse struct {
	Items []InventoryEntryItemResponse `json:"items"`
	Count int                          `json:"count"`
}

// InventoryItemDetailResponse returns own item detail.
type InventoryItemDetailResponse struct {
	ItemID         string     `json:"item_id"`
	ItemType       string     `json:"item_type"`
	Quantity       int        `json:"quantity"`
	Equipped       bool       `json:"equipped"`
	Stackable      bool       `json:"stackable"`
	Equipable      bool       `json:"equipable"`
	Consumable     bool       `json:"consumable"`
	MaxStack       int        `json:"max_stack"`
	LastSourceType string     `json:"last_source_type,omitempty"`
	LastSourceRef  string     `json:"last_source_ref,omitempty"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
