package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/entity"
	chapterrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/repository"
)

func (s *ChapterService) SoftDeleteChapter(ctx context.Context, request dto.SoftDeleteChapterRequest) (dto.OperationResponse, error) {
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
		return dto.OperationResponse{Status: "already_deleted"}, nil
	}

	now := s.now().UTC()
	chapter.DeletedAt = &now
	chapter.PublishState = entity.PublishStateArchived
	chapter.ArchivedAt = &now
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

	return dto.OperationResponse{Status: "soft_deleted"}, nil
}

func (s *ChapterService) RestoreChapter(ctx context.Context, request dto.RestoreChapterRequest) (dto.OperationResponse, error) {
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
	if chapter.DeletedAt == nil {
		return dto.OperationResponse{Status: "already_active"}, nil
	}

	chapter.DeletedAt = nil
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

	return dto.OperationResponse{Status: "restored"}, nil
}
