package service

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/entity"
)

func (s *NotificationService) GetPreference(ctx context.Context, request dto.GetPreferenceRequest) (dto.PreferenceResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.PreferenceResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.PreferenceResponse{}, err
	}

	preference, err := s.store.GetPreferenceByUserID(ctx, userID)
	if err != nil {
		return dto.PreferenceResponse{}, err
	}

	return mapPreference(preference), nil
}

func (s *NotificationService) UpdatePreference(ctx context.Context, request dto.UpdatePreferenceRequest) (dto.PreferenceResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.PreferenceResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.PreferenceResponse{}, err
	}

	mutedCategories, err := parseMutedCategories(request.MutedCategories)
	if err != nil {
		return dto.PreferenceResponse{}, err
	}
	if request.QuietHoursEnabled && request.QuietHoursStart == request.QuietHoursEnd {
		return dto.PreferenceResponse{}, ErrValidation
	}

	now := s.now().UTC()
	preference := entity.Preference{
		UserID:            userID,
		MutedCategories:   mutedCategories,
		QuietHoursEnabled: request.QuietHoursEnabled,
		QuietHoursStart:   request.QuietHoursStart,
		QuietHoursEnd:     request.QuietHoursEnd,
		InAppEnabled:      request.InAppEnabled,
		EmailEnabled:      request.EmailEnabled,
		PushEnabled:       request.PushEnabled,
		DigestEnabled:     request.DigestEnabled,
		UpdatedAt:         now,
	}

	persisted, err := s.store.UpsertPreference(ctx, preference)
	if err != nil {
		return dto.PreferenceResponse{}, err
	}

	return mapPreference(persisted), nil
}
