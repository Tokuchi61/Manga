package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/entity"
	chapterrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/repository"
)

func (s *ChapterService) UpdateAccess(ctx context.Context, request dto.UpdateAccessRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	chapterID, err := parseID(request.ChapterID, "chapter_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}

	chapter, err := s.store.GetChapterByID(ctx, chapterID)
	if err != nil {
		if errors.Is(err, chapterrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrChapterNotFound
		}
		return dto.OperationResponse{}, err
	}
	if chapter.DeletedAt != nil {
		return dto.OperationResponse{}, ErrInvalidStateTransition
	}

	if request.ReadAccessLevel != nil {
		value, parseErr := toReadAccessLevel(*request.ReadAccessLevel, chapter.ReadAccessLevel)
		if parseErr != nil {
			return dto.OperationResponse{}, parseErr
		}
		chapter.ReadAccessLevel = value
	}
	if request.InheritAccessFromManga != nil {
		chapter.InheritAccessFromManga = *request.InheritAccessFromManga
	}
	if request.VIPOnly != nil {
		chapter.VIPOnly = *request.VIPOnly
	}
	if request.EarlyAccessEnabled != nil {
		chapter.EarlyAccessEnabled = *request.EarlyAccessEnabled
	}
	if request.EarlyAccessLevel != nil {
		value, parseErr := toEarlyAccessLevel(*request.EarlyAccessLevel, chapter.EarlyAccessEnabled, chapter.VIPOnly)
		if parseErr != nil {
			return dto.OperationResponse{}, parseErr
		}
		chapter.EarlyAccessLevel = value
	}
	if request.EarlyAccessFallbackAccess != nil {
		value, parseErr := toReadAccessLevel(*request.EarlyAccessFallbackAccess, chapter.EarlyAccessFallbackAccess)
		if parseErr != nil {
			return dto.OperationResponse{}, parseErr
		}
		chapter.EarlyAccessFallbackAccess = value
	}
	if request.EarlyAccessStartAt != nil {
		value, parseErr := parseRFC3339Ptr(request.EarlyAccessStartAt, "early_access_start_at")
		if parseErr != nil {
			return dto.OperationResponse{}, parseErr
		}
		chapter.EarlyAccessStartAt = value
	}
	if request.EarlyAccessEndAt != nil {
		value, parseErr := parseRFC3339Ptr(request.EarlyAccessEndAt, "early_access_end_at")
		if parseErr != nil {
			return dto.OperationResponse{}, parseErr
		}
		chapter.EarlyAccessEndAt = value
	}
	if request.PreviewEnabled != nil {
		chapter.PreviewEnabled = *request.PreviewEnabled
	}
	if request.PreviewPageCount != nil {
		chapter.PreviewPageCount = pagePreviewCount(chapter.PageCount, *request.PreviewPageCount, chapter.PreviewEnabled)
	}

	if chapter.VIPOnly {
		chapter.ReadAccessLevel = entity.ReadAccessVIP
		chapter.EarlyAccessEnabled = false
		chapter.EarlyAccessLevel = entity.EarlyAccessNone
		chapter.EarlyAccessFallbackAccess = entity.ReadAccessVIP
		chapter.EarlyAccessStartAt = nil
		chapter.EarlyAccessEndAt = nil
	}
	if chapter.EarlyAccessEnabled {
		if chapter.EarlyAccessEndAt != nil && chapter.EarlyAccessStartAt != nil && chapter.EarlyAccessEndAt.Before(*chapter.EarlyAccessStartAt) {
			return dto.OperationResponse{}, fmt.Errorf("%w: early_access_end_at must be after early_access_start_at", ErrValidation)
		}
	}
	if !chapter.EarlyAccessEnabled {
		chapter.EarlyAccessLevel = entity.EarlyAccessNone
		chapter.EarlyAccessStartAt = nil
		chapter.EarlyAccessEndAt = nil
	}

	chapter.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateChapter(ctx, chapter); err != nil {
		if errors.Is(err, chapterrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrChapterNotFound
		}
		if errors.Is(err, chapterrepository.ErrConflict) {
			return dto.OperationResponse{}, ErrChapterAlreadyExists
		}
		return dto.OperationResponse{}, err
	}

	return dto.OperationResponse{Status: "access_updated"}, nil
}
