package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/dto"
	chapterrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/repository"
)

func (s *ChapterService) UpdateChapter(ctx context.Context, request dto.UpdateChapterRequest) (dto.OperationResponse, error) {
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

	if request.Slug != nil {
		slug, slugErr := ensureSlug(*request.Slug, chapter.Title, chapter.DisplayNumber)
		if slugErr != nil {
			return dto.OperationResponse{}, slugErr
		}
		chapter.Slug = slug
	}
	if request.Title != nil {
		chapter.Title = *request.Title
	}
	if request.Summary != nil {
		chapter.Summary = *request.Summary
	}
	if request.SequenceNo != nil {
		chapter.SequenceNo = *request.SequenceNo
	}
	if request.DisplayNumber != nil {
		chapter.DisplayNumber = *request.DisplayNumber
	}
	if request.Pages != nil {
		pages, pageErr := toPages(*request.Pages, s.now().UTC())
		if pageErr != nil {
			return dto.OperationResponse{}, pageErr
		}
		chapter.Pages = pages
		chapter.PageCount = len(pages)
		if chapter.PreviewPageCount > chapter.PageCount {
			chapter.PreviewPageCount = chapter.PageCount
		}
	}
	if chapter.DisplayNumber == "" {
		chapter.DisplayNumber = fmt.Sprintf("%d", chapter.SequenceNo)
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

	return dto.OperationResponse{Status: "chapter_updated"}, nil
}

func (s *ChapterService) ReorderChapter(ctx context.Context, request dto.ReorderChapterRequest) (dto.OperationResponse, error) {
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

	chapter.SequenceNo = request.SequenceNo
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

	return dto.OperationResponse{Status: "chapter_reordered"}, nil
}
