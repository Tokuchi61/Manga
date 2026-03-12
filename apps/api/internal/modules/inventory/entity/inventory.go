package entity

import "time"

// ItemDefinition stores ownable inventory item metadata.
type ItemDefinition struct {
	ItemID     string
	ItemType   string
	Stackable  bool
	Equipable  bool
	Consumable bool
	MaxStack   int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// InventoryEntry stores user-owned inventory state.
type InventoryEntry struct {
	UserID         string
	ItemID         string
	Quantity       int
	Equipped       bool
	LastSourceType string
	LastSourceRef  string
	RequestID      string
	CorrelationID  string
	ExpiresAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// RuntimeConfig stores stage-15 inventory runtime controls.
type RuntimeConfig struct {
	ReadEnabled    bool
	ClaimEnabled   bool
	ConsumeEnabled bool
	EquipEnabled   bool
	UpdatedAt      time.Time
}
