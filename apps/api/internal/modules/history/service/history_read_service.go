package service

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/dto"
	historyrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/history/repository"
)

func (s *HistoryService) ListContinueReading(ctx context.Context, request dto.ListContinueReadingRequest) (dto.ListContinueReadingResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListContinueReadingResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.ListContinueReadingResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ListContinueReadingResponse{}, err
	}
	if !cfg.ContinueReadingEnabled {
		return dto.ListContinueReadingResponse{}, ErrContinueReadingDisabled
	}

	limit := request.Limit
	if limit <= 0 {
		limit = 50
	}

	entries, err := s.store.ListContinueReading(ctx, historyrepository.ContinueReadingQuery{
		UserID: userID,
		Limit:  limit,
		Offset: request.Offset,
		SortBy: parseSortBy(request.SortBy, "newest"),
	})
	if err != nil {
		return dto.ListContinueReadingResponse{}, err
	}

	items := make([]dto.HistoryEntryResponse, 0, len(entries))
	for _, entry := range entries {
		items = append(items, mapEntry(entry, false))
	}

	return dto.ListContinueReadingResponse{Items: items, Count: len(items)}, nil
}

func (s *HistoryService) ListLibrary(ctx context.Context, request dto.ListLibraryRequest) (dto.ListLibraryResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListLibraryResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.ListLibraryResponse{}, err
	}
	status, err := parseOptionalStatus(request.Status)
	if err != nil {
		return dto.ListLibraryResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ListLibraryResponse{}, err
	}
	if !cfg.LibraryEnabled {
		return dto.ListLibraryResponse{}, ErrLibraryDisabled
	}

	limit := request.Limit
	if limit <= 0 {
		limit = 50
	}

	entries, err := s.store.ListLibrary(ctx, historyrepository.LibraryQuery{
		UserID:     userID,
		Status:     string(status),
		Bookmarked: request.Bookmarked,
		Favorited:  request.Favorited,
		SharedOnly: request.SharedOnly,
		Limit:      limit,
		Offset:     request.Offset,
		SortBy:     parseSortBy(request.SortBy, "newest"),
	})
	if err != nil {
		return dto.ListLibraryResponse{}, err
	}

	items := make([]dto.HistoryEntryResponse, 0, len(entries))
	for _, entry := range entries {
		items = append(items, mapEntry(entry, false))
	}

	return dto.ListLibraryResponse{Items: items, Count: len(items)}, nil
}

func (s *HistoryService) ListPublicLibrary(ctx context.Context, request dto.ListPublicLibraryRequest) (dto.ListLibraryResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListLibraryResponse{}, err
	}

	ownerUserID, err := parseID(request.OwnerUserID, "owner_user_id")
	if err != nil {
		return dto.ListLibraryResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ListLibraryResponse{}, err
	}
	if !cfg.LibraryEnabled {
		return dto.ListLibraryResponse{}, ErrLibraryDisabled
	}

	limit := request.Limit
	if limit <= 0 {
		limit = 50
	}

	entries, err := s.store.ListLibrary(ctx, historyrepository.LibraryQuery{
		UserID:     ownerUserID,
		SharedOnly: true,
		Limit:      limit,
		Offset:     request.Offset,
		SortBy:     parseSortBy(request.SortBy, "newest"),
	})
	if err != nil {
		return dto.ListLibraryResponse{}, err
	}

	items := make([]dto.HistoryEntryResponse, 0, len(entries))
	for _, entry := range entries {
		items = append(items, mapEntry(entry, true))
	}

	return dto.ListLibraryResponse{Items: items, Count: len(items)}, nil
}

func (s *HistoryService) ListTimeline(ctx context.Context, request dto.ListTimelineRequest) (dto.ListTimelineResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListTimelineResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.ListTimelineResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ListTimelineResponse{}, err
	}
	if !cfg.TimelineEnabled {
		return dto.ListTimelineResponse{}, ErrTimelineDisabled
	}

	limit := request.Limit
	if limit <= 0 {
		limit = 50
	}

	items, err := s.store.ListTimeline(ctx, historyrepository.TimelineQuery{
		UserID: userID,
		Event:  parseOptionalTimelineEvent(request.Event),
		Limit:  limit,
		Offset: request.Offset,
		SortBy: parseSortBy(request.SortBy, "newest"),
	})
	if err != nil {
		return dto.ListTimelineResponse{}, err
	}

	result := make([]dto.TimelineItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, mapTimeline(item))
	}

	return dto.ListTimelineResponse{Items: result, Count: len(result)}, nil
}
