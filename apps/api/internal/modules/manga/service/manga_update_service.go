package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/dto"
	mangarepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/repository"
)

func (s *MangaService) UpdateMangaMetadata(ctx context.Context, request dto.UpdateMangaMetadataRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	mangaID, err := parseID(request.MangaID, "manga_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}

	manga, err := s.store.GetMangaByID(ctx, mangaID)
	if err != nil {
		if errors.Is(err, mangarepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrMangaNotFound
		}
		return dto.OperationResponse{}, err
	}

	if manga.DeletedAt != nil {
		return dto.OperationResponse{}, ErrInvalidStateTransition
	}

	if request.Title != nil {
		manga.Title = normalizeFreeText(*request.Title)
	}
	if request.Slug != nil {
		slug, slugErr := ensureSlug(*request.Slug, manga.Title)
		if slugErr != nil {
			return dto.OperationResponse{}, slugErr
		}
		manga.Slug = slug
	}
	if request.AlternativeTitles != nil {
		manga.AlternativeTitles = normalizeList(*request.AlternativeTitles)
	}
	if request.Summary != nil {
		manga.Summary = normalizeFreeText(*request.Summary)
	}
	if request.ShortSummary != nil {
		manga.ShortSummary = normalizeFreeText(*request.ShortSummary)
	}
	if request.CoverImageURL != nil {
		manga.CoverImageURL = normalizeFreeText(*request.CoverImageURL)
	}
	if request.BannerImageURL != nil {
		manga.BannerImageURL = normalizeFreeText(*request.BannerImageURL)
	}
	if request.SEOTitle != nil {
		manga.SEOTitle = normalizeFreeText(*request.SEOTitle)
	}
	if request.SEODescription != nil {
		manga.SEODescription = normalizeFreeText(*request.SEODescription)
	}
	if request.Genres != nil {
		manga.Genres = normalizeList(*request.Genres)
	}
	if request.Tags != nil {
		manga.Tags = normalizeList(*request.Tags)
	}
	if request.Themes != nil {
		manga.Themes = normalizeList(*request.Themes)
	}
	if request.ContentWarnings != nil {
		manga.ContentWarnings = normalizeList(*request.ContentWarnings)
	}
	if request.ReleaseSchedule != nil {
		manga.ReleaseSchedule = normalizeFreeText(*request.ReleaseSchedule)
	}
	if request.TranslationGroup != nil {
		manga.TranslationGroup = normalizeFreeText(*request.TranslationGroup)
	}
	if manga.ShortSummary == "" {
		manga.ShortSummary = manga.Summary
	}

	manga.ContentVersion++
	manga.UpdatedAt = s.now().UTC()

	if err := s.store.UpdateManga(ctx, manga); err != nil {
		if errors.Is(err, mangarepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrMangaNotFound
		}
		if errors.Is(err, mangarepository.ErrConflict) {
			return dto.OperationResponse{}, ErrMangaAlreadyExists
		}
		return dto.OperationResponse{}, err
	}

	return dto.OperationResponse{Status: "metadata_updated"}, nil
}
