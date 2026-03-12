package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/entity"
	mangarepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/repository"
)

func (s *MangaService) ListPublicManga(ctx context.Context, request dto.ListMangaRequest) (dto.ListMangaResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListMangaResponse{}, err
	}

	sortBy := request.SortBy
	if sortBy == "" {
		sortBy = "newest"
	}
	limit := request.Limit
	if limit <= 0 {
		limit = 20
	}

	items, err := s.store.ListManga(ctx, mangarepository.ListQuery{
		Search:             request.Search,
		Genre:              request.Genre,
		Tag:                request.Tag,
		Theme:              request.Theme,
		ContentWarning:     request.ContentWarning,
		SortBy:             sortBy,
		Limit:              limit,
		Offset:             request.Offset,
		IncludeUnpublished: false,
		IncludeHidden:      false,
		IncludeDeleted:     false,
	})
	if err != nil {
		return dto.ListMangaResponse{}, err
	}

	responseItems := make([]dto.MangaListItemResponse, 0, len(items))
	for _, item := range items {
		responseItems = append(responseItems, toListItem(item))
	}

	return dto.ListMangaResponse{Items: responseItems, Count: len(responseItems)}, nil
}

func (s *MangaService) GetPublicMangaDetail(ctx context.Context, request dto.GetMangaDetailRequest) (dto.MangaDetailResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.MangaDetailResponse{}, err
	}

	mangaID, err := parseID(request.MangaID, "manga_id")
	if err != nil {
		return dto.MangaDetailResponse{}, err
	}

	manga, err := s.store.GetMangaByID(ctx, mangaID)
	if err != nil {
		if errors.Is(err, mangarepository.ErrNotFound) {
			return dto.MangaDetailResponse{}, ErrMangaNotFound
		}
		return dto.MangaDetailResponse{}, err
	}

	if manga.DeletedAt != nil || manga.PublishState != entity.PublishStatePublished || manga.Visibility != entity.VisibilityPublic {
		return dto.MangaDetailResponse{}, ErrMangaNotVisible
	}

	return toDetail(manga), nil
}

func toListItem(manga entity.Manga) dto.MangaListItemResponse {
	return dto.MangaListItemResponse{
		MangaID:          manga.ID,
		Slug:             manga.Slug,
		Title:            manga.Title,
		ShortSummary:     manga.ShortSummary,
		CoverImageURL:    manga.CoverImageURL,
		Genres:           append([]string(nil), manga.Genres...),
		Tags:             append([]string(nil), manga.Tags...),
		PublishState:     string(manga.PublishState),
		Visibility:       string(manga.Visibility),
		Featured:         manga.Featured,
		Recommended:      manga.Recommended,
		ChapterCount:     manga.ChapterCount,
		CommentCount:     manga.CommentCount,
		ViewCount:        manga.ViewCount,
		ReleaseSchedule:  manga.ReleaseSchedule,
		TranslationGroup: manga.TranslationGroup,
		PublishedAt:      manga.PublishedAt,
	}
}

func toDetail(manga entity.Manga) dto.MangaDetailResponse {
	return dto.MangaDetailResponse{
		MangaID:            manga.ID,
		Slug:               manga.Slug,
		Title:              manga.Title,
		AlternativeTitles:  append([]string(nil), manga.AlternativeTitles...),
		Summary:            manga.Summary,
		ShortSummary:       manga.ShortSummary,
		CoverImageURL:      manga.CoverImageURL,
		BannerImageURL:     manga.BannerImageURL,
		SEOTitle:           manga.SEOTitle,
		SEODescription:     manga.SEODescription,
		Genres:             append([]string(nil), manga.Genres...),
		Tags:               append([]string(nil), manga.Tags...),
		Themes:             append([]string(nil), manga.Themes...),
		ContentWarnings:    append([]string(nil), manga.ContentWarnings...),
		PublishState:       string(manga.PublishState),
		Visibility:         string(manga.Visibility),
		Featured:           manga.Featured,
		Recommended:        manga.Recommended,
		CollectionKeys:     append([]string(nil), manga.CollectionKeys...),
		DefaultReadAccess:  string(manga.DefaultReadAccess),
		EarlyAccessEnabled: manga.EarlyAccessEnabled,
		EarlyAccessLevel:   string(manga.EarlyAccessLevel),
		ReleaseSchedule:    manga.ReleaseSchedule,
		TranslationGroup:   manga.TranslationGroup,
		ChapterCount:       manga.ChapterCount,
		CommentCount:       manga.CommentCount,
		ViewCount:          manga.ViewCount,
		ContentVersion:     manga.ContentVersion,
		ScheduledAt:        manga.ScheduledAt,
		PublishedAt:        manga.PublishedAt,
		ArchivedAt:         manga.ArchivedAt,
		CreatedAt:          manga.CreatedAt,
		UpdatedAt:          manga.UpdatedAt,
	}
}
