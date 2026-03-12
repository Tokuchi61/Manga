package repository

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/entity"
)

var (
	ErrNotFound = errors.New("inventory_repository_not_found")
	ErrConflict = errors.New("inventory_repository_conflict")
)

// Store defines inventory persistence boundary.
type Store interface {
	UpsertItemDefinition(ctx context.Context, definition entity.ItemDefinition) error
	GetItemDefinition(ctx context.Context, itemID string) (entity.ItemDefinition, error)
	ListItemDefinitions(ctx context.Context, itemType string) ([]entity.ItemDefinition, error)

	GetInventoryEntry(ctx context.Context, userID string, itemID string) (entity.InventoryEntry, error)
	UpsertInventoryEntry(ctx context.Context, entry entity.InventoryEntry) error
	DeleteInventoryEntry(ctx context.Context, userID string, itemID string) error
	ListInventoryEntries(ctx context.Context, userID string, itemType string, equippedOnly bool, sortBy string, limit int, offset int) ([]entity.InventoryEntry, error)

	GetGrantByDedup(ctx context.Context, dedupKey string) (entity.InventoryEntry, error)
	PutGrantDedup(ctx context.Context, dedupKey string, entry entity.InventoryEntry) error

	GetRuntimeConfig(ctx context.Context) (entity.RuntimeConfig, error)
	UpdateRuntimeConfig(ctx context.Context, cfg entity.RuntimeConfig) error
}
