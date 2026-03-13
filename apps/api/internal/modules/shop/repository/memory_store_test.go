package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStoreProductOfferRoundTrip(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 19, 12, 0, 0, 0, time.UTC)
	start := now.Add(-2 * time.Hour)
	end := now.Add(24 * time.Hour)

	err := store.UpsertProductDefinition(context.Background(), entity.ProductDefinition{
		ProductID:       "product_avatar_01",
		Name:            "Avatar Frame",
		Category:        "cosmetic",
		State:           entity.ProductStateActive,
		InventoryItemID: "avatar_frame_gold",
		SinglePurchase:  true,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	require.NoError(t, err)

	err = store.UpsertOfferDefinition(context.Background(), entity.OfferDefinition{
		OfferID:         "offer_avatar_01",
		ProductID:       "product_avatar_01",
		Title:           "Launch Offer",
		Visibility:      entity.OfferVisibilityVisible,
		PriceMana:       500,
		DiscountPercent: 10,
		Active:          true,
		StartsAt:        &start,
		EndsAt:          &end,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	require.NoError(t, err)

	product, err := store.GetProductDefinition(context.Background(), "product_avatar_01")
	require.NoError(t, err)
	require.Equal(t, "Avatar Frame", product.Name)

	offer, err := store.GetOfferDefinition(context.Background(), "offer_avatar_01")
	require.NoError(t, err)
	require.Equal(t, 500, offer.PriceMana)
}

func TestMemoryStorePurchaseDedupAndSnapshot(t *testing.T) {
	store := NewMemoryStore()
	userID := uuid.NewString()
	now := time.Date(2026, 3, 19, 12, 0, 0, 0, time.UTC)

	intent := entity.PurchaseIntent{
		IntentID:       uuid.NewString(),
		UserID:         userID,
		ProductID:      "product_avatar_01",
		OfferID:        "offer_avatar_01",
		FinalPriceMana: 450,
		Currency:       entity.CurrencyMana,
		Status:         entity.PurchaseStatusDeliveryPending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, store.CreatePurchaseIntent(context.Background(), intent))
	require.NoError(t, store.PutPurchaseDedup(context.Background(), "dedup-shop-purchase", intent))

	_, err := store.GetPurchaseDedup(context.Background(), "dedup-shop-purchase")
	require.NoError(t, err)

	payload, err := store.Snapshot()
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	restored := NewMemoryStore()
	require.NoError(t, restored.RestoreSnapshot(payload))

	readIntent, err := restored.GetPurchaseIntent(context.Background(), intent.IntentID)
	require.NoError(t, err)
	require.Equal(t, intent.FinalPriceMana, readIntent.FinalPriceMana)
}
