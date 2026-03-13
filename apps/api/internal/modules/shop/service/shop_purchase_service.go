package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/entity"
	shoprepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/repository"
	"github.com/google/uuid"
)

func (s *ShopService) CreatePurchaseIntent(ctx context.Context, request dto.CreatePurchaseIntentRequest) (dto.CreatePurchaseIntentResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CreatePurchaseIntentResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.CreatePurchaseIntentResponse{}, err
	}
	if err := s.requireCatalogEnabled(cfg.CatalogEnabled); err != nil {
		return dto.CreatePurchaseIntentResponse{}, err
	}
	if err := s.requirePurchaseEnabled(cfg.PurchaseEnabled); err != nil {
		return dto.CreatePurchaseIntentResponse{}, err
	}

	now := s.now().UTC()
	product, offer, err := s.resolveOfferForPurchase(ctx, now, cfg, request.ProductID, request.OfferID, true)
	if err != nil {
		return dto.CreatePurchaseIntentResponse{}, err
	}
	if normalizeValue(product.State) != entity.ProductStateActive {
		return dto.CreatePurchaseIntentResponse{}, ErrNotFound
	}

	if product.VIPRequired && !request.ActorVIP {
		return dto.CreatePurchaseIntentResponse{}, ErrEligibilityFailed
	}
	if request.ActorLevel < product.MinLevel {
		return dto.CreatePurchaseIntentResponse{}, ErrEligibilityFailed
	}

	if strings.TrimSpace(request.RequestID) != "" {
		dedupKey := buildPurchaseDedupKey(request.ActorUserID, offer.OfferID, request.RequestID)
		existingIntent, dedupErr := s.store.GetPurchaseDedup(ctx, dedupKey)
		if dedupErr == nil {
			return dto.CreatePurchaseIntentResponse{
				Status:         "idempotent",
				IntentID:       existingIntent.IntentID,
				ProductID:      existingIntent.ProductID,
				OfferID:        existingIntent.OfferID,
				FinalPriceMana: existingIntent.FinalPriceMana,
				Currency:       existingIntent.Currency,
				PurchaseStatus: existingIntent.Status,
			}, nil
		}
		if dedupErr != nil && !errors.Is(dedupErr, shoprepository.ErrNotFound) {
			return dto.CreatePurchaseIntentResponse{}, dedupErr
		}
	}

	if product.SinglePurchase {
		latest, latestErr := s.store.FindLatestPurchaseByUserProduct(ctx, request.ActorUserID, product.ProductID)
		if latestErr == nil {
			switch normalizeValue(latest.Status) {
			case entity.PurchaseStatusDeliveryPending, entity.PurchaseStatusCompleted, entity.PurchaseStatusRecoveryRequired:
				return dto.CreatePurchaseIntentResponse{}, ErrAlreadyOwned
			}
		}
		if latestErr != nil && !errors.Is(latestErr, shoprepository.ErrNotFound) {
			return dto.CreatePurchaseIntentResponse{}, latestErr
		}
	}

	intent := entity.PurchaseIntent{
		IntentID:       uuid.NewString(),
		UserID:         request.ActorUserID,
		ProductID:      product.ProductID,
		OfferID:        offer.OfferID,
		FinalPriceMana: calculateFinalPrice(offer.PriceMana, offer.DiscountPercent),
		Currency:       entity.CurrencyMana,
		Status:         entity.PurchaseStatusDeliveryPending,
		RequestID:      request.RequestID,
		CorrelationID:  request.CorrelationID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := s.store.CreatePurchaseIntent(ctx, intent); err != nil {
		return dto.CreatePurchaseIntentResponse{}, err
	}

	if strings.TrimSpace(request.RequestID) != "" {
		dedupKey := buildPurchaseDedupKey(request.ActorUserID, offer.OfferID, request.RequestID)
		if err := s.store.PutPurchaseDedup(ctx, dedupKey, intent); err != nil {
			return dto.CreatePurchaseIntentResponse{}, err
		}
	}

	return dto.CreatePurchaseIntentResponse{
		Status:         "intent_created",
		IntentID:       intent.IntentID,
		ProductID:      intent.ProductID,
		OfferID:        intent.OfferID,
		FinalPriceMana: intent.FinalPriceMana,
		Currency:       intent.Currency,
		PurchaseStatus: intent.Status,
	}, nil
}

func (s *ShopService) RequestPurchaseRecovery(ctx context.Context, request dto.RequestPurchaseRecoveryRequest) (dto.RequestPurchaseRecoveryResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RequestPurchaseRecoveryResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RequestPurchaseRecoveryResponse{}, err
	}
	if err := s.requirePurchaseEnabled(cfg.PurchaseEnabled); err != nil {
		return dto.RequestPurchaseRecoveryResponse{}, err
	}

	intent, err := s.store.GetPurchaseIntent(ctx, request.IntentID)
	if err != nil {
		if errors.Is(err, shoprepository.ErrNotFound) {
			return dto.RequestPurchaseRecoveryResponse{}, ErrNotFound
		}
		return dto.RequestPurchaseRecoveryResponse{}, err
	}
	if normalizeValue(intent.UserID) != normalizeValue(request.ActorUserID) {
		return dto.RequestPurchaseRecoveryResponse{}, ErrForbiddenAction
	}
	if normalizeValue(intent.Status) == entity.PurchaseStatusRecoveryRequired {
		return dto.RequestPurchaseRecoveryResponse{Status: "idempotent", IntentID: intent.IntentID, PurchaseStatus: intent.Status}, nil
	}
	if normalizeValue(intent.Status) == entity.PurchaseStatusCompleted {
		return dto.RequestPurchaseRecoveryResponse{}, ErrConflict
	}

	intent.Status = entity.PurchaseStatusRecoveryRequired
	intent.UpdatedAt = s.now().UTC()
	if err := s.store.UpdatePurchaseIntent(ctx, intent); err != nil {
		return dto.RequestPurchaseRecoveryResponse{}, err
	}

	return dto.RequestPurchaseRecoveryResponse{Status: "recovery_requested", IntentID: intent.IntentID, PurchaseStatus: intent.Status}, nil
}
