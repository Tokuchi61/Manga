package dto

import "time"

// ListCatalogRequest resolves actor catalog listing.
type ListCatalogRequest struct {
	ActorUserID     string `json:"-" validate:"required,uuid4"`
	IncludeCampaign bool   `json:"-"`
	Limit           int    `json:"-" validate:"omitempty,min=1,max=100"`
	Offset          int    `json:"-" validate:"omitempty,min=0,max=10000"`
}

// CatalogItemResponse is public shop catalog item payload.
type CatalogItemResponse struct {
	ProductID       string `json:"product_id"`
	Name            string `json:"name"`
	Category        string `json:"category"`
	InventoryItemID string `json:"inventory_item_id,omitempty"`
	OfferID         string `json:"offer_id"`
	OfferTitle      string `json:"offer_title"`
	Visibility      string `json:"visibility"`
	PriceMana       int    `json:"price_mana"`
	DiscountPercent int    `json:"discount_percent"`
	FinalPriceMana  int    `json:"final_price_mana"`
}

// ListCatalogResponse wraps catalog listing payload.
type ListCatalogResponse struct {
	Items []CatalogItemResponse `json:"items"`
	Count int                   `json:"count"`
}

// GetCatalogItemRequest resolves actor catalog item detail.
type GetCatalogItemRequest struct {
	ActorUserID     string `json:"-" validate:"required,uuid4"`
	ProductID       string `json:"-" validate:"required,max=64"`
	IncludeCampaign bool   `json:"-"`
}

// CatalogItemDetailResponse is public shop product detail payload.
type CatalogItemDetailResponse struct {
	ProductID       string                 `json:"product_id"`
	Name            string                 `json:"name"`
	Category        string                 `json:"category"`
	State           string                 `json:"state"`
	InventoryItemID string                 `json:"inventory_item_id,omitempty"`
	SlotID          string                 `json:"slot_id,omitempty"`
	SinglePurchase  bool                   `json:"single_purchase"`
	Offers          []OfferPreviewResponse `json:"offers"`
	Count           int                    `json:"count"`
}

// OfferPreviewResponse returns public offer detail in catalog item response.
type OfferPreviewResponse struct {
	OfferID         string     `json:"offer_id"`
	Title           string     `json:"title"`
	Visibility      string     `json:"visibility"`
	PriceMana       int        `json:"price_mana"`
	DiscountPercent int        `json:"discount_percent"`
	FinalPriceMana  int        `json:"final_price_mana"`
	StartsAt        *time.Time `json:"starts_at,omitempty"`
	EndsAt          *time.Time `json:"ends_at,omitempty"`
}
