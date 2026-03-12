package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/entity"
	mangarepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/repository"
)

func (s *MangaService) SoftDeleteManga(ctx context.Context, request dto.SoftDeleteMangaRequest) (dto.OperationResponse, error) {
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
		return dto.OperationResponse{Status: "already_deleted"}, nil
	}

	now := s.now().UTC()
	manga.DeletedAt = &now
	manga.Visibility = entity.VisibilityHidden
	manga.PublishState = entity.PublishStateArchived
	manga.ArchivedAt = &now
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

	return dto.OperationResponse{Status: "soft_deleted"}, nil
}

func (s *MangaService) RestoreManga(ctx context.Context, request dto.RestoreMangaRequest) (dto.OperationResponse, error) {
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
	if manga.DeletedAt == nil {
		return dto.OperationResponse{Status: "already_active"}, nil
	}

	manga.DeletedAt = nil
	if manga.Visibility == entity.VisibilityHidden {
		manga.Visibility = entity.VisibilityPublic
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

	return dto.OperationResponse{Status: "restored"}, nil
}
