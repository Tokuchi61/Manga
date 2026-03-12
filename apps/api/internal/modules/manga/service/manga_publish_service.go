package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/entity"
	mangarepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/repository"
)

func (s *MangaService) UpdatePublishState(ctx context.Context, request dto.UpdatePublishStateRequest) (dto.OperationResponse, error) {
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

	now := s.now().UTC()
	action := normalizeValue(request.Action)
	scheduledAt, err := parseRFC3339Ptr(request.ScheduledAt, "scheduled_at")
	if err != nil {
		return dto.OperationResponse{}, err
	}

	switch action {
	case "draft":
		manga.PublishState = entity.PublishStateDraft
		manga.ScheduledAt = nil
		manga.ArchivedAt = nil
	case "schedule":
		if scheduledAt == nil {
			return dto.OperationResponse{}, fmt.Errorf("%w: schedule action requires scheduled_at", ErrValidation)
		}
		manga.PublishState = entity.PublishStateScheduled
		manga.ScheduledAt = scheduledAt
		manga.ArchivedAt = nil
	case "publish":
		manga.PublishState = entity.PublishStatePublished
		manga.PublishedAt = &now
		manga.ScheduledAt = nil
		manga.ArchivedAt = nil
	case "archive":
		manga.PublishState = entity.PublishStateArchived
		manga.ArchivedAt = &now
		manga.ScheduledAt = nil
	case "unpublish":
		manga.PublishState = entity.PublishStateUnpublished
		manga.ArchivedAt = nil
		manga.ScheduledAt = nil
	default:
		return dto.OperationResponse{}, fmt.Errorf("%w: unsupported publish action", ErrValidation)
	}

	manga.ContentVersion++
	manga.UpdatedAt = now

	if err := s.store.UpdateManga(ctx, manga); err != nil {
		if errors.Is(err, mangarepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrMangaNotFound
		}
		if errors.Is(err, mangarepository.ErrConflict) {
			return dto.OperationResponse{}, ErrMangaAlreadyExists
		}
		return dto.OperationResponse{}, err
	}

	return dto.OperationResponse{Status: "publish_state_updated"}, nil
}

func (s *MangaService) UpdateVisibility(ctx context.Context, request dto.UpdateVisibilityRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	mangaID, err := parseID(request.MangaID, "manga_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}

	visibility, err := toVisibility(request.Visibility)
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

	manga.Visibility = visibility
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
	return dto.OperationResponse{Status: "visibility_updated"}, nil
}

func (s *MangaService) UpdateEditorial(ctx context.Context, request dto.UpdateEditorialRequest) (dto.OperationResponse, error) {
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

	manga.Featured = boolValue(request.Featured, manga.Featured)
	manga.Recommended = boolValue(request.Recommended, manga.Recommended)
	if request.CollectionKeys != nil {
		manga.CollectionKeys = normalizeList(*request.CollectionKeys)
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

	return dto.OperationResponse{Status: "editorial_updated"}, nil
}

func (s *MangaService) SyncCounters(ctx context.Context, request dto.SyncCountersRequest) (dto.OperationResponse, error) {
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

	if request.ChapterCount != nil {
		if *request.ChapterCount < 0 {
			return dto.OperationResponse{}, ErrInvalidCounterValue
		}
		manga.ChapterCount = *request.ChapterCount
	}
	if request.CommentCount != nil {
		if *request.CommentCount < 0 {
			return dto.OperationResponse{}, ErrInvalidCounterValue
		}
		manga.CommentCount = *request.CommentCount
	}
	if request.ViewCount != nil {
		if *request.ViewCount < 0 {
			return dto.OperationResponse{}, ErrInvalidCounterValue
		}
		manga.ViewCount = *request.ViewCount
	}
	if request.SignalVersion != nil && *request.SignalVersion > manga.ContentVersion {
		manga.ContentVersion = *request.SignalVersion
	}
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

	return dto.OperationResponse{Status: "counters_synced"}, nil
}
