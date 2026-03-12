package service

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/dto"
)

func (s *InventoryService) ClaimInventoryItem(ctx context.Context, request dto.ClaimInventoryItemRequest) (dto.ClaimInventoryItemResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ClaimInventoryItemResponse{}, err
	}
	if err := validateQuantity(request.Quantity); err != nil {
		return dto.ClaimInventoryItemResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ClaimInventoryItemResponse{}, err
	}
	if err := s.requireClaimEnabled(cfg.ClaimEnabled); err != nil {
		return dto.ClaimInventoryItemResponse{}, err
	}

	entry, definition, created, err := s.applyGrant(ctx, request.ActorUserID, request.ItemID, request.Quantity, request.SourceType, request.SourceRef, request.RequestID, request.CorrelationID, request.ExpiresAt)
	if err != nil {
		return dto.ClaimInventoryItemResponse{}, err
	}

	return dto.ClaimInventoryItemResponse{Status: "claimed", Created: created, Inventory: toEntryResponse(definition, entry)}, nil
}
