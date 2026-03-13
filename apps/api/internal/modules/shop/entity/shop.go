package entity

import "time"

const (
	ProductStateDraft    = "draft"
	ProductStateActive   = "active"
	ProductStateArchived = "archived"
)

const (
	OfferVisibilityVisible      = "visible"
	OfferVisibilityHidden       = "hidden"
	OfferVisibilityCampaignOnly = "campaign_only"
)

const (
	PurchaseStatusDeliveryPending  = "delivery_pending"
	PurchaseStatusCompleted        = "completed"
	PurchaseStatusRecoveryRequired = "recovery_required"
	PurchaseStatusBlocked          = "blocked"
)

const CurrencyMana = "mana"

// ProductDefinition stores stage-18 product catalog metadata.
type ProductDefinition struct {
	ProductID       string
	Name            string
	Category        string
	State           string
	InventoryItemID string
	SlotID          string
	SinglePurchase  bool
	VIPRequired     bool
	MinLevel        int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// OfferDefinition stores stage-18 offer metadata.
type OfferDefinition struct {
	OfferID         string
	ProductID       string
	Title           string
	Visibility      string
	PriceMana       int
	DiscountPercent int
	Active          bool
	StartsAt        *time.Time
	EndsAt          *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// PurchaseIntent stores stage-18 purchase intent state.
type PurchaseIntent struct {
	IntentID       string
	UserID         string
	ProductID      string
	OfferID        string
	FinalPriceMana int
	Currency       string
	Status         string
	RequestID      string
	CorrelationID  string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// RuntimeConfig stores stage-18 runtime controls.
type RuntimeConfig struct {
	CatalogEnabled  bool
	PurchaseEnabled bool
	CampaignEnabled bool
	UpdatedAt       time.Time
}
