package dto

import "time"

// UpdateCatalogStateRequest updates catalog read surface runtime state.
type UpdateCatalogStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdatePurchaseStateRequest updates purchase write surface runtime state.
type UpdatePurchaseStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateCampaignStateRequest updates campaign surface runtime state.
type UpdateCampaignStateRequest struct {
	Enabled bool `json:"enabled"`
}

// RuntimeConfigResponse is shop runtime control payload.
type RuntimeConfigResponse struct {
	CatalogEnabled  bool      `json:"catalog_enabled"`
	PurchaseEnabled bool      `json:"purchase_enabled"`
	CampaignEnabled bool      `json:"campaign_enabled"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// UpsertProductDefinitionRequest creates or updates product definition.
type UpsertProductDefinitionRequest struct {
	ProductID       string `json:"product_id" validate:"required,max=64"`
	Name            string `json:"name" validate:"required,max=160"`
	Category        string `json:"category" validate:"required,max=64"`
	State           string `json:"state" validate:"required,oneof=draft active archived"`
	InventoryItemID string `json:"inventory_item_id,omitempty" validate:"omitempty,max=128"`
	SlotID          string `json:"slot_id,omitempty" validate:"omitempty,max=64"`
	SinglePurchase  bool   `json:"single_purchase"`
	VIPRequired     bool   `json:"vip_required"`
	MinLevel        int    `json:"min_level" validate:"omitempty,min=0,max=1000"`
}

// ProductDefinitionResponse returns product definition payload.
type ProductDefinitionResponse struct {
	ProductID       string    `json:"product_id"`
	Name            string    `json:"name"`
	Category        string    `json:"category"`
	State           string    `json:"state"`
	InventoryItemID string    `json:"inventory_item_id,omitempty"`
	SlotID          string    `json:"slot_id,omitempty"`
	SinglePurchase  bool      `json:"single_purchase"`
	VIPRequired     bool      `json:"vip_required"`
	MinLevel        int       `json:"min_level"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ListProductDefinitionsRequest resolves product definitions.
type ListProductDefinitionsRequest struct {
	State  string `json:"-" validate:"omitempty,oneof=draft active archived"`
	Limit  int    `json:"-" validate:"omitempty,min=1,max=100"`
	Offset int    `json:"-" validate:"omitempty,min=0,max=10000"`
}

// ListProductDefinitionsResponse wraps product definition list payload.
type ListProductDefinitionsResponse struct {
	Items []ProductDefinitionResponse `json:"items"`
	Count int                         `json:"count"`
}

// UpsertOfferDefinitionRequest creates or updates offer definition.
type UpsertOfferDefinitionRequest struct {
	OfferID         string     `json:"offer_id" validate:"required,max=64"`
	ProductID       string     `json:"product_id" validate:"required,max=64"`
	Title           string     `json:"title" validate:"required,max=160"`
	Visibility      string     `json:"visibility" validate:"required,oneof=visible hidden campaign_only"`
	PriceMana       int        `json:"price_mana" validate:"required,min=1,max=1000000"`
	DiscountPercent int        `json:"discount_percent" validate:"omitempty,min=0,max=100"`
	Active          bool       `json:"active"`
	StartsAt        *time.Time `json:"starts_at,omitempty"`
	EndsAt          *time.Time `json:"ends_at,omitempty"`
}

// OfferDefinitionResponse returns offer definition payload.
type OfferDefinitionResponse struct {
	OfferID         string     `json:"offer_id"`
	ProductID       string     `json:"product_id"`
	Title           string     `json:"title"`
	Visibility      string     `json:"visibility"`
	PriceMana       int        `json:"price_mana"`
	DiscountPercent int        `json:"discount_percent"`
	FinalPriceMana  int        `json:"final_price_mana"`
	Active          bool       `json:"active"`
	StartsAt        *time.Time `json:"starts_at,omitempty"`
	EndsAt          *time.Time `json:"ends_at,omitempty"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ListOfferDefinitionsRequest resolves offer definitions.
type ListOfferDefinitionsRequest struct {
	ProductID  string `json:"-" validate:"omitempty,max=64"`
	Visibility string `json:"-" validate:"omitempty,oneof=visible hidden campaign_only"`
	ActiveOnly bool   `json:"-"`
	Limit      int    `json:"-" validate:"omitempty,min=1,max=100"`
	Offset     int    `json:"-" validate:"omitempty,min=0,max=10000"`
}

// ListOfferDefinitionsResponse wraps offer definition list payload.
type ListOfferDefinitionsResponse struct {
	Items []OfferDefinitionResponse `json:"items"`
	Count int                       `json:"count"`
}
