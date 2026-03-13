package service

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/dto"
	paymentrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPaymentServiceCheckoutCallbackWalletAndRefundFlow(t *testing.T) {
	store := paymentrepository.NewMemoryStore()
	svc := New(store, validation.New())
	now := time.Date(2026, 3, 20, 14, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }

	actorID := uuid.NewString()
	setupActivePackage(t, svc)

	checkoutRes, err := svc.StartCheckoutSession(context.Background(), dto.StartCheckoutSessionRequest{
		ActorUserID: actorID,
		PackageID:   "mana_pack_small",
		RequestID:   "req-payment-1",
	})
	require.NoError(t, err)
	require.Equal(t, "checkout_started", checkoutRes.Status)

	idempotentRes, err := svc.StartCheckoutSession(context.Background(), dto.StartCheckoutSessionRequest{
		ActorUserID: actorID,
		PackageID:   "mana_pack_small",
		RequestID:   "req-payment-1",
	})
	require.NoError(t, err)
	require.Equal(t, "idempotent", idempotentRes.Status)

	callbackRes, err := svc.ProcessProviderCallback(context.Background(), dto.ProcessProviderCallbackRequest{
		ProviderEventID:   "evt-payment-1",
		SessionID:         checkoutRes.SessionID,
		ProviderReference: "provider-ref-1",
		Status:            "success",
	})
	require.NoError(t, err)
	require.Equal(t, "callback_accepted", callbackRes.Status)
	require.Equal(t, "success", callbackRes.TransactionStatus)

	callbackIdempotentRes, err := svc.ProcessProviderCallback(context.Background(), dto.ProcessProviderCallbackRequest{
		ProviderEventID: "evt-payment-1",
		SessionID:       checkoutRes.SessionID,
		Status:          "success",
	})
	require.NoError(t, err)
	require.Equal(t, "idempotent", callbackIdempotentRes.Status)

	walletRes, err := svc.GetOwnWallet(context.Background(), dto.GetOwnWalletRequest{ActorUserID: actorID})
	require.NoError(t, err)
	require.Equal(t, 500, walletRes.BalanceMana)

	transactionsRes, err := svc.ListOwnTransactions(context.Background(), dto.ListOwnTransactionsRequest{ActorUserID: actorID})
	require.NoError(t, err)
	require.Equal(t, 1, transactionsRes.Count)
	require.Equal(t, "success", transactionsRes.Items[0].Status)

	refundRes, err := svc.ProcessRefund(context.Background(), dto.ProcessRefundRequest{
		TransactionID: checkoutRes.TransactionID,
		ReasonCode:    "manual_refund",
	})
	require.NoError(t, err)
	require.Equal(t, "refund_completed", refundRes.Status)
	require.Equal(t, 0, refundRes.BalanceMana)
}

func TestPaymentServiceRuntimeTogglesAffectFlows(t *testing.T) {
	store := paymentrepository.NewMemoryStore()
	svc := New(store, validation.New())
	actorID := uuid.NewString()
	setupActivePackage(t, svc)

	_, err := svc.UpdateCheckoutState(context.Background(), dto.UpdateCheckoutStateRequest{Enabled: false})
	require.NoError(t, err)

	_, err = svc.StartCheckoutSession(context.Background(), dto.StartCheckoutSessionRequest{
		ActorUserID: actorID,
		PackageID:   "mana_pack_small",
		RequestID:   "req-disabled-checkout",
	})
	require.True(t, errors.Is(err, ErrCheckoutDisabled))

	_, err = svc.UpdateCheckoutState(context.Background(), dto.UpdateCheckoutStateRequest{Enabled: true})
	require.NoError(t, err)
	_, err = svc.UpdateTransactionReadState(context.Background(), dto.UpdateTransactionReadStateRequest{Enabled: false})
	require.NoError(t, err)

	_, err = svc.GetOwnWallet(context.Background(), dto.GetOwnWalletRequest{ActorUserID: actorID})
	require.True(t, errors.Is(err, ErrTransactionReadDisabled))

	_, err = svc.UpdateCallbackIntakeState(context.Background(), dto.UpdateCallbackIntakeStateRequest{Paused: true})
	require.NoError(t, err)

	_, err = svc.ProcessProviderCallback(context.Background(), dto.ProcessProviderCallbackRequest{
		ProviderEventID: "evt-paused",
		SessionID:       "session-any",
		Status:          "success",
	})
	require.True(t, errors.Is(err, ErrCallbackIntakePaused))
}

func TestVerifyProviderCallbackSignature(t *testing.T) {
	timestampNow := time.Date(2026, 3, 20, 14, 0, 0, 0, time.UTC)
	svc := New(paymentrepository.NewMemoryStore(), validation.New(), Config{
		CallbackSigningSecret: "test-callback-secret",
		CallbackTimestampSkew: 5 * time.Minute,
	})
	svc.now = func() time.Time { return timestampNow }

	payload := []byte(`{"provider_event_id":"evt-1","status":"success"}`)
	timestamp := strconv.FormatInt(timestampNow.Unix(), 10)
	signature := SignProviderCallback("test-callback-secret", timestamp, payload)

	require.NoError(t, svc.VerifyProviderCallbackSignature(payload, signature, timestamp))
	require.ErrorIs(t, svc.VerifyProviderCallbackSignature(payload, "bad-signature", timestamp), ErrCallbackSignature)
	require.ErrorIs(t, svc.VerifyProviderCallbackSignature(payload, signature, strconv.FormatInt(timestampNow.Add(-10*time.Minute).Unix(), 10)), ErrCallbackTimestamp)
}

func setupActivePackage(t *testing.T, svc *PaymentService) {
	t.Helper()

	_, err := svc.UpsertManaPackage(context.Background(), dto.UpsertManaPackageRequest{
		PackageID:     "mana_pack_small",
		Name:          "Small Mana Pack",
		Description:   "500 mana",
		ManaAmount:    500,
		PriceAmount:   499,
		PriceCurrency: "TRY",
		Active:        true,
		DisplayOrder:  10,
		Provider:      "mock_provider",
		ProviderSKU:   "sku-small",
	})
	require.NoError(t, err)
}
