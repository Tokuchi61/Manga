package service

import (
	"context"
	"errors"

	mangacontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/contract"
	mangarepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/repository"
)

// GetChapterAccessDefaults exposes manga-owned chapter default access metadata.
func (s *MangaService) GetChapterAccessDefaults(ctx context.Context, mangaID string) (mangacontract.ChapterAccessDefaults, error) {
	parsedID, err := parseID(mangaID, "manga_id")
	if err != nil {
		return mangacontract.ChapterAccessDefaults{}, err
	}

	manga, err := s.store.GetMangaByID(ctx, parsedID)
	if err != nil {
		if errors.Is(err, mangarepository.ErrNotFound) {
			return mangacontract.ChapterAccessDefaults{}, ErrMangaNotFound
		}
		return mangacontract.ChapterAccessDefaults{}, err
	}

	return mangacontract.ChapterAccessDefaults{
		MangaID:            manga.ID,
		ReadAccessLevel:    string(manga.DefaultReadAccess),
		EarlyAccessEnabled: manga.EarlyAccessEnabled,
		EarlyAccessLevel:   string(manga.EarlyAccessLevel),
		UpdatedAt:          manga.UpdatedAt,
	}, nil
}
