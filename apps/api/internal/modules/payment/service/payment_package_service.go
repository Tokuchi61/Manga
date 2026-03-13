package service

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/dto"
)

func (s *PaymentService) ListManaPackages(ctx context.Context, request dto.ListManaPackagesRequest) (dto.ListManaPackagesResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListManaPackagesResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ListManaPackagesResponse{}, err
	}
	if err := s.requireManaPurchaseEnabled(cfg.ManaPurchaseEnabled); err != nil {
		return dto.ListManaPackagesResponse{}, err
	}

	packages, err := s.store.ListManaPackages(ctx, true, 0, 0)
	if err != nil {
		return dto.ListManaPackagesResponse{}, err
	}

	items := make([]dto.ManaPackageResponse, 0, len(packages))
	for _, pkg := range packages {
		items = append(items, toManaPackageResponse(pkg))
	}
	return dto.ListManaPackagesResponse{Items: items, Count: len(items)}, nil
}
