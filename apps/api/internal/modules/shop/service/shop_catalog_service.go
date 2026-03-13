package service

import (
	"context"
	"errors"
	"sort"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/entity"
	shoprepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/repository"
)

func (s *ShopService) ListCatalog(ctx context.Context, request dto.ListCatalogRequest) (dto.ListCatalogResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListCatalogResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ListCatalogResponse{}, err
	}
	if err := s.requireCatalogEnabled(cfg.CatalogEnabled); err != nil {
		return dto.ListCatalogResponse{}, err
	}

	products, err := s.store.ListProductDefinitions(ctx, entity.ProductStateActive, 0, 0)
	if err != nil {
		return dto.ListCatalogResponse{}, err
	}

	now := s.now().UTC()
	items := make([]dto.CatalogItemResponse, 0, len(products))
	for _, product := range products {
		offers, offerErr := s.store.ListOfferDefinitions(ctx, product.ProductID, "", true, 0, 0)
		if offerErr != nil {
			return dto.ListCatalogResponse{}, offerErr
		}
		offer, chooseErr := chooseVisibleOffer(offers, now, request.IncludeCampaign, cfg.CampaignEnabled)
		if chooseErr != nil {
			if errors.Is(chooseErr, ErrNotFound) {
				continue
			}
			return dto.ListCatalogResponse{}, chooseErr
		}
		items = append(items, toCatalogItemResponse(product, offer))
	}

	sort.Slice(items, func(i int, j int) bool {
		if items[i].Category == items[j].Category {
			return items[i].ProductID < items[j].ProductID
		}
		return items[i].Category < items[j].Category
	})
	items = applyOffsetLimit(items, request.Offset, request.Limit)

	return dto.ListCatalogResponse{Items: items, Count: len(items)}, nil
}

func (s *ShopService) GetCatalogItem(ctx context.Context, request dto.GetCatalogItemRequest) (dto.CatalogItemDetailResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CatalogItemDetailResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.CatalogItemDetailResponse{}, err
	}
	if err := s.requireCatalogEnabled(cfg.CatalogEnabled); err != nil {
		return dto.CatalogItemDetailResponse{}, err
	}

	product, err := s.store.GetProductDefinition(ctx, request.ProductID)
	if err != nil {
		if errors.Is(err, shoprepository.ErrNotFound) {
			return dto.CatalogItemDetailResponse{}, ErrNotFound
		}
		return dto.CatalogItemDetailResponse{}, err
	}
	if normalizeValue(product.State) != entity.ProductStateActive {
		return dto.CatalogItemDetailResponse{}, ErrNotFound
	}

	offers, err := s.store.ListOfferDefinitions(ctx, product.ProductID, "", true, 0, 0)
	if err != nil {
		return dto.CatalogItemDetailResponse{}, err
	}

	now := s.now().UTC()
	visibleOffers := make([]dto.OfferPreviewResponse, 0, len(offers))
	for _, offer := range offers {
		if !isOfferVisibleForActor(offer, now, request.IncludeCampaign, cfg.CampaignEnabled) {
			continue
		}
		visibleOffers = append(visibleOffers, toOfferPreviewResponse(offer))
	}
	if len(visibleOffers) == 0 {
		return dto.CatalogItemDetailResponse{}, ErrNotFound
	}

	sort.Slice(visibleOffers, func(i int, j int) bool {
		if visibleOffers[i].FinalPriceMana == visibleOffers[j].FinalPriceMana {
			return visibleOffers[i].OfferID < visibleOffers[j].OfferID
		}
		return visibleOffers[i].FinalPriceMana < visibleOffers[j].FinalPriceMana
	})

	return dto.CatalogItemDetailResponse{
		ProductID:       product.ProductID,
		Name:            product.Name,
		Category:        product.Category,
		State:           product.State,
		InventoryItemID: product.InventoryItemID,
		SlotID:          product.SlotID,
		SinglePurchase:  product.SinglePurchase,
		Offers:          visibleOffers,
		Count:           len(visibleOffers),
	}, nil
}
