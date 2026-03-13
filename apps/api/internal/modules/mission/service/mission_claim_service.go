package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/dto"
	missionrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/repository"
)

func (s *MissionService) ClaimMission(ctx context.Context, request dto.ClaimMissionRequest) (dto.ClaimMissionResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ClaimMissionResponse{}, err
	}
	request.RequestID = strings.TrimSpace(request.RequestID)
	if request.RequestID == "" {
		return dto.ClaimMissionResponse{}, ErrValidation
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ClaimMissionResponse{}, err
	}
	if err := s.requireClaimEnabled(cfg.ClaimEnabled); err != nil {
		return dto.ClaimMissionResponse{}, err
	}

	definition, err := s.store.GetMissionDefinition(ctx, request.MissionID)
	if err != nil {
		if errors.Is(err, missionrepository.ErrNotFound) {
			return dto.ClaimMissionResponse{}, ErrNotFound
		}
		return dto.ClaimMissionResponse{}, err
	}

	now := s.now().UTC()
	if !isMissionActive(definition, now) {
		return dto.ClaimMissionResponse{}, ErrMissionInactive
	}

	periodKey := buildPeriodKey(definition.Category, now, cfg.DailyResetHourUTC)
	dedupKey := buildClaimDedupKey(request.ActorUserID, request.MissionID, periodKey, request.RequestID)
	dedupProgress, dedupErr := s.store.GetClaimDedup(ctx, dedupKey)
	if dedupErr == nil {
		item := toProgressItemResponse(definition, dedupProgress, now, cfg.DailyResetHourUTC)
		return dto.ClaimMissionResponse{
			Status:         "idempotent",
			RewardItemID:   definition.RewardItemID,
			RewardQuantity: definition.RewardQuantity,
			Mission:        item,
		}, nil
	}
	if dedupErr != nil && !errors.Is(dedupErr, missionrepository.ErrNotFound) {
		return dto.ClaimMissionResponse{}, dedupErr
	}

	progress, progressErr := s.store.GetMissionProgress(ctx, request.ActorUserID, request.MissionID, periodKey)
	if progressErr != nil {
		if errors.Is(progressErr, missionrepository.ErrNotFound) {
			return dto.ClaimMissionResponse{}, ErrMissionNotEligible
		}
		return dto.ClaimMissionResponse{}, progressErr
	}
	if !progress.Completed {
		return dto.ClaimMissionResponse{}, ErrMissionNotEligible
	}
	if progress.Claimed {
		return dto.ClaimMissionResponse{}, ErrAlreadyClaimed
	}

	claimedAt := now
	progress.Claimed = true
	progress.ClaimedAt = &claimedAt
	progress.LastRequestID = request.RequestID
	progress.LastCorrelationID = request.CorrelationID
	progress.UpdatedAt = now
	if progress.CreatedAt.IsZero() {
		progress.CreatedAt = now
	}

	if err := s.store.UpsertMissionProgress(ctx, progress); err != nil {
		return dto.ClaimMissionResponse{}, err
	}
	if err := s.store.PutClaimDedup(ctx, dedupKey, progress); err != nil {
		return dto.ClaimMissionResponse{}, err
	}

	item := toProgressItemResponse(definition, progress, now, cfg.DailyResetHourUTC)
	return dto.ClaimMissionResponse{
		Status:         "claim_requested",
		RewardItemID:   definition.RewardItemID,
		RewardQuantity: definition.RewardQuantity,
		Mission:        item,
	}, nil
}
