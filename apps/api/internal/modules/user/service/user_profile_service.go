package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/entity"
	userrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/repository"
)

func (s *UserService) GetPublicProfile(ctx context.Context, request dto.GetPublicProfileRequest) (dto.PublicProfileResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.PublicProfileResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.PublicProfileResponse{}, err
	}

	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, userrepository.ErrNotFound) {
			return dto.PublicProfileResponse{}, ErrUserNotFound
		}
		return dto.PublicProfileResponse{}, err
	}

	if user.AccountState == entity.AccountStateBanned || user.ProfileVisibility == entity.ProfileVisibilityPrivate {
		return dto.PublicProfileResponse{}, ErrProfileNotVisible
	}

	return dto.PublicProfileResponse{
		UserID:            user.ID,
		Username:          user.Username,
		DisplayName:       user.DisplayName,
		Bio:               user.Bio,
		AvatarURL:         user.AvatarURL,
		BannerURL:         user.BannerURL,
		ProfileVisibility: string(user.ProfileVisibility),
		VIPBadgeVisible:   user.VIPActive,
	}, nil
}

func (s *UserService) GetOwnProfile(ctx context.Context, request dto.GetOwnProfileRequest) (dto.OwnProfileResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OwnProfileResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.OwnProfileResponse{}, err
	}
	viewerID, err := parseID(request.ViewerID, "viewer_id")
	if err != nil {
		return dto.OwnProfileResponse{}, err
	}
	if viewerID != userID {
		return dto.OwnProfileResponse{}, ErrProfileNotVisible
	}

	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, userrepository.ErrNotFound) {
			return dto.OwnProfileResponse{}, ErrUserNotFound
		}
		return dto.OwnProfileResponse{}, err
	}

	return dto.OwnProfileResponse{
		UserID:                      user.ID,
		CredentialID:                user.CredentialID,
		Username:                    user.Username,
		DisplayName:                 user.DisplayName,
		Bio:                         user.Bio,
		AvatarURL:                   user.AvatarURL,
		BannerURL:                   user.BannerURL,
		ProfileVisibility:           string(user.ProfileVisibility),
		HistoryVisibilityPreference: string(user.HistoryVisibilityPreference),
		AccountState:                string(user.AccountState),
		VIPActive:                   user.VIPActive,
		VIPFrozen:                   user.VIPFrozen,
		VIPStartedAt:                user.VIPStartedAt,
		VIPEndsAt:                   user.VIPEndsAt,
		VIPFrozenAt:                 user.VIPFrozenAt,
		VIPFreezeReason:             user.VIPFreezeReason,
	}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, request dto.UpdateProfileRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	viewerID, err := parseID(request.ViewerID, "viewer_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	if viewerID != userID {
		return dto.OperationResponse{}, ErrProfileNotVisible
	}

	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, userrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrUserNotFound
		}
		return dto.OperationResponse{}, err
	}
	if user.AccountState == entity.AccountStateDeactivated {
		return dto.OperationResponse{}, ErrAccountDeactivated
	}
	if user.AccountState == entity.AccountStateBanned {
		return dto.OperationResponse{}, ErrAccountBanned
	}

	if request.DisplayName != nil {
		user.DisplayName = strings.TrimSpace(*request.DisplayName)
	}
	if request.Bio != nil {
		user.Bio = strings.TrimSpace(*request.Bio)
	}
	if request.AvatarURL != nil {
		user.AvatarURL = strings.TrimSpace(*request.AvatarURL)
	}
	if request.BannerURL != nil {
		user.BannerURL = strings.TrimSpace(*request.BannerURL)
	}

	user.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateUser(ctx, user); err != nil {
		if errors.Is(err, userrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrUserNotFound
		}
		if errors.Is(err, userrepository.ErrConflict) {
			return dto.OperationResponse{}, ErrUserAlreadyExists
		}
		return dto.OperationResponse{}, err
	}

	return dto.OperationResponse{Status: "profile_updated"}, nil
}
