package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/entity"
	chapterrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/repository"
	"github.com/google/uuid"
)

func (s *ChapterService) CreateChapter(ctx context.Context, request dto.CreateChapterRequest) (dto.CreateChapterResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CreateChapterResponse{}, err
	}

	mangaID, err := parseID(request.MangaID, "manga_id")
	if err != nil {
		return dto.CreateChapterResponse{}, err
	}
	if s.mangaLookup != nil {
		exists, lookupErr := s.mangaLookup.TargetExists(ctx, mangaID)
		if lookupErr != nil {
			return dto.CreateChapterResponse{}, lookupErr
		}
		if !exists {
			return dto.CreateChapterResponse{}, ErrMangaNotFound
		}
	}
	slug, err := ensureSlug(request.Slug, request.Title, request.DisplayNumber)
	if err != nil {
		return dto.CreateChapterResponse{}, err
	}

	pages, err := toPages(request.Pages, s.now().UTC())
	if err != nil {
		return dto.CreateChapterResponse{}, err
	}

	publishState, err := toPublishState(request.PublishState, entity.PublishStateDraft)
	if err != nil {
		return dto.CreateChapterResponse{}, err
	}
	readAccessLevel, err := toReadAccessLevel(request.ReadAccessLevel, entity.ReadAccessAuthenticated)
	if err != nil {
		return dto.CreateChapterResponse{}, err
	}
	fallbackAccess, err := toReadAccessLevel(request.EarlyAccessFallbackAccess, entity.ReadAccessAuthenticated)
	if err != nil {
		return dto.CreateChapterResponse{}, err
	}
	earlyAccessLevel, err := toEarlyAccessLevel(request.EarlyAccessLevel, request.EarlyAccessEnabled, request.VIPOnly)
	if err != nil {
		return dto.CreateChapterResponse{}, err
	}

	if request.VIPOnly {
		readAccessLevel = entity.ReadAccessVIP
		request.EarlyAccessEnabled = false
		earlyAccessLevel = entity.EarlyAccessNone
		fallbackAccess = entity.ReadAccessVIP
	}

	earlyStart := request.EarlyAccessStartAt
	earlyEnd := request.EarlyAccessEndAt
	if request.EarlyAccessEnabled {
		if earlyStart != nil && earlyEnd != nil && earlyEnd.Before(*earlyStart) {
			return dto.CreateChapterResponse{}, fmt.Errorf("%w: early_access_end_at must be after early_access_start_at", ErrValidation)
		}
	}
	if !request.EarlyAccessEnabled {
		earlyStart = nil
		earlyEnd = nil
	}

	now := s.now().UTC()
	scheduledAt, publishedAt, archivedAt, err := createPublishTimes(publishState, request.ScheduledAt, now)
	if err != nil {
		return dto.CreateChapterResponse{}, err
	}

	chapter := entity.Chapter{
		ID:                        uuid.NewString(),
		MangaID:                   mangaID,
		Slug:                      slug,
		Title:                     request.Title,
		Summary:                   request.Summary,
		SequenceNo:                request.SequenceNo,
		DisplayNumber:             request.DisplayNumber,
		PublishState:              publishState,
		ReadAccessLevel:           readAccessLevel,
		InheritAccessFromManga:    request.InheritAccessFromManga,
		VIPOnly:                   request.VIPOnly,
		EarlyAccessEnabled:        request.EarlyAccessEnabled,
		EarlyAccessLevel:          earlyAccessLevel,
		EarlyAccessStartAt:        earlyStart,
		EarlyAccessEndAt:          earlyEnd,
		EarlyAccessFallbackAccess: fallbackAccess,
		PreviewEnabled:            request.PreviewEnabled,
		PreviewPageCount:          pagePreviewCount(len(pages), request.PreviewPageCount, request.PreviewEnabled),
		MediaHealthStatus:         entity.MediaHealthHealthy,
		IntegrityStatus:           entity.IntegrityUnknown,
		PageCount:                 len(pages),
		Pages:                     pages,
		ScheduledAt:               scheduledAt,
		PublishedAt:               publishedAt,
		ArchivedAt:                archivedAt,
		CreatedAt:                 now,
		UpdatedAt:                 now,
	}

	if chapter.DisplayNumber == "" {
		chapter.DisplayNumber = fmt.Sprintf("%d", chapter.SequenceNo)
	}

	if err := s.store.CreateChapter(ctx, chapter); err != nil {
		if errors.Is(err, chapterrepository.ErrConflict) {
			return dto.CreateChapterResponse{}, ErrChapterAlreadyExists
		}
		return dto.CreateChapterResponse{}, err
	}

	return dto.CreateChapterResponse{
		ChapterID:       chapter.ID,
		MangaID:         chapter.MangaID,
		SequenceNo:      chapter.SequenceNo,
		PublishState:    string(chapter.PublishState),
		ReadAccessLevel: string(chapter.ReadAccessLevel),
	}, nil
}

func toPages(input []dto.ChapterPageRequest, now time.Time) ([]entity.ChapterPage, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("%w: pages cannot be empty", ErrValidation)
	}

	result := make([]entity.ChapterPage, 0, len(input))
	seen := make(map[int]struct{}, len(input))
	for _, page := range input {
		if _, exists := seen[page.PageNumber]; exists {
			return nil, fmt.Errorf("%w: duplicate page_number", ErrValidation)
		}
		seen[page.PageNumber] = struct{}{}
		result = append(result, entity.ChapterPage{
			PageNumber: page.PageNumber,
			MediaURL:   page.MediaURL,
			Width:      page.Width,
			Height:     page.Height,
			LongStrip:  page.LongStrip,
			Checksum:   page.Checksum,
			CDNHealthy: true,
			UpdatedAt:  now,
		})
	}

	sort.Slice(result, func(i, j int) bool { return result[i].PageNumber < result[j].PageNumber })
	for i := 1; i < len(result); i++ {
		if result[i].PageNumber != result[i-1].PageNumber+1 {
			return nil, fmt.Errorf("%w: page_number must be continuous", ErrValidation)
		}
	}

	return result, nil
}

func createPublishTimes(state entity.PublishState, scheduledAt *time.Time, now time.Time) (*time.Time, *time.Time, *time.Time, error) {
	var scheduled *time.Time
	var published *time.Time
	var archived *time.Time

	switch state {
	case entity.PublishStateScheduled:
		if scheduledAt == nil {
			return nil, nil, nil, fmt.Errorf("%w: scheduled publish_state requires scheduled_at", ErrValidation)
		}
		resolved := scheduledAt.UTC()
		scheduled = &resolved
	case entity.PublishStatePublished:
		resolved := now
		published = &resolved
	case entity.PublishStateArchived:
		resolved := now
		archived = &resolved
	case entity.PublishStateDraft, entity.PublishStateUnpublished:
		// noop
	default:
		return nil, nil, nil, errors.Join(ErrValidation, errors.New("unsupported publish_state"))
	}

	return scheduled, published, archived, nil
}
