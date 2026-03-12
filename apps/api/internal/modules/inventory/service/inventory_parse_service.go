package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/entity"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/repository"
)

const defaultStackCap = 999999

func normalize(raw string) string {
	return strings.TrimSpace(strings.ToLower(raw))
}

func normalizeItemType(raw string) string {
	return strings.TrimSpace(strings.ToLower(raw))
}

func normalizeSort(raw string) string {
	value := normalize(raw)
	if value == "oldest" {
		return "oldest"
	}
	return "newest"
}

func stackCap(definition entity.ItemDefinition) int {
	if !definition.Stackable {
		return 1
	}
	if definition.MaxStack <= 0 {
		return defaultStackCap
	}
	return definition.MaxStack
}

func buildDedupKey(userID string, requestID string, itemID string, sourceType string, sourceRef string) string {
	if strings.TrimSpace(requestID) == "" {
		return ""
	}
	return normalize(userID) + ":" + normalize(requestID) + ":" + normalize(itemID) + ":" + normalize(sourceType) + ":" + normalize(sourceRef)
}

func (s *InventoryService) resolveDefinition(ctx context.Context, itemID string) (entity.ItemDefinition, error) {
	definition, err := s.store.GetItemDefinition(ctx, itemID)
	if err == nil {
		return definition, nil
	}
	if err == repository.ErrNotFound {
		return entity.ItemDefinition{}, ErrNotFound
	}
	return entity.ItemDefinition{}, err
}

func toEntryResponse(definition entity.ItemDefinition, entry entity.InventoryEntry) dto.InventoryEntryItemResponse {
	return dto.InventoryEntryItemResponse{
		ItemID:         entry.ItemID,
		ItemType:       definition.ItemType,
		Quantity:       entry.Quantity,
		Equipped:       entry.Equipped,
		Stackable:      definition.Stackable,
		Equipable:      definition.Equipable,
		Consumable:     definition.Consumable,
		MaxStack:       stackCap(definition),
		LastSourceType: entry.LastSourceType,
		LastSourceRef:  entry.LastSourceRef,
		ExpiresAt:      entry.ExpiresAt,
		UpdatedAt:      entry.UpdatedAt,
	}
}

func toDetailResponse(definition entity.ItemDefinition, entry entity.InventoryEntry) dto.InventoryItemDetailResponse {
	return dto.InventoryItemDetailResponse{
		ItemID:         entry.ItemID,
		ItemType:       definition.ItemType,
		Quantity:       entry.Quantity,
		Equipped:       entry.Equipped,
		Stackable:      definition.Stackable,
		Equipable:      definition.Equipable,
		Consumable:     definition.Consumable,
		MaxStack:       stackCap(definition),
		LastSourceType: entry.LastSourceType,
		LastSourceRef:  entry.LastSourceRef,
		ExpiresAt:      entry.ExpiresAt,
		CreatedAt:      entry.CreatedAt,
		UpdatedAt:      entry.UpdatedAt,
	}
}

func toDefinitionResponse(definition entity.ItemDefinition) dto.ItemDefinitionResponse {
	return dto.ItemDefinitionResponse{
		ItemID:     definition.ItemID,
		ItemType:   definition.ItemType,
		Stackable:  definition.Stackable,
		Equipable:  definition.Equipable,
		Consumable: definition.Consumable,
		MaxStack:   stackCap(definition),
		UpdatedAt:  definition.UpdatedAt,
	}
}

func mapRuntimeConfig(cfg entity.RuntimeConfig) dto.RuntimeConfigResponse {
	return dto.RuntimeConfigResponse{
		ReadEnabled:    cfg.ReadEnabled,
		ClaimEnabled:   cfg.ClaimEnabled,
		ConsumeEnabled: cfg.ConsumeEnabled,
		EquipEnabled:   cfg.EquipEnabled,
		UpdatedAt:      cfg.UpdatedAt,
	}
}

func validateQuantity(quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("%w: invalid quantity", ErrValidation)
	}
	return nil
}
