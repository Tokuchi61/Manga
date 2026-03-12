package service

import (
	"context"
	"errors"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/entity"
	chapterrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/repository"
)

func (s *ChapterService) ReadChapter(ctx context.Context, request dto.ReadChapterRequest) (dto.ReadChapterResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ReadChapterResponse{}, err
	}

	chapterID, err := parseID(request.ChapterID, "chapter_id")
	if err != nil {
		return dto.ReadChapterResponse{}, err
	}

	chapter, err := s.store.GetChapterByID(ctx, chapterID)
	if err != nil {
		if errors.Is(err, chapterrepository.ErrNotFound) {
			return dto.ReadChapterResponse{}, ErrChapterNotFound
		}
		return dto.ReadChapterResponse{}, err
	}
	if chapter.DeletedAt != nil || chapter.PublishState != entity.PublishStatePublished {
		return dto.ReadChapterResponse{}, ErrChapterNotVisible
	}

	mode := normalizeValue(request.Mode)
	if mode == "" {
		mode = "full"
	}

	pages := chapter.Pages
	if mode == "preview" {
		if !chapter.PreviewEnabled {
			return dto.ReadChapterResponse{}, ErrChapterNotVisible
		}
		previewCount := pagePreviewCount(chapter.PageCount, chapter.PreviewPageCount, chapter.PreviewEnabled)
		if previewCount < len(pages) {
			pages = pages[:previewCount]
		}
	}

	readAt := s.now().UTC()
	if request.At != nil {
		readAt = request.At.UTC()
	}
	earlyAccessActive := isEarlyAccessActive(chapter, readAt)

	responsePages := make([]dto.ChapterPageResponse, 0, len(pages))
	for _, page := range pages {
		responsePages = append(responsePages, dto.ChapterPageResponse{
			PageNumber: page.PageNumber,
			MediaURL:   page.MediaURL,
			Width:      page.Width,
			Height:     page.Height,
			LongStrip:  page.LongStrip,
		})
	}

	return dto.ReadChapterResponse{
		ChapterID:                 chapter.ID,
		Mode:                      mode,
		PublishState:              string(chapter.PublishState),
		ReadAccessLevel:           string(chapter.ReadAccessLevel),
		VIPOnly:                   chapter.VIPOnly,
		EarlyAccessEnabled:        chapter.EarlyAccessEnabled,
		EarlyAccessLevel:          string(chapter.EarlyAccessLevel),
		EarlyAccessActive:         earlyAccessActive,
		EarlyAccessFallbackAccess: string(chapter.EarlyAccessFallbackAccess),
		Pages:                     responsePages,
		PageCount:                 len(responsePages),
	}, nil
}

func isEarlyAccessActive(chapter entity.Chapter, at time.Time) bool {
	if !chapter.EarlyAccessEnabled {
		return false
	}
	if chapter.EarlyAccessStartAt != nil && at.Before(chapter.EarlyAccessStartAt.UTC()) {
		return false
	}
	if chapter.EarlyAccessEndAt != nil && !at.Before(chapter.EarlyAccessEndAt.UTC()) {
		return false
	}
	return true
}
