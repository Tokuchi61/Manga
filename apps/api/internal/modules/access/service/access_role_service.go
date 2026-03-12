package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/entity"
	accessrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/access/repository"
	"github.com/google/uuid"
)

func (s *AccessService) CreateRole(ctx context.Context, request dto.CreateRoleRequest) (dto.CreateRoleResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CreateRoleResponse{}, err
	}

	now := s.now().UTC()
	role := entity.Role{
		ID:           uuid.New(),
		Name:         normalizeName(request.Name),
		Priority:     request.Priority,
		IsDefault:    request.IsDefault,
		IsSuperAdmin: request.IsSuperAdmin,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.store.CreateRole(ctx, role); err != nil {
		if errors.Is(err, accessrepository.ErrConflict) {
			return dto.CreateRoleResponse{}, ErrRoleAlreadyExists
		}
		return dto.CreateRoleResponse{}, err
	}

	return dto.CreateRoleResponse{
		RoleID:       role.ID.String(),
		Name:         role.Name,
		Priority:     role.Priority,
		IsDefault:    role.IsDefault,
		IsSuperAdmin: role.IsSuperAdmin,
		CreatedAt:    role.CreatedAt,
	}, nil
}
