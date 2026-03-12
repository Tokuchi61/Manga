package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/entity"
	userrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/repository"
)

func (s *UserService) UpdateAccountState(ctx context.Context, request dto.UpdateAccountStateRequest) (dto.AccountStateResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.AccountStateResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.AccountStateResponse{}, err
	}
	targetState, err := toAccountState(request.AccountState)
	if err != nil {
		return dto.AccountStateResponse{}, err
	}

	actorScope := strings.TrimSpace(strings.ToLower(request.ActorScope))
	if actorScope == "self" {
		actorUserID, parseErr := parseID(request.ActorUserID, "actor_user_id")
		if parseErr != nil {
			return dto.AccountStateResponse{}, parseErr
		}
		if actorUserID != userID {
			return dto.AccountStateResponse{}, ErrProfileNotVisible
		}
		if targetState == entity.AccountStateBanned {
			return dto.AccountStateResponse{}, ErrInvalidStateTransition
		}
	}

	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, userrepository.ErrNotFound) {
			return dto.AccountStateResponse{}, ErrUserNotFound
		}
		return dto.AccountStateResponse{}, err
	}

	user.AccountState = targetState
	user.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateUser(ctx, user); err != nil {
		if errors.Is(err, userrepository.ErrNotFound) {
			return dto.AccountStateResponse{}, ErrUserNotFound
		}
		return dto.AccountStateResponse{}, err
	}

	return dto.AccountStateResponse{
		Status:       "account_state_updated",
		AccountState: string(user.AccountState),
	}, nil
}
