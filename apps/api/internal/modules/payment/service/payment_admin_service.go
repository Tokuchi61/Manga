package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/entity"
	paymentrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/repository"
	"github.com/google/uuid"
)

func (s *PaymentService) GetRuntimeConfig(ctx context.Context) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return dto.RuntimeConfigResponse{
		ManaPurchaseEnabled:    cfg.ManaPurchaseEnabled,
		CheckoutEnabled:        cfg.CheckoutEnabled,
		TransactionReadEnabled: cfg.TransactionReadEnabled,
		CallbackIntakePaused:   cfg.CallbackIntakePaused,
		UpdatedAt:              cfg.UpdatedAt,
	}, nil
}

func (s *PaymentService) UpdateManaPurchaseState(ctx context.Context, request dto.UpdateManaPurchaseStateRequest) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.ManaPurchaseEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *PaymentService) UpdateCheckoutState(ctx context.Context, request dto.UpdateCheckoutStateRequest) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.CheckoutEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *PaymentService) UpdateTransactionReadState(ctx context.Context, request dto.UpdateTransactionReadStateRequest) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.TransactionReadEnabled = request.Enabled
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *PaymentService) UpdateCallbackIntakeState(ctx context.Context, request dto.UpdateCallbackIntakeStateRequest) (dto.RuntimeConfigResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	cfg.CallbackIntakePaused = request.Paused
	cfg.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateRuntimeConfig(ctx, cfg); err != nil {
		return dto.RuntimeConfigResponse{}, err
	}
	return s.GetRuntimeConfig(ctx)
}

func (s *PaymentService) UpsertManaPackage(ctx context.Context, request dto.UpsertManaPackageRequest) (dto.ManaPackageResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ManaPackageResponse{}, err
	}

	now := s.now().UTC()
	pkg := entity.ManaPackage{
		PackageID:     request.PackageID,
		Name:          request.Name,
		Description:   request.Description,
		ManaAmount:    request.ManaAmount,
		PriceAmount:   request.PriceAmount,
		PriceCurrency: request.PriceCurrency,
		Active:        request.Active,
		DisplayOrder:  request.DisplayOrder,
		Provider:      request.Provider,
		ProviderSKU:   request.ProviderSKU,
		UpdatedAt:     now,
	}

	existing, err := s.store.GetManaPackage(ctx, request.PackageID)
	if err == nil {
		pkg.CreatedAt = existing.CreatedAt
	} else if err != nil && !errors.Is(err, paymentrepository.ErrNotFound) {
		return dto.ManaPackageResponse{}, err
	}
	if pkg.CreatedAt.IsZero() {
		pkg.CreatedAt = now
	}

	if err := s.store.UpsertManaPackage(ctx, pkg); err != nil {
		return dto.ManaPackageResponse{}, err
	}
	return toManaPackageResponse(pkg), nil
}

func (s *PaymentService) ListAdminManaPackages(ctx context.Context, request dto.ListAdminManaPackagesRequest) (dto.ListManaPackagesResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListManaPackagesResponse{}, err
	}

	packages, err := s.store.ListManaPackages(ctx, request.ActiveOnly, request.Limit, request.Offset)
	if err != nil {
		return dto.ListManaPackagesResponse{}, err
	}

	items := make([]dto.ManaPackageResponse, 0, len(packages))
	for _, pkg := range packages {
		items = append(items, toManaPackageResponse(pkg))
	}
	return dto.ListManaPackagesResponse{Items: items, Count: len(items)}, nil
}

func (s *PaymentService) RunReconcile(ctx context.Context, request dto.RunReconcileRequest) (dto.RunReconcileResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RunReconcileResponse{}, err
	}

	userIDs := make([]string, 0)
	if request.ActorUserID != "" {
		userIDs = append(userIDs, request.ActorUserID)
	} else {
		ledgerUsers, err := s.store.ListLedgerUsers(ctx)
		if err != nil {
			return dto.RunReconcileResponse{}, err
		}
		userIDs = append(userIDs, ledgerUsers...)
	}

	corrected := 0
	now := s.now().UTC()
	for _, userID := range userIDs {
		entries, err := s.store.ListLedgerEntriesByUser(ctx, userID)
		if err != nil {
			return dto.RunReconcileResponse{}, err
		}
		expected := calculateLedgerBalance(entries)

		snapshot, err := s.resolveWalletBalance(ctx, userID)
		if err != nil {
			return dto.RunReconcileResponse{}, err
		}
		if snapshot.BalanceMana == expected {
			continue
		}
		snapshot.BalanceMana = expected
		snapshot.UpdatedAt = now
		if len(entries) > 0 {
			snapshot.LastLedgerEntryID = entries[len(entries)-1].EntryID
		}
		if err := s.store.UpsertBalanceSnapshot(ctx, snapshot); err != nil {
			return dto.RunReconcileResponse{}, err
		}
		corrected++
	}

	return dto.RunReconcileResponse{ScannedUsers: len(userIDs), CorrectedUsers: corrected}, nil
}

func (s *PaymentService) ProcessRefund(ctx context.Context, request dto.ProcessRefundRequest) (dto.AdminTransactionActionResponse, error) {
	return s.processDebitAction(ctx, request.TransactionID, request.ReasonCode, entity.TransactionStatusRefunded, "refund_completed")
}

func (s *PaymentService) ProcessReversal(ctx context.Context, request dto.ProcessReversalRequest) (dto.AdminTransactionActionResponse, error) {
	return s.processDebitAction(ctx, request.TransactionID, request.ReasonCode, entity.TransactionStatusReversed, "reversal_completed")
}

func (s *PaymentService) processDebitAction(ctx context.Context, transactionID string, reasonCode string, targetStatus string, statusLabel string) (dto.AdminTransactionActionResponse, error) {
	tx, err := s.store.GetTransaction(ctx, transactionID)
	if err != nil {
		if errors.Is(err, paymentrepository.ErrNotFound) {
			return dto.AdminTransactionActionResponse{}, ErrNotFound
		}
		return dto.AdminTransactionActionResponse{}, err
	}

	currentStatus := normalizeValue(tx.Status)
	if currentStatus != entity.TransactionStatusSuccess {
		return dto.AdminTransactionActionResponse{}, ErrConflict
	}

	now := s.now().UTC()
	tx.Status = targetStatus
	tx.UpdatedAt = now
	if err := s.store.UpdateTransaction(ctx, tx); err != nil {
		return dto.AdminTransactionActionResponse{}, err
	}

	entry := entity.LedgerEntry{
		EntryID:       uuid.NewString(),
		UserID:        tx.UserID,
		TransactionID: tx.TransactionID,
		EntryType:     entity.LedgerEntryTypeDebit,
		AmountMana:    tx.AmountMana,
		ReasonCode:    reasonCode,
		CreatedAt:     now,
	}
	if err := s.store.CreateLedgerEntry(ctx, entry); err != nil {
		return dto.AdminTransactionActionResponse{}, err
	}

	snapshot, err := s.resolveWalletBalance(ctx, tx.UserID)
	if err != nil {
		return dto.AdminTransactionActionResponse{}, err
	}
	snapshot.BalanceMana -= tx.AmountMana
	snapshot.LastLedgerEntryID = entry.EntryID
	snapshot.UpdatedAt = now
	if err := s.store.UpsertBalanceSnapshot(ctx, snapshot); err != nil {
		return dto.AdminTransactionActionResponse{}, err
	}

	return dto.AdminTransactionActionResponse{
		Status:        statusLabel,
		TransactionID: tx.TransactionID,
		BalanceMana:   snapshot.BalanceMana,
		UpdatedAt:     now,
	}, nil
}
