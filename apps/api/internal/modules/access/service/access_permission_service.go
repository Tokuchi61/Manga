package service

import (
	"context"
	"errors"

	accesscontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/access/contract"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/entity"
	accessrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/access/repository"
	"github.com/google/uuid"
)

func (s *AccessService) CreatePermission(ctx context.Context, request dto.CreatePermissionRequest) (dto.CreatePermissionResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CreatePermissionResponse{}, err
	}

	audienceKind, err := toAudienceKind(request.AudienceKind)
	if err != nil {
		return dto.CreatePermissionResponse{}, err
	}

	now := s.now().UTC()
	permission := entity.Permission{
		ID:           uuid.New(),
		Name:         normalizeName(request.Name),
		Module:       normalizeName(request.Module),
		Surface:      normalizeName(request.Surface),
		Action:       normalizeName(request.Action),
		AudienceKind: audienceKind,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.store.CreatePermission(ctx, permission); err != nil {
		if errors.Is(err, accessrepository.ErrConflict) {
			return dto.CreatePermissionResponse{}, ErrPermissionAlreadyExists
		}
		return dto.CreatePermissionResponse{}, err
	}

	return dto.CreatePermissionResponse{
		PermissionID: permission.ID.String(),
		Name:         permission.Name,
		Module:       permission.Module,
		Surface:      permission.Surface,
		Action:       permission.Action,
		AudienceKind: permission.AudienceKind,
		CreatedAt:    permission.CreatedAt,
	}, nil
}

func (s *AccessService) ListCanonicalPermissions(_ context.Context) (dto.ListCanonicalPermissionsResponse, error) {
	permissions := make([]string, 0, len(accesscontract.CanonicalPermissions))
	permissions = append(permissions, accesscontract.CanonicalPermissions...)
	return dto.ListCanonicalPermissionsResponse{Permissions: permissions}, nil
}
