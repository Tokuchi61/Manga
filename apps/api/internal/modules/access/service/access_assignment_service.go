package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/entity"
	accessrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/access/repository"
	"github.com/google/uuid"
)

func (s *AccessService) AssignPermissionToRole(ctx context.Context, roleID string, request dto.AssignPermissionToRoleRequest) (dto.AssignPermissionToRoleResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.AssignPermissionToRoleResponse{}, err
	}

	parsedRoleID, err := parseUUID(roleID, "role_id")
	if err != nil {
		return dto.AssignPermissionToRoleResponse{}, err
	}

	permission, err := s.store.GetPermissionByName(ctx, request.PermissionName)
	if err != nil {
		if errors.Is(err, accessrepository.ErrNotFound) {
			return dto.AssignPermissionToRoleResponse{}, ErrPermissionNotFound
		}
		return dto.AssignPermissionToRoleResponse{}, err
	}

	if err := s.store.AttachPermissionToRole(ctx, parsedRoleID, permission.ID, s.now().UTC()); err != nil {
		if errors.Is(err, accessrepository.ErrNotFound) {
			return dto.AssignPermissionToRoleResponse{}, ErrRoleNotFound
		}
		if errors.Is(err, accessrepository.ErrConflict) {
			return dto.AssignPermissionToRoleResponse{}, ErrPermissionAlreadyExists
		}
		return dto.AssignPermissionToRoleResponse{}, err
	}

	return dto.AssignPermissionToRoleResponse{
		Status:         "role_permission_assigned",
		RoleID:         parsedRoleID.String(),
		PermissionName: permission.Name,
	}, nil
}

func (s *AccessService) AssignRoleToUser(ctx context.Context, userID string, request dto.AssignRoleToUserRequest) (dto.AssignRoleToUserResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.AssignRoleToUserResponse{}, err
	}

	parsedUserID, err := parseUUID(userID, "user_id")
	if err != nil {
		return dto.AssignRoleToUserResponse{}, err
	}

	role, err := s.store.GetRoleByName(ctx, request.RoleName)
	if err != nil {
		if errors.Is(err, accessrepository.ErrNotFound) {
			return dto.AssignRoleToUserResponse{}, ErrRoleNotFound
		}
		return dto.AssignRoleToUserResponse{}, err
	}

	if request.ExpiresAt != nil {
		expiresAt := request.ExpiresAt.UTC()
		if !expiresAt.After(s.now().UTC()) {
			return dto.AssignRoleToUserResponse{}, ErrValidation
		}
	}

	assignment := entity.UserRole{
		UserID:    parsedUserID.String(),
		RoleID:    role.ID,
		ExpiresAt: request.ExpiresAt,
		CreatedAt: s.now().UTC(),
	}
	if err := s.store.AssignRoleToUser(ctx, assignment); err != nil {
		if errors.Is(err, accessrepository.ErrNotFound) {
			return dto.AssignRoleToUserResponse{}, ErrRoleNotFound
		}
		if errors.Is(err, accessrepository.ErrConflict) {
			return dto.AssignRoleToUserResponse{}, ErrRoleAlreadyExists
		}
		return dto.AssignRoleToUserResponse{}, err
	}

	return dto.AssignRoleToUserResponse{
		Status:    "user_role_assigned",
		UserID:    assignment.UserID,
		RoleName:  role.Name,
		ExpiresAt: assignment.ExpiresAt,
	}, nil
}

func (s *AccessService) CreateTemporaryGrant(ctx context.Context, userID string, request dto.CreateTemporaryGrantRequest) (dto.CreateTemporaryGrantResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CreateTemporaryGrantResponse{}, err
	}

	parsedUserID, err := parseUUID(userID, "user_id")
	if err != nil {
		return dto.CreateTemporaryGrantResponse{}, err
	}

	permission, err := s.store.GetPermissionByName(ctx, request.PermissionName)
	if err != nil {
		if errors.Is(err, accessrepository.ErrNotFound) {
			return dto.CreateTemporaryGrantResponse{}, ErrPermissionNotFound
		}
		return dto.CreateTemporaryGrantResponse{}, err
	}

	expiresAt := request.ExpiresAt.UTC()
	if !expiresAt.After(s.now().UTC()) {
		return dto.CreateTemporaryGrantResponse{}, ErrValidation
	}

	grant := entity.TemporaryGrant{
		ID:           uuid.New(),
		UserID:       parsedUserID.String(),
		PermissionID: permission.ID,
		Reason:       request.Reason,
		ExpiresAt:    expiresAt,
		CreatedAt:    s.now().UTC(),
	}
	if err := s.store.CreateTemporaryGrant(ctx, grant); err != nil {
		if errors.Is(err, accessrepository.ErrNotFound) {
			return dto.CreateTemporaryGrantResponse{}, ErrPermissionNotFound
		}
		return dto.CreateTemporaryGrantResponse{}, err
	}

	return dto.CreateTemporaryGrantResponse{
		Status:         "temporary_grant_created",
		GrantID:        grant.ID.String(),
		UserID:         grant.UserID,
		PermissionName: permission.Name,
		ExpiresAt:      grant.ExpiresAt,
	}, nil
}
