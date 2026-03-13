package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/dto"
	shoprepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestShopServiceCatalogPurchaseAndRecoveryFlow(t *testing.T) {
	store := shoprepository.NewMemoryStore()
	svc := New(store, validation.New())
	now := time.Date(2026, 3, 19, 12, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }

	actorID := uuid.NewString()
	setupActiveProductAndOffer(t, svc, now)

	catalogRes, err := svc.ListCatalog(context.Background(), dto.ListCatalogRequest{ActorUserID: actorID})
	require.NoError(t, err)
	require.Equal(t, 1, catalogRes.Count)

	purchaseRes, err := svc.CreatePurchaseIntent(context.Background(), dto.CreatePurchaseIntentRequest{
		ActorUserID: actorID,
		ProductID:   "product_avatar_01",
		RequestID:   "req-shop-purchase-1",
	})
	require.NoError(t, err)
	require.Equal(t, "intent_created", purchaseRes.Status)

	idempotentRes, err := svc.CreatePurchaseIntent(context.Background(), dto.CreatePurchaseIntentRequest{
		ActorUserID: actorID,
		ProductID:   "product_avatar_01",
		RequestID:   "req-shop-purchase-1",
	})
	require.NoError(t, err)
	require.Equal(t, "idempotent", idempotentRes.Status)

	_, err = svc.CreatePurchaseIntent(context.Background(), dto.CreatePurchaseIntentRequest{
		ActorUserID: actorID,
		ProductID:   "product_avatar_01",
		RequestID:   "req-shop-purchase-2",
	})
	require.True(t, errors.Is(err, ErrAlreadyOwned))

	recoveryRes, err := svc.RequestPurchaseRecovery(context.Background(), dto.RequestPurchaseRecoveryRequest{
		ActorUserID: actorID,
		IntentID:    purchaseRes.IntentID,
	})
	require.NoError(t, err)
	require.Equal(t, "recovery_requested", recoveryRes.Status)
}

func TestShopServiceRuntimeTogglesAffectFlows(t *testing.T) {
	store := shoprepository.NewMemoryStore()
	svc := New(store, validation.New())
	now := time.Date(2026, 3, 19, 12, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }
	actorID := uuid.NewString()
	setupActiveProductAndOffer(t, svc, now)

	_, err := svc.UpdateCatalogState(context.Background(), dto.UpdateCatalogStateRequest{Enabled: false})
	require.NoError(t, err)
	_, err = svc.ListCatalog(context.Background(), dto.ListCatalogRequest{ActorUserID: actorID})
	require.True(t, errors.Is(err, ErrCatalogDisabled))

	_, err = svc.UpdateCatalogState(context.Background(), dto.UpdateCatalogStateRequest{Enabled: true})
	require.NoError(t, err)
	_, err = svc.UpdatePurchaseState(context.Background(), dto.UpdatePurchaseStateRequest{Enabled: false})
	require.NoError(t, err)
	_, err = svc.CreatePurchaseIntent(context.Background(), dto.CreatePurchaseIntentRequest{ActorUserID: actorID, ProductID: "product_avatar_01"})
	require.True(t, errors.Is(err, ErrPurchaseDisabled))
}

func setupActiveProductAndOffer(t *testing.T, svc *ShopService, now time.Time) {
	t.Helper()
	startsAt := now.Add(-24 * time.Hour)
	endsAt := now.Add(24 * time.Hour)

	_, err := svc.UpsertProductDefinition(context.Background(), dto.UpsertProductDefinitionRequest{
		ProductID:       "product_avatar_01",
		Name:            "Avatar Frame",
		Category:        "cosmetic",
		State:           "active",
		InventoryItemID: "avatar_frame_gold",
		SinglePurchase:  true,
	})
	require.NoError(t, err)

	_, err = svc.UpsertOfferDefinition(context.Background(), dto.UpsertOfferDefinitionRequest{
		OfferID:         "offer_avatar_01",
		ProductID:       "product_avatar_01",
		Title:           "Launch Offer",
		Visibility:      "visible",
		PriceMana:       500,
		DiscountPercent: 10,
		Active:          true,
		StartsAt:        &startsAt,
		EndsAt:          &endsAt,
	})
	require.NoError(t, err)
}
