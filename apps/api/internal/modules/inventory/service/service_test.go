package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/dto"
	inventoryrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestInventoryServiceClaimConsumeEquipAndListFlow(t *testing.T) {
	store := inventoryrepository.NewMemoryStore()
	svc := New(store, validation.New())
	now := time.Date(2026, 3, 14, 2, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }

	actorID := uuid.NewString()

	_, err := svc.UpsertItemDefinition(context.Background(), dto.UpsertItemDefinitionRequest{
		ItemID:     "mana_potion",
		ItemType:   "consumable",
		Stackable:  true,
		Consumable: true,
		MaxStack:   99,
	})
	require.NoError(t, err)

	claimRes, err := svc.ClaimInventoryItem(context.Background(), dto.ClaimInventoryItemRequest{
		ActorUserID: actorID,
		ItemID:      "mana_potion",
		Quantity:    5,
		SourceType:  "mission",
		SourceRef:   "daily_1",
		RequestID:   "req-inventory-claim-1",
	})
	require.NoError(t, err)
	require.True(t, claimRes.Created)
	require.Equal(t, 5, claimRes.Inventory.Quantity)

	consumeRes, err := svc.ConsumeInventoryItem(context.Background(), dto.ConsumeInventoryItemRequest{
		ActorUserID: actorID,
		ItemID:      "mana_potion",
		Quantity:    2,
	})
	require.NoError(t, err)
	require.Equal(t, 3, consumeRes.RemainingQuantity)

	listRes, err := svc.ListInventoryEntries(context.Background(), dto.ListInventoryEntriesRequest{ActorUserID: actorID})
	require.NoError(t, err)
	require.Equal(t, 1, listRes.Count)
	require.Equal(t, 3, listRes.Items[0].Quantity)
}

func TestInventoryServiceEquipAndRevokeFlow(t *testing.T) {
	store := inventoryrepository.NewMemoryStore()
	svc := New(store, validation.New())
	actorID := uuid.NewString()
	adminTargetID := uuid.NewString()

	_, err := svc.UpsertItemDefinition(context.Background(), dto.UpsertItemDefinitionRequest{
		ItemID:    "avatar_frame_gold",
		ItemType:  "cosmetic",
		Stackable: false,
		Equipable: true,
	})
	require.NoError(t, err)

	grantRes, err := svc.AdminGrantItem(context.Background(), dto.AdminGrantItemRequest{
		TargetUserID: adminTargetID,
		ItemID:       "avatar_frame_gold",
		Quantity:     1,
		SourceType:   "admin",
		SourceRef:    "manual_grant",
		RequestID:    "req-admin-grant-1",
	})
	require.NoError(t, err)
	require.True(t, grantRes.Created)

	_, err = svc.UpdateEquipState(context.Background(), dto.UpdateEquipStateRequest{
		ActorUserID: adminTargetID,
		ItemID:      "avatar_frame_gold",
		Enabled:     true,
	})
	require.NoError(t, err)

	revokeRes, err := svc.AdminRevokeItem(context.Background(), dto.AdminRevokeItemRequest{
		TargetUserID: adminTargetID,
		ItemID:       "avatar_frame_gold",
		Quantity:     1,
	})
	require.NoError(t, err)
	require.Equal(t, 0, revokeRes.RemainingQuantity)

	_, err = svc.GetInventoryItemDetail(context.Background(), dto.GetInventoryItemDetailRequest{ActorUserID: adminTargetID, ItemID: "avatar_frame_gold"})
	require.True(t, errors.Is(err, ErrNotFound))

	_, err = svc.UpdateEquipState(context.Background(), dto.UpdateEquipStateRequest{ActorUserID: actorID, ItemID: "avatar_frame_gold", Enabled: true})
	require.True(t, errors.Is(err, ErrNotFound))
}

func TestInventoryServiceRuntimeTogglesAffectOperations(t *testing.T) {
	store := inventoryrepository.NewMemoryStore()
	svc := New(store, validation.New())
	actorID := uuid.NewString()

	_, err := svc.UpsertItemDefinition(context.Background(), dto.UpsertItemDefinitionRequest{
		ItemID:     "ticket_epic",
		ItemType:   "token",
		Stackable:  true,
		Consumable: true,
	})
	require.NoError(t, err)

	_, err = svc.UpdateClaimState(context.Background(), dto.UpdateClaimStateRequest{Enabled: false})
	require.NoError(t, err)
	_, err = svc.ClaimInventoryItem(context.Background(), dto.ClaimInventoryItemRequest{ActorUserID: actorID, ItemID: "ticket_epic", Quantity: 1, SourceType: "mission", RequestID: "req-claim-disabled"})
	require.True(t, errors.Is(err, ErrClaimDisabled))

	_, err = svc.UpdateClaimState(context.Background(), dto.UpdateClaimStateRequest{Enabled: true})
	require.NoError(t, err)
	_, err = svc.ClaimInventoryItem(context.Background(), dto.ClaimInventoryItemRequest{ActorUserID: actorID, ItemID: "ticket_epic", Quantity: 1, SourceType: "mission", RequestID: "req-claim-enable"})
	require.NoError(t, err)

	_, err = svc.UpdateConsumeState(context.Background(), dto.UpdateConsumeStateRequest{Enabled: false})
	require.NoError(t, err)
	_, err = svc.ConsumeInventoryItem(context.Background(), dto.ConsumeInventoryItemRequest{ActorUserID: actorID, ItemID: "ticket_epic", Quantity: 1})
	require.True(t, errors.Is(err, ErrConsumeDisabled))

	_, err = svc.UpdateReadState(context.Background(), dto.UpdateReadStateRequest{Enabled: false})
	require.NoError(t, err)
	_, err = svc.ListInventoryEntries(context.Background(), dto.ListInventoryEntriesRequest{ActorUserID: actorID})
	require.True(t, errors.Is(err, ErrReadDisabled))
}
