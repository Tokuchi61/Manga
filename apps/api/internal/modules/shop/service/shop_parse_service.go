package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/entity"
	shoprepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/repository"
)

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func calculateFinalPrice(price int, discountPercent int) int {
	if discountPercent <= 0 {
		return price
	}
	if discountPercent > 100 {
		discountPercent = 100
	}
	finalPrice := price - (price*discountPercent)/100
	if finalPrice <= 0 {
		return 1
	}
	return finalPrice
}

func isOfferWithinWindow(offer entity.OfferDefinition, now time.Time) bool {
	utcNow := now.UTC()
	if offer.StartsAt != nil && utcNow.Before(offer.StartsAt.UTC()) {
		return false
	}
	if offer.EndsAt != nil && utcNow.After(offer.EndsAt.UTC()) {
		return false
	}
	return true
}

func isOfferVisibleForActor(offer entity.OfferDefinition, now time.Time, includeCampaign bool, campaignEnabled bool) bool {
	if !offer.Active {
		return false
	}
	if !isOfferWithinWindow(offer, now) {
		return false
	}
	switch normalizeValue(offer.Visibility) {
	case entity.OfferVisibilityHidden:
		return false
	case entity.OfferVisibilityCampaignOnly:
		if !campaignEnabled || !includeCampaign {
			return false
		}
	}
	return true
}

func buildPurchaseDedupKey(userID string, offerID string, requestID string) string {
	return normalizeValue(userID) + ":" + normalizeValue(offerID) + ":" + normalizeValue(requestID)
}

func toProductDefinitionResponse(product entity.ProductDefinition) dto.ProductDefinitionResponse {
	return dto.ProductDefinitionResponse{
		ProductID:       product.ProductID,
		Name:            product.Name,
		Category:        product.Category,
		State:           product.State,
		InventoryItemID: product.InventoryItemID,
		SlotID:          product.SlotID,
		SinglePurchase:  product.SinglePurchase,
		VIPRequired:     product.VIPRequired,
		MinLevel:        product.MinLevel,
		UpdatedAt:       product.UpdatedAt,
	}
}

func toOfferDefinitionResponse(offer entity.OfferDefinition) dto.OfferDefinitionResponse {
	return dto.OfferDefinitionResponse{
		OfferID:         offer.OfferID,
		ProductID:       offer.ProductID,
		Title:           offer.Title,
		Visibility:      offer.Visibility,
		PriceMana:       offer.PriceMana,
		DiscountPercent: offer.DiscountPercent,
		FinalPriceMana:  calculateFinalPrice(offer.PriceMana, offer.DiscountPercent),
		Active:          offer.Active,
		StartsAt:        offer.StartsAt,
		EndsAt:          offer.EndsAt,
		UpdatedAt:       offer.UpdatedAt,
	}
}

func toOfferPreviewResponse(offer entity.OfferDefinition) dto.OfferPreviewResponse {
	return dto.OfferPreviewResponse{
		OfferID:         offer.OfferID,
		Title:           offer.Title,
		Visibility:      offer.Visibility,
		PriceMana:       offer.PriceMana,
		DiscountPercent: offer.DiscountPercent,
		FinalPriceMana:  calculateFinalPrice(offer.PriceMana, offer.DiscountPercent),
		StartsAt:        offer.StartsAt,
		EndsAt:          offer.EndsAt,
	}
}

func toCatalogItemResponse(product entity.ProductDefinition, offer entity.OfferDefinition) dto.CatalogItemResponse {
	return dto.CatalogItemResponse{
		ProductID:       product.ProductID,
		Name:            product.Name,
		Category:        product.Category,
		InventoryItemID: product.InventoryItemID,
		OfferID:         offer.OfferID,
		OfferTitle:      offer.Title,
		Visibility:      offer.Visibility,
		PriceMana:       offer.PriceMana,
		DiscountPercent: offer.DiscountPercent,
		FinalPriceMana:  calculateFinalPrice(offer.PriceMana, offer.DiscountPercent),
	}
}

func chooseVisibleOffer(offers []entity.OfferDefinition, now time.Time, includeCampaign bool, campaignEnabled bool) (entity.OfferDefinition, error) {
	candidates := make([]entity.OfferDefinition, 0, len(offers))
	for _, offer := range offers {
		if !isOfferVisibleForActor(offer, now, includeCampaign, campaignEnabled) {
			continue
		}
		candidates = append(candidates, offer)
	}
	if len(candidates) == 0 {
		return entity.OfferDefinition{}, ErrNotFound
	}

	sort.Slice(candidates, func(i int, j int) bool {
		leftPrice := calculateFinalPrice(candidates[i].PriceMana, candidates[i].DiscountPercent)
		rightPrice := calculateFinalPrice(candidates[j].PriceMana, candidates[j].DiscountPercent)
		if leftPrice == rightPrice {
			return candidates[i].UpdatedAt.After(candidates[j].UpdatedAt)
		}
		return leftPrice < rightPrice
	})

	return candidates[0], nil
}

func applyOffsetLimit[T any](items []T, offset int, limit int) []T {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(items) {
		return []T{}
	}
	if limit <= 0 {
		limit = 50
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return append([]T(nil), items[offset:end]...)
}

func (s *ShopService) resolveOfferForPurchase(ctx context.Context, now time.Time, cfg entity.RuntimeConfig, productID string, offerID string, includeCampaign bool) (entity.ProductDefinition, entity.OfferDefinition, error) {
	if strings.TrimSpace(offerID) != "" {
		offer, err := s.store.GetOfferDefinition(ctx, offerID)
		if err != nil {
			if errors.Is(err, shoprepository.ErrNotFound) {
				return entity.ProductDefinition{}, entity.OfferDefinition{}, ErrNotFound
			}
			return entity.ProductDefinition{}, entity.OfferDefinition{}, err
		}
		product, err := s.store.GetProductDefinition(ctx, offer.ProductID)
		if err != nil {
			if errors.Is(err, shoprepository.ErrNotFound) {
				return entity.ProductDefinition{}, entity.OfferDefinition{}, ErrNotFound
			}
			return entity.ProductDefinition{}, entity.OfferDefinition{}, err
		}
		if strings.TrimSpace(productID) != "" && normalizeValue(productID) != normalizeValue(product.ProductID) {
			return entity.ProductDefinition{}, entity.OfferDefinition{}, fmt.Errorf("%w: product_offer_mismatch", ErrValidation)
		}
		if normalizeValue(offer.Visibility) == entity.OfferVisibilityCampaignOnly && !cfg.CampaignEnabled {
			return entity.ProductDefinition{}, entity.OfferDefinition{}, ErrCampaignDisabled
		}
		if !isOfferVisibleForActor(offer, now, includeCampaign, cfg.CampaignEnabled) {
			return entity.ProductDefinition{}, entity.OfferDefinition{}, ErrNotFound
		}
		return product, offer, nil
	}

	if strings.TrimSpace(productID) == "" {
		return entity.ProductDefinition{}, entity.OfferDefinition{}, fmt.Errorf("%w: offer_or_product_required", ErrValidation)
	}

	product, err := s.store.GetProductDefinition(ctx, productID)
	if err != nil {
		if errors.Is(err, shoprepository.ErrNotFound) {
			return entity.ProductDefinition{}, entity.OfferDefinition{}, ErrNotFound
		}
		return entity.ProductDefinition{}, entity.OfferDefinition{}, err
	}
	offers, err := s.store.ListOfferDefinitions(ctx, product.ProductID, "", true, 0, 0)
	if err != nil {
		return entity.ProductDefinition{}, entity.OfferDefinition{}, err
	}
	offer, err := chooseVisibleOffer(offers, now, includeCampaign, cfg.CampaignEnabled)
	if err != nil {
		return entity.ProductDefinition{}, entity.OfferDefinition{}, err
	}
	return product, offer, nil
}
