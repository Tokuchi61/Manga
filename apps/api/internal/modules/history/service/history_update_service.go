package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/entity"
	historyrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/history/repository"
	"github.com/google/uuid"
)

func (s *HistoryService) UpdateBookmark(ctx context.Context, request dto.UpdateBookmarkRequest) (dto.HistoryEntryResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.HistoryEntryResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.HistoryEntryResponse{}, err
	}
	mangaID, err := parseID(request.MangaID, "manga_id")
	if err != nil {
		return dto.HistoryEntryResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.HistoryEntryResponse{}, err
	}
	if !cfg.BookmarkWriteEnabled {
		return dto.HistoryEntryResponse{}, ErrBookmarkWriteDisabled
	}

	entry, err := s.store.GetLibraryEntry(ctx, userID, mangaID)
	if err != nil {
		if !errors.Is(err, historyrepository.ErrNotFound) {
			return dto.HistoryEntryResponse{}, err
		}
		now := s.now().UTC()
		entry = entity.LibraryEntry{
			ID:         uuid.NewString(),
			UserID:     userID,
			MangaID:    mangaID,
			Status:     entity.ReadingStatusDropped,
			LastReadAt: now,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
	}

	now := s.now().UTC()
	entry.Bookmarked = request.Bookmarked
	entry.Favorited = request.Favorited
	entry.UpdatedAt = now
	if entry.LastReadAt.IsZero() {
		entry.LastReadAt = now
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = now
	}

	persisted, err := s.store.UpsertLibraryEntry(ctx, entry)
	if err != nil {
		return dto.HistoryEntryResponse{}, err
	}

	return mapEntry(persisted, false), nil
}

func (s *HistoryService) UpdateShare(ctx context.Context, request dto.UpdateShareRequest) (dto.HistoryEntryResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.HistoryEntryResponse{}, err
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.HistoryEntryResponse{}, err
	}
	mangaID, err := parseID(request.MangaID, "manga_id")
	if err != nil {
		return dto.HistoryEntryResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.HistoryEntryResponse{}, err
	}
	if !cfg.LibraryEnabled {
		return dto.HistoryEntryResponse{}, ErrLibraryDisabled
	}

	entry, err := s.store.GetLibraryEntry(ctx, userID, mangaID)
	if err != nil {
		if errors.Is(err, historyrepository.ErrNotFound) {
			return dto.HistoryEntryResponse{}, ErrNotFound
		}
		return dto.HistoryEntryResponse{}, err
	}

	entry.SharePublic = request.SharePublic
	entry.UpdatedAt = s.now().UTC()

	persisted, err := s.store.UpsertLibraryEntry(ctx, entry)
	if err != nil {
		return dto.HistoryEntryResponse{}, err
	}

	return mapEntry(persisted, false), nil
}
