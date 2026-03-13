package service

import (
	"context"
	"errors"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/entity"
	paymentrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/repository"
	"github.com/google/uuid"
)

func (s *PaymentService) StartCheckoutSession(ctx context.Context, request dto.StartCheckoutSessionRequest) (dto.StartCheckoutSessionResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.StartCheckoutSessionResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.StartCheckoutSessionResponse{}, err
	}
	if err := s.requireManaPurchaseEnabled(cfg.ManaPurchaseEnabled); err != nil {
		return dto.StartCheckoutSessionResponse{}, err
	}
	if err := s.requireCheckoutEnabled(cfg.CheckoutEnabled); err != nil {
		return dto.StartCheckoutSessionResponse{}, err
	}

	now := s.now().UTC()
	dedupKey := buildCheckoutDedupKey(request.ActorUserID, request.RequestID)

	dedupTx, err := s.store.GetCheckoutDedup(ctx, dedupKey)
	if err == nil {
		session, sessionErr := s.store.GetProviderSession(ctx, dedupTx.TransactionID)
		if sessionErr != nil {
			return dto.StartCheckoutSessionResponse{}, sessionErr
		}
		return dto.StartCheckoutSessionResponse{
			Status:        "idempotent",
			TransactionID: dedupTx.TransactionID,
			SessionID:     session.SessionID,
			Provider:      dedupTx.Provider,
			ExpiresAt:     session.ExpiresAt,
			AmountMana:    dedupTx.AmountMana,
			PriceAmount:   dedupTx.MoneyAmount,
			PriceCurrency: dedupTx.MoneyCurrency,
		}, nil
	}
	if err != nil && !errors.Is(err, paymentrepository.ErrNotFound) {
		return dto.StartCheckoutSessionResponse{}, err
	}

	pkg, err := s.store.GetManaPackage(ctx, request.PackageID)
	if err != nil {
		if errors.Is(err, paymentrepository.ErrNotFound) {
			return dto.StartCheckoutSessionResponse{}, ErrNotFound
		}
		return dto.StartCheckoutSessionResponse{}, err
	}
	if !pkg.Active {
		return dto.StartCheckoutSessionResponse{}, ErrNotFound
	}

	source := normalizeValue(request.Source)
	if source == "" {
		source = entity.CheckoutSourceExternalProvider
	}

	transactionID := uuid.NewString()
	sessionID := transactionID
	expiresAt := now.Add(15 * time.Minute)

	tx := entity.Transaction{
		TransactionID: transactionID,
		UserID:        request.ActorUserID,
		PackageID:     pkg.PackageID,
		AmountMana:    pkg.ManaAmount,
		MoneyAmount:   pkg.PriceAmount,
		MoneyCurrency: pkg.PriceCurrency,
		Source:        source,
		Status:        entity.TransactionStatusPending,
		Provider:      pkg.Provider,
		RequestID:     request.RequestID,
		CorrelationID: request.CorrelationID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	session := entity.ProviderSession{
		SessionID:         sessionID,
		TransactionID:     transactionID,
		UserID:            request.ActorUserID,
		PackageID:         pkg.PackageID,
		Provider:          pkg.Provider,
		ProviderReference: "pending",
		Status:            entity.ProviderSessionStatusStarted,
		ExpiresAt:         &expiresAt,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := s.store.CreateTransaction(ctx, tx); err != nil {
		return dto.StartCheckoutSessionResponse{}, err
	}
	if err := s.store.CreateProviderSession(ctx, session); err != nil {
		return dto.StartCheckoutSessionResponse{}, err
	}
	if err := s.store.PutCheckoutDedup(ctx, dedupKey, tx); err != nil {
		return dto.StartCheckoutSessionResponse{}, err
	}

	return dto.StartCheckoutSessionResponse{
		Status:        "checkout_started",
		TransactionID: tx.TransactionID,
		SessionID:     session.SessionID,
		Provider:      tx.Provider,
		ExpiresAt:     session.ExpiresAt,
		AmountMana:    tx.AmountMana,
		PriceAmount:   tx.MoneyAmount,
		PriceCurrency: tx.MoneyCurrency,
	}, nil
}
