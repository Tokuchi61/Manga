package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/dto"
	chapterrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/repository"
)

func (s *ChapterService) UpdateMediaHealth(ctx context.Context, request dto.UpdateMediaHealthRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	chapterID, err := parseID(request.ChapterID, "chapter_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	status, err := toMediaHealthStatus(request.MediaHealthStatus)
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

	chapter.MediaHealthStatus = status
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
	return dto.OperationResponse{Status: "media_health_updated"}, nil
}

func (s *ChapterService) UpdateIntegrity(ctx context.Context, request dto.UpdateIntegrityRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	chapterID, err := parseID(request.ChapterID, "chapter_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	status, err := toIntegrityStatus(request.IntegrityStatus)
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

	chapter.IntegrityStatus = status
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
	return dto.OperationResponse{Status: "integrity_updated"}, nil
}
