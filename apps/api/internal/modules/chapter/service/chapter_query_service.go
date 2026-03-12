package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/entity"
	chapterrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/repository"
)

func (s *ChapterService) ListChaptersByManga(ctx context.Context, request dto.ListChapterRequest) (dto.ListChapterResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListChapterResponse{}, err
	}

	mangaID, err := parseID(request.MangaID, "manga_id")
	if err != nil {
		return dto.ListChapterResponse{}, err
	}

	sortBy := request.SortBy
	if sortBy == "" {
		sortBy = "sequence_newest"
	}
	limit := request.Limit
	if limit <= 0 {
		limit = 50
	}

	chapters, err := s.store.ListChaptersByManga(ctx, chapterrepository.ListQuery{
		MangaID:            mangaID,
		SortBy:             sortBy,
		Limit:              limit,
		Offset:             request.Offset,
		IncludeUnpublished: request.IncludeUnpublished,
		IncludeDeleted:     false,
	})
	if err != nil {
		return dto.ListChapterResponse{}, err
	}

	items := make([]dto.ChapterListItemResponse, 0, len(chapters))
	for _, chapter := range chapters {
		if chapter.DeletedAt != nil {
			continue
		}
		items = append(items, toListItem(chapter))
	}

	return dto.ListChapterResponse{Items: items, Count: len(items)}, nil
}

func (s *ChapterService) GetChapterDetail(ctx context.Context, request dto.GetChapterDetailRequest) (dto.ChapterDetailResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ChapterDetailResponse{}, err
	}

	chapterID, err := parseID(request.ChapterID, "chapter_id")
	if err != nil {
		return dto.ChapterDetailResponse{}, err
	}

	chapter, err := s.store.GetChapterByID(ctx, chapterID)
	if err != nil {
		if errors.Is(err, chapterrepository.ErrNotFound) {
			return dto.ChapterDetailResponse{}, ErrChapterNotFound
		}
		return dto.ChapterDetailResponse{}, err
	}

	if chapter.DeletedAt != nil || chapter.PublishState != entity.PublishStatePublished {
		return dto.ChapterDetailResponse{}, ErrChapterNotVisible
	}

	return toDetail(chapter), nil
}

func (s *ChapterService) GetNavigation(ctx context.Context, request dto.NavigationRequest) (dto.NavigationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.NavigationResponse{}, err
	}

	chapterID, err := parseID(request.ChapterID, "chapter_id")
	if err != nil {
		return dto.NavigationResponse{}, err
	}

	current, err := s.store.GetChapterByID(ctx, chapterID)
	if err != nil {
		if errors.Is(err, chapterrepository.ErrNotFound) {
			return dto.NavigationResponse{}, ErrChapterNotFound
		}
		return dto.NavigationResponse{}, err
	}

	navigation, err := s.store.ResolveNavigation(ctx, current.MangaID, current.ID)
	if err != nil {
		return dto.NavigationResponse{}, err
	}
	if !navigation.FoundCurrent {
		return dto.NavigationResponse{}, ErrChapterNotVisible
	}

	response := dto.NavigationResponse{CurrentChapterID: current.ID}
	response.FirstChapterID = navigation.FirstID
	response.LastChapterID = navigation.LastID
	response.PreviousChapterID = navigation.PreviousID
	response.NextChapterID = navigation.NextID

	return response, nil
}

func toListItem(chapter entity.Chapter) dto.ChapterListItemResponse {
	return dto.ChapterListItemResponse{
		ChapterID:          chapter.ID,
		MangaID:            chapter.MangaID,
		Slug:               chapter.Slug,
		Title:              chapter.Title,
		SequenceNo:         chapter.SequenceNo,
		DisplayNumber:      chapter.DisplayNumber,
		PublishState:       string(chapter.PublishState),
		ReadAccessLevel:    string(chapter.ReadAccessLevel),
		VIPOnly:            chapter.VIPOnly,
		EarlyAccessEnabled: chapter.EarlyAccessEnabled,
		PreviewEnabled:     chapter.PreviewEnabled,
		PreviewPageCount:   chapter.PreviewPageCount,
		PageCount:          chapter.PageCount,
		PublishedAt:        chapter.PublishedAt,
	}
}

func toDetail(chapter entity.Chapter) dto.ChapterDetailResponse {
	return dto.ChapterDetailResponse{
		ChapterID:                 chapter.ID,
		MangaID:                   chapter.MangaID,
		Slug:                      chapter.Slug,
		Title:                     chapter.Title,
		Summary:                   chapter.Summary,
		SequenceNo:                chapter.SequenceNo,
		DisplayNumber:             chapter.DisplayNumber,
		PublishState:              string(chapter.PublishState),
		ReadAccessLevel:           string(chapter.ReadAccessLevel),
		InheritAccessFromManga:    chapter.InheritAccessFromManga,
		VIPOnly:                   chapter.VIPOnly,
		EarlyAccessEnabled:        chapter.EarlyAccessEnabled,
		EarlyAccessLevel:          string(chapter.EarlyAccessLevel),
		EarlyAccessStartAt:        chapter.EarlyAccessStartAt,
		EarlyAccessEndAt:          chapter.EarlyAccessEndAt,
		EarlyAccessFallbackAccess: string(chapter.EarlyAccessFallbackAccess),
		PreviewEnabled:            chapter.PreviewEnabled,
		PreviewPageCount:          chapter.PreviewPageCount,
		MediaHealthStatus:         string(chapter.MediaHealthStatus),
		IntegrityStatus:           string(chapter.IntegrityStatus),
		PageCount:                 chapter.PageCount,
		ScheduledAt:               chapter.ScheduledAt,
		PublishedAt:               chapter.PublishedAt,
		ArchivedAt:                chapter.ArchivedAt,
		CreatedAt:                 chapter.CreatedAt,
		UpdatedAt:                 chapter.UpdatedAt,
	}
}
