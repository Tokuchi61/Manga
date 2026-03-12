package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/entity"
	chapterrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/repository"
)

func (s *ChapterService) UpdatePublishState(ctx context.Context, request dto.UpdatePublishStateRequest) (dto.OperationResponse, error) {
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

	now := s.now().UTC()
	action := normalizeValue(request.Action)
	scheduledAt, err := parseRFC3339Ptr(request.ScheduledAt, "scheduled_at")
	if err != nil {
		return dto.OperationResponse{}, err
	}

	switch action {
	case "draft":
		chapter.PublishState = entity.PublishStateDraft
		chapter.ScheduledAt = nil
		chapter.ArchivedAt = nil
	case "schedule":
		if scheduledAt == nil {
			return dto.OperationResponse{}, fmt.Errorf("%w: schedule action requires scheduled_at", ErrValidation)
		}
		chapter.PublishState = entity.PublishStateScheduled
		chapter.ScheduledAt = scheduledAt
		chapter.ArchivedAt = nil
	case "publish":
		chapter.PublishState = entity.PublishStatePublished
		chapter.PublishedAt = &now
		chapter.ScheduledAt = nil
		chapter.ArchivedAt = nil
	case "archive":
		chapter.PublishState = entity.PublishStateArchived
		chapter.ArchivedAt = &now
		chapter.ScheduledAt = nil
	case "unpublish":
		chapter.PublishState = entity.PublishStateUnpublished
		chapter.ScheduledAt = nil
		chapter.ArchivedAt = nil
	default:
		return dto.OperationResponse{}, fmt.Errorf("%w: unsupported publish action", ErrValidation)
	}

	chapter.UpdatedAt = now
	if err := s.store.UpdateChapter(ctx, chapter); err != nil {
		if errors.Is(err, chapterrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrChapterNotFound
		}
		if errors.Is(err, chapterrepository.ErrConflict) {
			return dto.OperationResponse{}, ErrChapterAlreadyExists
		}
		return dto.OperationResponse{}, err
	}

	return dto.OperationResponse{Status: "publish_state_updated"}, nil
}
