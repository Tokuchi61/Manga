package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/dto"
	paymentrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/repository"
)

func (s *PaymentService) GetOwnWallet(ctx context.Context, request dto.GetOwnWalletRequest) (dto.WalletResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.WalletResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.WalletResponse{}, err
	}
	if err := s.requireTransactionReadEnabled(cfg.TransactionReadEnabled); err != nil {
		return dto.WalletResponse{}, err
	}

	snapshot, err := s.store.GetBalanceSnapshot(ctx, request.ActorUserID)
	if err != nil {
		if errors.Is(err, paymentrepository.ErrNotFound) {
			return dto.WalletResponse{UserID: request.ActorUserID, BalanceMana: 0, UpdatedAt: s.now().UTC()}, nil
		}
		return dto.WalletResponse{}, err
	}

	return dto.WalletResponse{
		UserID:            snapshot.UserID,
		BalanceMana:       snapshot.BalanceMana,
		LastLedgerEntryID: snapshot.LastLedgerEntryID,
		UpdatedAt:         snapshot.UpdatedAt,
	}, nil
}

func (s *PaymentService) ListOwnTransactions(ctx context.Context, request dto.ListOwnTransactionsRequest) (dto.ListOwnTransactionsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListOwnTransactionsResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ListOwnTransactionsResponse{}, err
	}
	if err := s.requireTransactionReadEnabled(cfg.TransactionReadEnabled); err != nil {
		return dto.ListOwnTransactionsResponse{}, err
	}

	transactions, err := s.store.ListTransactionsByUser(ctx, request.ActorUserID, request.Status, request.Limit, request.Offset)
	if err != nil {
		return dto.ListOwnTransactionsResponse{}, err
	}

	items := make([]dto.TransactionResponse, 0, len(transactions))
	for _, tx := range transactions {
		items = append(items, toTransactionResponse(tx))
	}
	return dto.ListOwnTransactionsResponse{Items: items, Count: len(items)}, nil
}
