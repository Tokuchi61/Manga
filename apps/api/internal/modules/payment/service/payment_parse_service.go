package service

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/entity"
	paymentrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/repository"
)

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func buildCheckoutDedupKey(userID string, requestID string) string {
	return normalizeValue(userID) + ":" + normalizeValue(requestID)
}

func toManaPackageResponse(pkg entity.ManaPackage) dto.ManaPackageResponse {
	return dto.ManaPackageResponse{
		PackageID:     pkg.PackageID,
		Name:          pkg.Name,
		Description:   pkg.Description,
		ManaAmount:    pkg.ManaAmount,
		PriceAmount:   pkg.PriceAmount,
		PriceCurrency: pkg.PriceCurrency,
		Active:        pkg.Active,
		DisplayOrder:  pkg.DisplayOrder,
		Provider:      pkg.Provider,
		ProviderSKU:   pkg.ProviderSKU,
		UpdatedAt:     pkg.UpdatedAt,
	}
}

func toTransactionResponse(tx entity.Transaction) dto.TransactionResponse {
	return dto.TransactionResponse{
		TransactionID:     tx.TransactionID,
		PackageID:         tx.PackageID,
		AmountMana:        tx.AmountMana,
		MoneyAmount:       tx.MoneyAmount,
		MoneyCurrency:     tx.MoneyCurrency,
		Source:            tx.Source,
		Status:            tx.Status,
		Provider:          tx.Provider,
		ProviderReference: tx.ProviderReference,
		RequestID:         tx.RequestID,
		CorrelationID:     tx.CorrelationID,
		UpdatedAt:         tx.UpdatedAt,
	}
}

func (s *PaymentService) resolveWalletBalance(ctx context.Context, userID string) (entity.BalanceSnapshot, error) {
	snapshot, err := s.store.GetBalanceSnapshot(ctx, userID)
	if err == nil {
		return snapshot, nil
	}
	if !errors.Is(err, paymentrepository.ErrNotFound) {
		return entity.BalanceSnapshot{}, err
	}

	now := s.now().UTC()
	snapshot = entity.BalanceSnapshot{
		UserID:      userID,
		BalanceMana: 0,
		UpdatedAt:   now,
	}
	if upsertErr := s.store.UpsertBalanceSnapshot(ctx, snapshot); upsertErr != nil {
		return entity.BalanceSnapshot{}, upsertErr
	}
	return snapshot, nil
}

func sortPackages(items []entity.ManaPackage) {
	sort.Slice(items, func(i int, j int) bool {
		if items[i].DisplayOrder == items[j].DisplayOrder {
			return items[i].UpdatedAt.After(items[j].UpdatedAt)
		}
		return items[i].DisplayOrder < items[j].DisplayOrder
	})
}

func calculateLedgerBalance(entries []entity.LedgerEntry) int {
	balance := 0
	for _, entry := range entries {
		switch normalizeValue(entry.EntryType) {
		case entity.LedgerEntryTypeCredit:
			balance += entry.AmountMana
		case entity.LedgerEntryTypeDebit:
			balance -= entry.AmountMana
		}
	}
	return balance
}
