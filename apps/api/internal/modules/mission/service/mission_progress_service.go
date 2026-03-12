package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/dto"
	missionrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/repository"
)

func (s *MissionService) IngestMissionProgress(ctx context.Context, request dto.IngestMissionProgressRequest) (dto.IngestMissionProgressResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.IngestMissionProgressResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.IngestMissionProgressResponse{}, err
	}
	if err := s.requireProgressIngestEnabled(cfg.ProgressIngestEnabled); err != nil {
		return dto.IngestMissionProgressResponse{}, err
	}

	definition, err := s.store.GetMissionDefinition(ctx, request.MissionID)
	if err != nil {
		if errors.Is(err, missionrepository.ErrNotFound) {
			return dto.IngestMissionProgressResponse{}, ErrNotFound
		}
		return dto.IngestMissionProgressResponse{}, err
	}

	now := s.now().UTC()
	if !isMissionActive(definition, now) {
		return dto.IngestMissionProgressResponse{}, ErrMissionInactive
	}

	periodKey := buildPeriodKey(definition.Category, now, cfg.DailyResetHourUTC)
	if strings.TrimSpace(request.RequestID) != "" {
		dedupKey := buildProgressDedupKey(request.ActorUserID, request.MissionID, periodKey, request.RequestID)
		dedupProgress, dedupErr := s.store.GetProgressDedup(ctx, dedupKey)
		if dedupErr == nil {
			item := toProgressItemResponse(definition, dedupProgress, now, cfg.DailyResetHourUTC)
			return dto.IngestMissionProgressResponse{Status: "idempotent", Completed: dedupProgress.Completed, Mission: item}, nil
		}
		if dedupErr != nil && !errors.Is(dedupErr, missionrepository.ErrNotFound) {
			return dto.IngestMissionProgressResponse{}, dedupErr
		}
	}

	progress, progressErr := s.store.GetMissionProgress(ctx, request.ActorUserID, request.MissionID, periodKey)
	if progressErr != nil && !errors.Is(progressErr, missionrepository.ErrNotFound) {
		return dto.IngestMissionProgressResponse{}, progressErr
	}
	if errors.Is(progressErr, missionrepository.ErrNotFound) {
		progress = fallbackProgress(definition, request.ActorUserID, periodKey, now)
	}
	if progress.Claimed {
		return dto.IngestMissionProgressResponse{}, ErrAlreadyClaimed
	}

	wasCompleted := progress.Completed
	progress.ProgressCount += request.Delta
	if progress.ProgressCount >= definition.TargetCount {
		progress.Completed = true
		if progress.CompletedAt == nil {
			completedAt := now
			progress.CompletedAt = &completedAt
		}
	}
	progress.LastRequestID = request.RequestID
	progress.LastCorrelationID = request.CorrelationID
	progress.UpdatedAt = now
	if progress.CreatedAt.IsZero() {
		progress.CreatedAt = now
	}

	if err := s.store.UpsertMissionProgress(ctx, progress); err != nil {
		return dto.IngestMissionProgressResponse{}, err
	}

	if strings.TrimSpace(request.RequestID) != "" {
		dedupKey := buildProgressDedupKey(request.ActorUserID, request.MissionID, periodKey, request.RequestID)
		if err := s.store.PutProgressDedup(ctx, dedupKey, progress); err != nil {
			return dto.IngestMissionProgressResponse{}, err
		}
	}

	status := "progressed"
	if progress.Completed {
		status = "completed"
		if wasCompleted {
			status = "progressed"
		}
	}

	item := toProgressItemResponse(definition, progress, now, cfg.DailyResetHourUTC)
	return dto.IngestMissionProgressResponse{Status: status, Completed: progress.Completed, Mission: item}, nil
}
