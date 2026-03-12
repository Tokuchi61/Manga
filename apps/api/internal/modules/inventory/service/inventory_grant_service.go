package service

import (
	"context"
	"strings"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/entity"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/repository"
)

func (s *InventoryService) applyGrant(ctx context.Context, userID string, itemID string, quantity int, sourceType string, sourceRef string, requestID string, correlationID string, expiresAt *time.Time) (entity.InventoryEntry, entity.ItemDefinition, bool, error) {
	definition, err := s.resolveDefinition(ctx, itemID)
	if err != nil {
		return entity.InventoryEntry{}, entity.ItemDefinition{}, false, err
	}

	if !definition.Stackable && quantity > 1 {
		return entity.InventoryEntry{}, entity.ItemDefinition{}, false, ErrConflict
	}

	dedupKey := buildDedupKey(userID, requestID, itemID, sourceType, sourceRef)
	if dedupKey != "" {
		dedupEntry, err := s.store.GetGrantByDedup(ctx, dedupKey)
		if err == nil {
			return dedupEntry, definition, false, nil
		}
		if err != repository.ErrNotFound {
			return entity.InventoryEntry{}, entity.ItemDefinition{}, false, err
		}
	}

	now := s.now().UTC()
	entry, err := s.store.GetInventoryEntry(ctx, userID, itemID)
	if err != nil && err != repository.ErrNotFound {
		return entity.InventoryEntry{}, entity.ItemDefinition{}, false, err
	}

	created := false
	if err == repository.ErrNotFound {
		entry = entity.InventoryEntry{
			UserID:    userID,
			ItemID:    itemID,
			Quantity:  quantity,
			Equipped:  false,
			CreatedAt: now,
			UpdatedAt: now,
		}
		created = true
	} else {
		if !definition.Stackable {
			return entity.InventoryEntry{}, entity.ItemDefinition{}, false, ErrItemAlreadyOwned
		}
		nextQuantity := entry.Quantity + quantity
		if nextQuantity > stackCap(definition) {
			return entity.InventoryEntry{}, entity.ItemDefinition{}, false, ErrConflict
		}
		entry.Quantity = nextQuantity
		entry.UpdatedAt = now
	}

	entry.LastSourceType = normalize(sourceType)
	entry.LastSourceRef = strings.TrimSpace(sourceRef)
	entry.RequestID = strings.TrimSpace(requestID)
	entry.CorrelationID = strings.TrimSpace(correlationID)
	if expiresAt != nil {
		t := expiresAt.UTC()
		entry.ExpiresAt = &t
	}

	if err := s.store.UpsertInventoryEntry(ctx, entry); err != nil {
		return entity.InventoryEntry{}, entity.ItemDefinition{}, false, err
	}
	if dedupKey != "" {
		if err := s.store.PutGrantDedup(ctx, dedupKey, entry); err != nil {
			return entity.InventoryEntry{}, entity.ItemDefinition{}, false, err
		}
	}

	return entry, definition, created, nil
}

func (s *InventoryService) AdminGrantItem(ctx context.Context, request dto.AdminGrantItemRequest) (dto.AdminGrantItemResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.AdminGrantItemResponse{}, err
	}
	if err := validateQuantity(request.Quantity); err != nil {
		return dto.AdminGrantItemResponse{}, err
	}

	entry, definition, created, err := s.applyGrant(ctx, request.TargetUserID, request.ItemID, request.Quantity, request.SourceType, request.SourceRef, request.RequestID, request.CorrelationID, request.ExpiresAt)
	if err != nil {
		return dto.AdminGrantItemResponse{}, err
	}

	return dto.AdminGrantItemResponse{Status: "granted", Created: created, Inventory: toEntryResponse(definition, entry)}, nil
}

func (s *InventoryService) AdminRevokeItem(ctx context.Context, request dto.AdminRevokeItemRequest) (dto.AdminRevokeItemResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.AdminRevokeItemResponse{}, err
	}
	if err := validateQuantity(request.Quantity); err != nil {
		return dto.AdminRevokeItemResponse{}, err
	}

	entry, err := s.store.GetInventoryEntry(ctx, request.TargetUserID, request.ItemID)
	if err != nil {
		if err == repository.ErrNotFound {
			return dto.AdminRevokeItemResponse{}, ErrNotFound
		}
		return dto.AdminRevokeItemResponse{}, err
	}
	if entry.Quantity < request.Quantity {
		return dto.AdminRevokeItemResponse{}, ErrInsufficientQuantity
	}

	now := s.now().UTC()
	entry.Quantity -= request.Quantity
	entry.UpdatedAt = now
	if entry.Quantity <= 0 {
		if err := s.store.DeleteInventoryEntry(ctx, request.TargetUserID, request.ItemID); err != nil {
			if err == repository.ErrNotFound {
				return dto.AdminRevokeItemResponse{}, ErrNotFound
			}
			return dto.AdminRevokeItemResponse{}, err
		}
		return dto.AdminRevokeItemResponse{Status: "revoked", ItemID: request.ItemID, RemainingQuantity: 0, RevokedAt: now}, nil
	}

	if err := s.store.UpsertInventoryEntry(ctx, entry); err != nil {
		return dto.AdminRevokeItemResponse{}, err
	}
	return dto.AdminRevokeItemResponse{Status: "revoked", ItemID: request.ItemID, RemainingQuantity: entry.Quantity, RevokedAt: now}, nil
}
