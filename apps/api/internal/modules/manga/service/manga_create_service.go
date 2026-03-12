package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/entity"
	mangarepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/repository"
	"github.com/google/uuid"
)

func (s *MangaService) CreateManga(ctx context.Context, request dto.CreateMangaRequest) (dto.CreateMangaResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CreateMangaResponse{}, err
	}

	slug, err := ensureSlug(request.Slug, request.Title)
	if err != nil {
		return dto.CreateMangaResponse{}, err
	}
	readAccessLevel, err := toReadAccessLevel(request.DefaultReadAccessLevel)
	if err != nil {
		return dto.CreateMangaResponse{}, err
	}
	earlyAccessLevel, err := toEarlyAccessLevel(request.DefaultEarlyAccessLevel, request.DefaultEarlyAccessEnabled)
	if err != nil {
		return dto.CreateMangaResponse{}, err
	}
	publishState, err := toPublishState(request.PublishState, entity.PublishStateDraft)
	if err != nil {
		return dto.CreateMangaResponse{}, err
	}
	scheduledAt, err := parseScheduledAt(request.ScheduledAt)
	if err != nil {
		return dto.CreateMangaResponse{}, err
	}

	now := s.now().UTC()
	var publishedAt *time.Time
	var archivedAt *time.Time

	switch publishState {
	case entity.PublishStateScheduled:
		if scheduledAt == nil {
			return dto.CreateMangaResponse{}, fmt.Errorf("%w: scheduled publish_state requires scheduled_at", ErrValidation)
		}
	case entity.PublishStatePublished:
		publishedAt = &now
		scheduledAt = nil
	case entity.PublishStateArchived:
		archivedAt = &now
		scheduledAt = nil
	}

	manga := entity.Manga{
		ID:                 uuid.NewString(),
		Slug:               slug,
		Title:              normalizeFreeText(request.Title),
		AlternativeTitles:  normalizeList(request.AlternativeTitles),
		Summary:            normalizeFreeText(request.Summary),
		ShortSummary:       normalizeFreeText(request.ShortSummary),
		CoverImageURL:      normalizeFreeText(request.CoverImageURL),
		BannerImageURL:     normalizeFreeText(request.BannerImageURL),
		SEOTitle:           normalizeFreeText(request.SEOTitle),
		SEODescription:     normalizeFreeText(request.SEODescription),
		Genres:             normalizeList(request.Genres),
		Tags:               normalizeList(request.Tags),
		Themes:             normalizeList(request.Themes),
		ContentWarnings:    normalizeList(request.ContentWarnings),
		PublishState:       publishState,
		Visibility:         entity.VisibilityPublic,
		DefaultReadAccess:  readAccessLevel,
		EarlyAccessEnabled: request.DefaultEarlyAccessEnabled,
		EarlyAccessLevel:   earlyAccessLevel,
		ReleaseSchedule:    normalizeFreeText(request.ReleaseSchedule),
		TranslationGroup:   normalizeFreeText(request.TranslationGroup),
		ContentVersion:     1,
		ScheduledAt:        scheduledAt,
		PublishedAt:        publishedAt,
		ArchivedAt:         archivedAt,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if manga.ShortSummary == "" {
		manga.ShortSummary = manga.Summary
	}

	if err := s.store.CreateManga(ctx, manga); err != nil {
		if errors.Is(err, mangarepository.ErrConflict) {
			return dto.CreateMangaResponse{}, ErrMangaAlreadyExists
		}
		return dto.CreateMangaResponse{}, err
	}

	return dto.CreateMangaResponse{
		MangaID:      manga.ID,
		Slug:         manga.Slug,
		PublishState: string(manga.PublishState),
		Visibility:   string(manga.Visibility),
	}, nil
}

func parseScheduledAt(value *time.Time) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}
	resolved := value.UTC()
	return &resolved, nil
}
