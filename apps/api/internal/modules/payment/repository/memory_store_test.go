package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorePackageTransactionLedgerAndSnapshot(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 19, 12, 0, 0, 0, time.UTC)
	userID := uuid.NewString()

	err := store.UpsertManaPackage(context.Background(), entity.ManaPackage{
		PackageID:     "mana_pack_small",
		Name:          "Small Mana Pack",
		ManaAmount:    500,
		PriceAmount:   499,
		PriceCurrency: "TRY",
		Active:        true,
		DisplayOrder:  10,
		Provider:      "mock_provider",
		ProviderSKU:   "sku-small",
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	require.NoError(t, err)

	tx := entity.Transaction{
		TransactionID: uuid.NewString(),
		UserID:        userID,
		PackageID:     "mana_pack_small",
		AmountMana:    500,
		MoneyAmount:   499,
		MoneyCurrency: "TRY",
		Source:        entity.CheckoutSourceExternalProvider,
		Status:        entity.TransactionStatusSuccess,
		Provider:      "mock_provider",
		RequestID:     "req-payment-1",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	require.NoError(t, store.CreateTransaction(context.Background(), tx))
	require.NoError(t, store.PutCheckoutDedup(context.Background(), "dedup-checkout", tx))

	entry := entity.LedgerEntry{
		EntryID:       uuid.NewString(),
		UserID:        userID,
		TransactionID: tx.TransactionID,
		EntryType:     entity.LedgerEntryTypeCredit,
		AmountMana:    500,
		ReasonCode:    "checkout_settled",
		CreatedAt:     now,
	}
	require.NoError(t, store.CreateLedgerEntry(context.Background(), entry))
	require.NoError(t, store.UpsertBalanceSnapshot(context.Background(), entity.BalanceSnapshot{
		UserID:            userID,
		BalanceMana:       500,
		LastLedgerEntryID: entry.EntryID,
		UpdatedAt:         now,
	}))

	payload, err := store.Snapshot()
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	restored := NewMemoryStore()
	require.NoError(t, restored.RestoreSnapshot(payload))

	readTx, err := restored.GetTransaction(context.Background(), tx.TransactionID)
	require.NoError(t, err)
	require.Equal(t, tx.AmountMana, readTx.AmountMana)

	readBalance, err := restored.GetBalanceSnapshot(context.Background(), userID)
	require.NoError(t, err)
	require.Equal(t, 500, readBalance.BalanceMana)
}
