package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/dto"
	rprepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/repository"
)

func (s *RoyalPassService) IngestRoyalPassProgress(ctx context.Context, request dto.IngestRoyalPassProgressRequest) (dto.IngestRoyalPassProgressResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.IngestRoyalPassProgressResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.IngestRoyalPassProgressResponse{}, err
	}
	if err := s.requireSeasonEnabled(cfg.SeasonEnabled); err != nil {
		return dto.IngestRoyalPassProgressResponse{}, err
	}

	now := s.now().UTC()
	season, err := s.resolveSeason(ctx, request.SeasonID, now)
	if err != nil {
		return dto.IngestRoyalPassProgressResponse{}, err
	}
	if !isSeasonWritable(season, now) {
		return dto.IngestRoyalPassProgressResponse{}, ErrSeasonUnavailable
	}

	if strings.TrimSpace(request.RequestID) != "" {
		dedupKey := buildProgressDedupKey(request.ActorUserID, season.SeasonID, request.RequestID)
		dedupProgress, dedupErr := s.store.GetProgressDedup(ctx, dedupKey)
		if dedupErr == nil {
			return dto.IngestRoyalPassProgressResponse{Status: "idempotent", SeasonID: season.SeasonID, Points: dedupProgress.Points}, nil
		}
		if dedupErr != nil && !errors.Is(dedupErr, rprepository.ErrNotFound) {
			return dto.IngestRoyalPassProgressResponse{}, dedupErr
		}
	}

	progress, progressErr := s.store.GetUserProgress(ctx, request.ActorUserID, season.SeasonID)
	if progressErr != nil {
		if !errors.Is(progressErr, rprepository.ErrNotFound) {
			return dto.IngestRoyalPassProgressResponse{}, progressErr
		}
		progress = fallbackUserProgress(request.ActorUserID, season.SeasonID, now)
	}

	progress.Points += request.Delta
	progress.LastRequestID = request.RequestID
	progress.LastCorrelationID = request.CorrelationID
	progress.UpdatedAt = now
	if progress.CreatedAt.IsZero() {
		progress.CreatedAt = now
	}

	if err := s.store.UpsertUserProgress(ctx, progress); err != nil {
		return dto.IngestRoyalPassProgressResponse{}, err
	}

	if strings.TrimSpace(request.RequestID) != "" {
		dedupKey := buildProgressDedupKey(request.ActorUserID, season.SeasonID, request.RequestID)
		if err := s.store.PutProgressDedup(ctx, dedupKey, progress); err != nil {
			return dto.IngestRoyalPassProgressResponse{}, err
		}
	}

	return dto.IngestRoyalPassProgressResponse{Status: "progressed", SeasonID: season.SeasonID, Points: progress.Points}, nil
}
