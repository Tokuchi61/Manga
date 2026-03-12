package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/entity"
	userrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/repository"
)

func (s *UserService) UpdateProfileVisibility(ctx context.Context, request dto.UpdateProfileVisibilityRequest) (dto.VisibilityResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.VisibilityResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.VisibilityResponse{}, err
	}
	viewerID, err := parseID(request.ViewerID, "viewer_id")
	if err != nil {
		return dto.VisibilityResponse{}, err
	}
	if viewerID != userID {
		return dto.VisibilityResponse{}, ErrProfileNotVisible
	}
	visibility, err := toProfileVisibility(request.ProfileVisibility)
	if err != nil {
		return dto.VisibilityResponse{}, err
	}

	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, userrepository.ErrNotFound) {
			return dto.VisibilityResponse{}, ErrUserNotFound
		}
		return dto.VisibilityResponse{}, err
	}
	if user.AccountState == entity.AccountStateDeactivated {
		return dto.VisibilityResponse{}, ErrAccountDeactivated
	}
	if user.AccountState == entity.AccountStateBanned {
		return dto.VisibilityResponse{}, ErrAccountBanned
	}

	user.ProfileVisibility = visibility
	user.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateUser(ctx, user); err != nil {
		if errors.Is(err, userrepository.ErrNotFound) {
			return dto.VisibilityResponse{}, ErrUserNotFound
		}
		return dto.VisibilityResponse{}, err
	}

	return dto.VisibilityResponse{
		Status:            "profile_visibility_updated",
		ProfileVisibility: string(user.ProfileVisibility),
	}, nil
}

func historyVisibleForViewer(global entity.HistoryVisibilityPreference, entryVisibility entity.HistoryVisibilityPreference, isOwner bool) bool {
	if isOwner {
		return true
	}
	if global == entity.HistoryVisibilityPrivate {
		return false
	}
	return entryVisibility == entity.HistoryVisibilityPublic
}

func (s *UserService) UpdateHistoryVisibility(ctx context.Context, request dto.UpdateHistoryVisibilityRequest) (dto.VisibilityResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.VisibilityResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.VisibilityResponse{}, err
	}
	viewerID, err := parseID(request.ViewerID, "viewer_id")
	if err != nil {
		return dto.VisibilityResponse{}, err
	}
	if viewerID != userID {
		return dto.VisibilityResponse{}, ErrProfileNotVisible
	}
	visibility, err := toHistoryVisibility(request.HistoryVisibilityPreference)
	if err != nil {
		return dto.VisibilityResponse{}, err
	}

	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, userrepository.ErrNotFound) {
			return dto.VisibilityResponse{}, ErrUserNotFound
		}
		return dto.VisibilityResponse{}, err
	}
	if user.AccountState == entity.AccountStateDeactivated {
		return dto.VisibilityResponse{}, ErrAccountDeactivated
	}
	if user.AccountState == entity.AccountStateBanned {
		return dto.VisibilityResponse{}, ErrAccountBanned
	}

	user.HistoryVisibilityPreference = visibility
	user.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateUser(ctx, user); err != nil {
		if errors.Is(err, userrepository.ErrNotFound) {
			return dto.VisibilityResponse{}, ErrUserNotFound
		}
		return dto.VisibilityResponse{}, err
	}

	return dto.VisibilityResponse{
		Status:                      "history_visibility_updated",
		HistoryVisibilityPreference: string(user.HistoryVisibilityPreference),
	}, nil
}
