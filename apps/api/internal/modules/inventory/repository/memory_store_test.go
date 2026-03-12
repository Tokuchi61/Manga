package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStoreItemDefinitionAndInventoryFlow(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 14, 1, 0, 0, 0, time.UTC)
	userID := uuid.NewString()

	require.NoError(t, store.UpsertItemDefinition(context.Background(), entity.ItemDefinition{
		ItemID:     "avatar_frame_gold",
		ItemType:   "cosmetic",
		Stackable:  false,
		Equipable:  true,
		Consumable: false,
		MaxStack:   1,
		CreatedAt:  now,
		UpdatedAt:  now,
	}))

	definition, err := store.GetItemDefinition(context.Background(), "avatar_frame_gold")
	require.NoError(t, err)
	require.Equal(t, "cosmetic", definition.ItemType)

	require.NoError(t, store.UpsertInventoryEntry(context.Background(), entity.InventoryEntry{
		UserID:         userID,
		ItemID:         "avatar_frame_gold",
		Quantity:       1,
		Equipped:       true,
		LastSourceType: "mission",
		LastSourceRef:  "mission_1",
		RequestID:      "req-inv-1",
		CreatedAt:      now,
		UpdatedAt:      now,
	}))

	entry, err := store.GetInventoryEntry(context.Background(), userID, "avatar_frame_gold")
	require.NoError(t, err)
	require.Equal(t, 1, entry.Quantity)
	require.True(t, entry.Equipped)

	items, err := store.ListInventoryEntries(context.Background(), userID, "", false, "newest", 20, 0)
	require.NoError(t, err)
	require.Len(t, items, 1)

	require.NoError(t, store.DeleteInventoryEntry(context.Background(), userID, "avatar_frame_gold"))
	_, err = store.GetInventoryEntry(context.Background(), userID, "avatar_frame_gold")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestMemoryStoreGrantDedupAndRuntimeConfig(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 14, 1, 10, 0, 0, time.UTC)
	userID := uuid.NewString()

	entry := entity.InventoryEntry{
		UserID:         userID,
		ItemID:         "mana_potion",
		Quantity:       3,
		LastSourceType: "mission",
		LastSourceRef:  "daily_login",
		RequestID:      "req-inv-dedup-1",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, store.PutGrantDedup(context.Background(), "dedup-1", entry))

	persisted, err := store.GetGrantByDedup(context.Background(), "dedup-1")
	require.NoError(t, err)
	require.Equal(t, 3, persisted.Quantity)

	cfg, err := store.GetRuntimeConfig(context.Background())
	require.NoError(t, err)
	require.True(t, cfg.ClaimEnabled)

	cfg.ClaimEnabled = false
	cfg.UpdatedAt = now
	require.NoError(t, store.UpdateRuntimeConfig(context.Background(), cfg))

	updatedCfg, err := store.GetRuntimeConfig(context.Background())
	require.NoError(t, err)
	require.False(t, updatedCfg.ClaimEnabled)
}

func TestMemoryStoreSnapshotRoundtrip(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 14, 1, 20, 0, 0, time.UTC)
	userID := uuid.NewString()

	require.NoError(t, store.UpsertItemDefinition(context.Background(), entity.ItemDefinition{
		ItemID:     "ticket_epic",
		ItemType:   "token",
		Stackable:  true,
		Equipable:  false,
		Consumable: true,
		MaxStack:   999,
		CreatedAt:  now,
		UpdatedAt:  now,
	}))
	require.NoError(t, store.UpsertInventoryEntry(context.Background(), entity.InventoryEntry{
		UserID:         userID,
		ItemID:         "ticket_epic",
		Quantity:       10,
		LastSourceType: "admin",
		LastSourceRef:  "seed",
		CreatedAt:      now,
		UpdatedAt:      now,
	}))

	payload, err := store.Snapshot()
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	restored := NewMemoryStore()
	require.NoError(t, restored.RestoreSnapshot(payload))

	entry, err := restored.GetInventoryEntry(context.Background(), userID, "ticket_epic")
	require.NoError(t, err)
	require.Equal(t, 10, entry.Quantity)
}
