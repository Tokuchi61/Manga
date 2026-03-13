package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/entity"
	shoprepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/repository"
)

func (s *ShopService) GetRuntimeConfig(ctx context.Context) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return dto.RuntimeConfigResponse{
		CatalogEnabled:  cfg.CatalogEnabled,
		PurchaseEnabled: cfg.PurchaseEnabled,
		CampaignEnabled: cfg.CampaignEnabled,
		UpdatedAt:       cfg.UpdatedAt,
	}, nil
}

func (s *ShopService) UpdateCatalogState(ctx context.Context, request dto.UpdateCatalogStateRequest) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.CatalogEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *ShopService) UpdatePurchaseState(ctx context.Context, request dto.UpdatePurchaseStateRequest) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.PurchaseEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *ShopService) UpdateCampaignState(ctx context.Context, request dto.UpdateCampaignStateRequest) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.CampaignEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *ShopService) UpsertProductDefinition(ctx context.Context, request dto.UpsertProductDefinitionRequest) (dto.ProductDefinitionResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ProductDefinitionResponse{}, err
	}

	now := s.now().UTC()
	product := entity.ProductDefinition{
		ProductID:       request.ProductID,
		Name:            request.Name,
		Category:        request.Category,
		State:           request.State,
		InventoryItemID: request.InventoryItemID,
		SlotID:          request.SlotID,
		SinglePurchase:  request.SinglePurchase,
		VIPRequired:     request.VIPRequired,
		MinLevel:        request.MinLevel,
		UpdatedAt:       now,
	}

	existing, err := s.store.GetProductDefinition(ctx, request.ProductID)
	if err == nil {
		product.CreatedAt = existing.CreatedAt
	} else if err != nil && !errors.Is(err, shoprepository.ErrNotFound) {
		return dto.ProductDefinitionResponse{}, err
	}
	if product.CreatedAt.IsZero() {
		product.CreatedAt = now
	}

	if err := s.store.UpsertProductDefinition(ctx, product); err != nil {
		return dto.ProductDefinitionResponse{}, err
	}
	return toProductDefinitionResponse(product), nil
}

func (s *ShopService) ListProductDefinitions(ctx context.Context, request dto.ListProductDefinitionsRequest) (dto.ListProductDefinitionsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListProductDefinitionsResponse{}, err
	}

	products, err := s.store.ListProductDefinitions(ctx, request.State, request.Limit, request.Offset)
	if err != nil {
		return dto.ListProductDefinitionsResponse{}, err
	}

	items := make([]dto.ProductDefinitionResponse, 0, len(products))
	for _, product := range products {
		items = append(items, toProductDefinitionResponse(product))
	}
	return dto.ListProductDefinitionsResponse{Items: items, Count: len(items)}, nil
}

func (s *ShopService) UpsertOfferDefinition(ctx context.Context, request dto.UpsertOfferDefinitionRequest) (dto.OfferDefinitionResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OfferDefinitionResponse{}, err
	}
	if request.StartsAt != nil && request.EndsAt != nil && request.StartsAt.After(*request.EndsAt) {
		return dto.OfferDefinitionResponse{}, fmt.Errorf("%w: starts_at_after_ends_at", ErrValidation)
	}

	if _, err := s.store.GetProductDefinition(ctx, request.ProductID); err != nil {
		if errors.Is(err, shoprepository.ErrNotFound) {
			return dto.OfferDefinitionResponse{}, ErrNotFound
		}
		return dto.OfferDefinitionResponse{}, err
	}

	now := s.now().UTC()
	offer := entity.OfferDefinition{
		OfferID:         request.OfferID,
		ProductID:       request.ProductID,
		Title:           request.Title,
		Visibility:      request.Visibility,
		PriceMana:       request.PriceMana,
		DiscountPercent: request.DiscountPercent,
		Active:          request.Active,
		StartsAt:        request.StartsAt,
		EndsAt:          request.EndsAt,
		UpdatedAt:       now,
	}

	existing, err := s.store.GetOfferDefinition(ctx, request.OfferID)
	if err == nil {
		offer.CreatedAt = existing.CreatedAt
	} else if err != nil && !errors.Is(err, shoprepository.ErrNotFound) {
		return dto.OfferDefinitionResponse{}, err
	}
	if offer.CreatedAt.IsZero() {
		offer.CreatedAt = now
	}

	if err := s.store.UpsertOfferDefinition(ctx, offer); err != nil {
		return dto.OfferDefinitionResponse{}, err
	}
	return toOfferDefinitionResponse(offer), nil
}

func (s *ShopService) ListOfferDefinitions(ctx context.Context, request dto.ListOfferDefinitionsRequest) (dto.ListOfferDefinitionsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListOfferDefinitionsResponse{}, err
	}

	offers, err := s.store.ListOfferDefinitions(ctx, request.ProductID, request.Visibility, request.ActiveOnly, request.Limit, request.Offset)
	if err != nil {
		return dto.ListOfferDefinitionsResponse{}, err
	}

	items := make([]dto.OfferDefinitionResponse, 0, len(offers))
	for _, offer := range offers {
		items = append(items, toOfferDefinitionResponse(offer))
	}
	return dto.ListOfferDefinitionsResponse{Items: items, Count: len(items)}, nil
}
