package service

import (
	"context"
	"errors"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/entity"
	missionrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/repository"
)

func (s *MissionService) ListActorMissions(ctx context.Context, request dto.ListActorMissionsRequest) (dto.ListActorMissionsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListActorMissionsResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ListActorMissionsResponse{}, err
	}
	if err := s.requireReadEnabled(cfg.ReadEnabled); err != nil {
		return dto.ListActorMissionsResponse{}, err
	}

	definitions, err := s.store.ListMissionDefinitions(ctx, request.Category, false, 0, 0)
	if err != nil {
		return dto.ListActorMissionsResponse{}, err
	}

	now := s.now().UTC()
	filteredItems := make([]dto.MissionProgressItemResponse, 0, len(definitions))
	for _, definition := range definitions {
		periodKey := buildPeriodKey(definition.Category, now, cfg.DailyResetHourUTC)
		progress, progressErr := s.store.GetMissionProgress(ctx, request.ActorUserID, definition.MissionID, periodKey)
		if progressErr != nil && !errors.Is(progressErr, missionrepository.ErrNotFound) {
			return dto.ListActorMissionsResponse{}, progressErr
		}

		if errors.Is(progressErr, missionrepository.ErrNotFound) {
			progress = fallbackProgress(definition, request.ActorUserID, periodKey, now)
		}
		item := toProgressItemResponse(definition, progress, now, cfg.DailyResetHourUTC)

		if request.State != "" && normalizeValue(item.Status) != normalizeValue(request.State) {
			continue
		}
		filteredItems = append(filteredItems, item)
	}

	items := applyOffsetLimit(filteredItems, request.Offset, request.Limit)
	return dto.ListActorMissionsResponse{Items: items, Count: len(items)}, nil
}

func (s *MissionService) GetActorMissionDetail(ctx context.Context, request dto.GetActorMissionDetailRequest) (dto.MissionDetailResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.MissionDetailResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.MissionDetailResponse{}, err
	}
	if err := s.requireReadEnabled(cfg.ReadEnabled); err != nil {
		return dto.MissionDetailResponse{}, err
	}

	definition, err := s.store.GetMissionDefinition(ctx, request.MissionID)
	if err != nil {
		if errors.Is(err, missionrepository.ErrNotFound) {
			return dto.MissionDetailResponse{}, ErrNotFound
		}
		return dto.MissionDetailResponse{}, err
	}

	now := s.now().UTC()
	periodKey := buildPeriodKey(definition.Category, now, cfg.DailyResetHourUTC)
	progress, progressErr := s.store.GetMissionProgress(ctx, request.ActorUserID, definition.MissionID, periodKey)
	if progressErr != nil && !errors.Is(progressErr, missionrepository.ErrNotFound) {
		return dto.MissionDetailResponse{}, progressErr
	}
	if errors.Is(progressErr, missionrepository.ErrNotFound) {
		progress = fallbackProgress(definition, request.ActorUserID, periodKey, now)
	}

	return toMissionDetailResponse(definition, progress, now, cfg.DailyResetHourUTC), nil
}

func fallbackProgress(definition entity.MissionDefinition, actorUserID string, periodKey string, now time.Time) entity.MissionProgress {
	progress := entity.MissionProgress{
		UserID:        actorUserID,
		MissionID:     definition.MissionID,
		PeriodKey:     periodKey,
		ProgressCount: 0,
		CreatedAt:     definition.CreatedAt,
		UpdatedAt:     definition.UpdatedAt,
	}
	if progress.UpdatedAt.IsZero() {
		progress.UpdatedAt = now.UTC()
	}
	if progress.CreatedAt.IsZero() {
		progress.CreatedAt = progress.UpdatedAt
	}
	return progress
}
