package dto

import "time"

// ListManaPackagesRequest resolves package listing for actor surface.
type ListManaPackagesRequest struct {
	ActorUserID string `json:"-" validate:"required,max=64"`
}

// ManaPackageResponse returns package payload.
type ManaPackageResponse struct {
	PackageID     string    `json:"package_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	ManaAmount    int       `json:"mana_amount"`
	PriceAmount   int       `json:"price_amount"`
	PriceCurrency string    `json:"price_currency"`
	Active        bool      `json:"active"`
	DisplayOrder  int       `json:"display_order"`
	Provider      string    `json:"provider"`
	ProviderSKU   string    `json:"provider_sku"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ListManaPackagesResponse wraps package list payload.
type ListManaPackagesResponse struct {
	Items []ManaPackageResponse `json:"items"`
	Count int                   `json:"count"`
}
