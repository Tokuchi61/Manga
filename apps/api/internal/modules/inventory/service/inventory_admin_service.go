package service

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/entity"
)

func (s *InventoryService) GetRuntimeConfig(ctx context.Context) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *InventoryService) UpdateReadState(ctx context.Context, request dto.UpdateReadStateRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.ReadEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()

	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *InventoryService) UpdateClaimState(ctx context.Context, request dto.UpdateClaimStateRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.ClaimEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()

	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *InventoryService) UpdateConsumeState(ctx context.Context, request dto.UpdateConsumeStateRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.ConsumeEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()

	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *InventoryService) UpdateEquipStateRuntime(ctx context.Context, request dto.UpdateEquipStateRuntimeRequest) (dto.RuntimeConfigResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.EquipEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()

	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return mapRuntimeConfig(cfg), nil
}

func (s *InventoryService) UpsertItemDefinition(ctx context.Context, request dto.UpsertItemDefinitionRequest) (dto.ItemDefinitionResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ItemDefinitionResponse{}, err
	}

	now := s.now().UTC()
	definition := entity.ItemDefinition{
		ItemID:     request.ItemID,
		ItemType:   normalizeItemType(request.ItemType),
		Stackable:  request.Stackable,
		Equipable:  request.Equipable,
		Consumable: request.Consumable,
		MaxStack:   request.MaxStack,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if !definition.Stackable {
		definition.MaxStack = 1
	} else if definition.MaxStack <= 0 {
		definition.MaxStack = defaultStackCap
	}

	if err := s.store.UpsertItemDefinition(ctx, definition); err != nil {
		return dto.ItemDefinitionResponse{}, err
	}
	return toDefinitionResponse(definition), nil
}

func (s *InventoryService) ListItemDefinitions(ctx context.Context, itemType string) (dto.ListItemDefinitionsResponse, error) {
	definitions, err := s.store.ListItemDefinitions(ctx, itemType)
	if err != nil {
		return dto.ListItemDefinitionsResponse{}, err
	}
	items := make([]dto.ItemDefinitionResponse, 0, len(definitions))
	for _, definition := range definitions {
		items = append(items, toDefinitionResponse(definition))
	}
	return dto.ListItemDefinitionsResponse{Items: items, Count: len(items)}, nil
}
