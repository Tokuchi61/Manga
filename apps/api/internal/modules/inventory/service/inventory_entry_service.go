package service

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/repository"
)

func (s *InventoryService) ListInventoryEntries(ctx context.Context, request dto.ListInventoryEntriesRequest) (dto.ListInventoryEntriesResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListInventoryEntriesResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ListInventoryEntriesResponse{}, err
	}
	if err := s.requireReadEnabled(cfg.ReadEnabled); err != nil {
		return dto.ListInventoryEntriesResponse{}, err
	}

	entries, err := s.store.ListInventoryEntries(ctx, request.ActorUserID, request.ItemType, request.EquippedOnly, normalizeSort(request.SortBy), request.Limit, request.Offset)
	if err != nil {
		return dto.ListInventoryEntriesResponse{}, err
	}

	items := make([]dto.InventoryEntryItemResponse, 0, len(entries))
	for _, entry := range entries {
		definition, err := s.resolveDefinition(ctx, entry.ItemID)
		if err != nil {
			if err == ErrNotFound {
				continue
			}
			return dto.ListInventoryEntriesResponse{}, err
		}
		items = append(items, toEntryResponse(definition, entry))
	}

	return dto.ListInventoryEntriesResponse{Items: items, Count: len(items)}, nil
}

func (s *InventoryService) GetInventoryItemDetail(ctx context.Context, request dto.GetInventoryItemDetailRequest) (dto.InventoryItemDetailResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.InventoryItemDetailResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.InventoryItemDetailResponse{}, err
	}
	if err := s.requireReadEnabled(cfg.ReadEnabled); err != nil {
		return dto.InventoryItemDetailResponse{}, err
	}

	entry, err := s.store.GetInventoryEntry(ctx, request.ActorUserID, request.ItemID)
	if err != nil {
		if err == repository.ErrNotFound {
			return dto.InventoryItemDetailResponse{}, ErrNotFound
		}
		return dto.InventoryItemDetailResponse{}, err
	}
	definition, err := s.resolveDefinition(ctx, entry.ItemID)
	if err != nil {
		return dto.InventoryItemDetailResponse{}, err
	}
	return toDetailResponse(definition, entry), nil
}

func (s *InventoryService) ConsumeInventoryItem(ctx context.Context, request dto.ConsumeInventoryItemRequest) (dto.ConsumeInventoryItemResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ConsumeInventoryItemResponse{}, err
	}
	if err := validateQuantity(request.Quantity); err != nil {
		return dto.ConsumeInventoryItemResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ConsumeInventoryItemResponse{}, err
	}
	if err := s.requireConsumeEnabled(cfg.ConsumeEnabled); err != nil {
		return dto.ConsumeInventoryItemResponse{}, err
	}

	definition, err := s.resolveDefinition(ctx, request.ItemID)
	if err != nil {
		return dto.ConsumeInventoryItemResponse{}, err
	}
	if !definition.Consumable {
		return dto.ConsumeInventoryItemResponse{}, ErrForbiddenAction
	}

	entry, err := s.store.GetInventoryEntry(ctx, request.ActorUserID, request.ItemID)
	if err != nil {
		if err == repository.ErrNotFound {
			return dto.ConsumeInventoryItemResponse{}, ErrNotFound
		}
		return dto.ConsumeInventoryItemResponse{}, err
	}
	if entry.Quantity < request.Quantity {
		return dto.ConsumeInventoryItemResponse{}, ErrInsufficientQuantity
	}

	now := s.now().UTC()
	entry.Quantity -= request.Quantity
	entry.RequestID = request.RequestID
	entry.CorrelationID = request.CorrelationID
	entry.UpdatedAt = now
	if entry.Quantity <= 0 {
		if err := s.store.DeleteInventoryEntry(ctx, request.ActorUserID, request.ItemID); err != nil {
			if err == repository.ErrNotFound {
				return dto.ConsumeInventoryItemResponse{}, ErrNotFound
			}
			return dto.ConsumeInventoryItemResponse{}, err
		}
		return dto.ConsumeInventoryItemResponse{Status: "consumed", ItemID: request.ItemID, RemainingQuantity: 0, ConsumedAt: now}, nil
	}

	if err := s.store.UpsertInventoryEntry(ctx, entry); err != nil {
		return dto.ConsumeInventoryItemResponse{}, err
	}

	return dto.ConsumeInventoryItemResponse{Status: "consumed", ItemID: request.ItemID, RemainingQuantity: entry.Quantity, ConsumedAt: now}, nil
}

func (s *InventoryService) UpdateEquipState(ctx context.Context, request dto.UpdateEquipStateRequest) (dto.UpdateEquipStateResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.UpdateEquipStateResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.UpdateEquipStateResponse{}, err
	}
	if err := s.requireEquipEnabled(cfg.EquipEnabled); err != nil {
		return dto.UpdateEquipStateResponse{}, err
	}

	definition, err := s.resolveDefinition(ctx, request.ItemID)
	if err != nil {
		return dto.UpdateEquipStateResponse{}, err
	}
	if !definition.Equipable {
		return dto.UpdateEquipStateResponse{}, ErrForbiddenAction
	}

	entry, err := s.store.GetInventoryEntry(ctx, request.ActorUserID, request.ItemID)
	if err != nil {
		if err == repository.ErrNotFound {
			return dto.UpdateEquipStateResponse{}, ErrNotFound
		}
		return dto.UpdateEquipStateResponse{}, err
	}

	now := s.now().UTC()
	entry.Equipped = request.Enabled
	entry.UpdatedAt = now
	if err := s.store.UpsertInventoryEntry(ctx, entry); err != nil {
		return dto.UpdateEquipStateResponse{}, err
	}

	return dto.UpdateEquipStateResponse{Status: "equip_state_updated", ItemID: request.ItemID, Equipped: request.Enabled, UpdatedAt: now}, nil
}
