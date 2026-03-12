package service

import (
	"context"
	"errors"

	chaptercontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/contract"
	chapterrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/repository"
)

// GetResumeAnchor exposes chapter-owned anchor surface for history module.
func (s *ChapterService) GetResumeAnchor(ctx context.Context, chapterID string, pageNumber int) (chaptercontract.ResumeAnchor, error) {
	parsedID, err := parseID(chapterID, "chapter_id")
	if err != nil {
		return chaptercontract.ResumeAnchor{}, err
	}

	chapter, err := s.store.GetChapterByID(ctx, parsedID)
	if err != nil {
		if errors.Is(err, chapterrepository.ErrNotFound) {
			return chaptercontract.ResumeAnchor{}, ErrChapterNotFound
		}
		return chaptercontract.ResumeAnchor{}, err
	}

	if pageNumber <= 0 {
		pageNumber = 1
	}
	if pageNumber > chapter.PageCount {
		pageNumber = chapter.PageCount
	}

	return chaptercontract.ResumeAnchor{
		ChapterID:  chapter.ID,
		MangaID:    chapter.MangaID,
		PageNumber: pageNumber,
		PageCount:  chapter.PageCount,
		UpdatedAt:  chapter.UpdatedAt,
	}, nil
}

// TargetExists exposes chapter target existence checks for consumer modules.
func (s *ChapterService) TargetExists(ctx context.Context, chapterID string) (bool, error) {
	parsedID, err := parseID(chapterID, "chapter_id")
	if err != nil {
		return false, nil
	}

	_, err = s.store.GetChapterByID(ctx, parsedID)
	if err != nil {
		if errors.Is(err, chapterrepository.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// BuildReadSignal creates stable chapter->history signal payload.
func (s *ChapterService) BuildReadSignal(chapterID string, mangaID string, pageNumber int, pageCount int, event string, requestID string, correlationID string) chaptercontract.ReadSignal {
	if pageNumber < 0 {
		pageNumber = 0
	}
	if pageCount < 0 {
		pageCount = 0
	}
	if event == "" {
		event = chaptercontract.EventReadCheckpoint
	}

	return chaptercontract.ReadSignal{
		Event:         event,
		ChapterID:     chapterID,
		MangaID:       mangaID,
		PageNumber:    pageNumber,
		PageCount:     pageCount,
		OccurredAt:    s.now().UTC(),
		RequestID:     requestID,
		CorrelationID: correlationID,
	}
}
