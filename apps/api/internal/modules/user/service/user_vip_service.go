package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/entity"
	userrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/repository"
)

func (s *UserService) UpdateVIPState(ctx context.Context, request dto.UpdateVIPStateRequest) (dto.VIPStateResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.VIPStateResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.VIPStateResponse{}, err
	}
	action, err := toVIPAction(request.Action)
	if err != nil {
		return dto.VIPStateResponse{}, err
	}

	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, userrepository.ErrNotFound) {
			return dto.VIPStateResponse{}, ErrUserNotFound
		}
		return dto.VIPStateResponse{}, err
	}

	now := s.now().UTC()
	switch action {
	case entity.VIPActionActivate:
		if request.EndsAt == nil {
			return dto.VIPStateResponse{}, ErrValidation
		}
		endAt := request.EndsAt.UTC()
		if !endAt.After(now) {
			return dto.VIPStateResponse{}, ErrInvalidStateTransition
		}
		if user.VIPStartedAt == nil {
			user.VIPStartedAt = &now
		}
		user.VIPEndsAt = &endAt
		user.VIPActive = true
		user.VIPFrozen = false
		user.VIPFrozenAt = nil
		user.VIPFreezeReason = ""
	case entity.VIPActionFreeze:
		if !user.VIPActive || user.VIPFrozen {
			return dto.VIPStateResponse{}, ErrInvalidStateTransition
		}
		user.VIPFrozen = true
		user.VIPFrozenAt = &now
		user.VIPFreezeReason = strings.TrimSpace(request.FreezeReason)
	case entity.VIPActionResume:
		if !user.VIPActive || !user.VIPFrozen {
			return dto.VIPStateResponse{}, ErrInvalidStateTransition
		}
		user.VIPFrozen = false
		user.VIPFrozenAt = nil
		user.VIPFreezeReason = ""
	case entity.VIPActionDeactivate:
		if !user.VIPActive {
			return dto.VIPStateResponse{}, ErrInvalidStateTransition
		}
		user.VIPActive = false
		user.VIPFrozen = false
		user.VIPStartedAt = nil
		user.VIPEndsAt = nil
		user.VIPFrozenAt = nil
		user.VIPFreezeReason = ""
	default:
		return dto.VIPStateResponse{}, ErrValidation
	}

	user.UpdatedAt = now
	if err := s.store.UpdateUser(ctx, user); err != nil {
		if errors.Is(err, userrepository.ErrNotFound) {
			return dto.VIPStateResponse{}, ErrUserNotFound
		}
		return dto.VIPStateResponse{}, err
	}

	return dto.VIPStateResponse{
		Status:    "vip_state_updated",
		VIPActive: user.VIPActive,
		VIPFrozen: user.VIPFrozen,
		VIPEndsAt: user.VIPEndsAt,
	}, nil
}
