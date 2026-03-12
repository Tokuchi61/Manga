package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/entity"
	userrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/repository"
	"github.com/google/uuid"
)

func (s *UserService) CreateUser(ctx context.Context, request dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CreateUserResponse{}, err
	}

	credentialID, err := parseID(request.CredentialID, "credential_id")
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	if _, err := s.store.GetUserByCredentialID(ctx, credentialID); err == nil {
		return dto.CreateUserResponse{}, ErrUserAlreadyExists
	} else if !errors.Is(err, userrepository.ErrNotFound) {
		return dto.CreateUserResponse{}, err
	}

	now := s.now().UTC()
	displayName := strings.TrimSpace(request.DisplayName)
	if displayName == "" {
		displayName = request.Username
	}

	user := entity.UserAccount{
		ID:                          uuid.NewString(),
		CredentialID:                credentialID,
		Username:                    request.Username,
		DisplayName:                 displayName,
		Bio:                         strings.TrimSpace(request.Bio),
		ProfileVisibility:           entity.ProfileVisibilityPublic,
		HistoryVisibilityPreference: entity.HistoryVisibilityPrivate,
		AccountState:                entity.AccountStateActive,
		CreatedAt:                   now,
		UpdatedAt:                   now,
	}

	if err := s.store.CreateUser(ctx, user); err != nil {
		if errors.Is(err, userrepository.ErrConflict) {
			return dto.CreateUserResponse{}, ErrUserAlreadyExists
		}
		return dto.CreateUserResponse{}, err
	}

	return dto.CreateUserResponse{
		UserID:       user.ID,
		CredentialID: user.CredentialID,
		Username:     user.Username,
		AccountState: string(user.AccountState),
	}, nil
}
