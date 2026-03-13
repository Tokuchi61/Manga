package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/entity"
	paymentrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/repository"
	"github.com/google/uuid"
)

func (s *PaymentService) ProcessProviderCallback(ctx context.Context, request dto.ProcessProviderCallbackRequest) (dto.ProcessProviderCallbackResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ProcessProviderCallbackResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ProcessProviderCallbackResponse{}, err
	}
	if err := s.requireCallbackIntakeOpen(cfg.CallbackIntakePaused); err != nil {
		return dto.ProcessProviderCallbackResponse{}, err
	}

	now := s.now().UTC()
	dedupTx, err := s.store.GetCallbackDedup(ctx, request.ProviderEventID)
	if err == nil {
		return dto.ProcessProviderCallbackResponse{
			Status:            "idempotent",
			TransactionID:     dedupTx.TransactionID,
			TransactionStatus: dedupTx.Status,
			ProcessedAt:       now,
		}, nil
	}
	if err != nil && !errors.Is(err, paymentrepository.ErrNotFound) {
		return dto.ProcessProviderCallbackResponse{}, err
	}

	session, err := s.store.GetProviderSession(ctx, request.SessionID)
	if err != nil {
		if errors.Is(err, paymentrepository.ErrNotFound) {
			return dto.ProcessProviderCallbackResponse{}, ErrNotFound
		}
		return dto.ProcessProviderCallbackResponse{}, err
	}

	tx, err := s.store.GetTransaction(ctx, session.TransactionID)
	if err != nil {
		if errors.Is(err, paymentrepository.ErrNotFound) {
			return dto.ProcessProviderCallbackResponse{}, ErrNotFound
		}
		return dto.ProcessProviderCallbackResponse{}, err
	}

	normalizedStatus := normalizeValue(request.Status)
	tx.ProviderReference = request.ProviderReference
	tx.UpdatedAt = now
	session.ProviderReference = request.ProviderReference
	session.UpdatedAt = now

	switch normalizedStatus {
	case "success":
		session.Status = entity.ProviderSessionStatusCompleted
		if normalizeValue(tx.Status) == entity.TransactionStatusPending {
			tx.Status = entity.TransactionStatusSuccess

			entry := entity.LedgerEntry{
				EntryID:       uuid.NewString(),
				UserID:        tx.UserID,
				TransactionID: tx.TransactionID,
				EntryType:     entity.LedgerEntryTypeCredit,
				AmountMana:    tx.AmountMana,
				ReasonCode:    "checkout_settled",
				CreatedAt:     now,
			}
			if err := s.store.CreateLedgerEntry(ctx, entry); err != nil {
				return dto.ProcessProviderCallbackResponse{}, err
			}

			snapshot, snapshotErr := s.resolveWalletBalance(ctx, tx.UserID)
			if snapshotErr != nil {
				return dto.ProcessProviderCallbackResponse{}, snapshotErr
			}
			snapshot.BalanceMana += tx.AmountMana
			snapshot.LastLedgerEntryID = entry.EntryID
			snapshot.UpdatedAt = now
			if err := s.store.UpsertBalanceSnapshot(ctx, snapshot); err != nil {
				return dto.ProcessProviderCallbackResponse{}, err
			}
		}
	case "failed":
		tx.Status = entity.TransactionStatusFailed
		session.Status = entity.ProviderSessionStatusCancelled
	case "cancelled":
		tx.Status = entity.TransactionStatusCancelled
		session.Status = entity.ProviderSessionStatusCancelled
	default:
		return dto.ProcessProviderCallbackResponse{}, ErrValidation
	}

	if err := s.store.UpdateTransaction(ctx, tx); err != nil {
		return dto.ProcessProviderCallbackResponse{}, err
	}
	if err := s.store.UpdateProviderSession(ctx, session); err != nil {
		return dto.ProcessProviderCallbackResponse{}, err
	}
	if err := s.store.PutCallbackDedup(ctx, request.ProviderEventID, tx); err != nil {
		return dto.ProcessProviderCallbackResponse{}, err
	}
	if err := s.publishCallbackEvent(ctx, request, tx, now); err != nil {
		return dto.ProcessProviderCallbackResponse{}, err
	}

	return dto.ProcessProviderCallbackResponse{
		Status:            "callback_accepted",
		TransactionID:     tx.TransactionID,
		TransactionStatus: tx.Status,
		ProcessedAt:       now,
	}, nil
}
