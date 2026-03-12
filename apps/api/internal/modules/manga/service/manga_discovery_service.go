package service

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/dto"
	mangarepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/repository"
)

func (s *MangaService) ListDiscovery(ctx context.Context, request dto.DiscoveryRequest) (dto.DiscoveryResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.DiscoveryResponse{}, err
	}

	mode := normalizeValue(request.Mode)
	if mode == "" {
		mode = "all"
	}
	sortBy := request.SortBy
	if sortBy == "" {
		sortBy = "popular"
	}
	limit := request.Limit
	if limit <= 0 {
		limit = 20
	}

	items, err := s.store.ListManga(ctx, mangarepository.ListQuery{
		SortBy:             sortBy,
		Limit:              limit,
		Offset:             0,
		IncludeUnpublished: false,
		IncludeHidden:      false,
		IncludeDeleted:     false,
	})
	if err != nil {
		return dto.DiscoveryResponse{}, err
	}

	collectionKey := normalizeValue(request.CollectionKey)
	result := make([]dto.MangaListItemResponse, 0, len(items))
	for _, item := range items {
		include := false
		switch mode {
		case "recommended":
			include = item.Recommended
		case "featured":
			include = item.Featured
		case "collection":
			include = collectionKey != "" && hasCollectionKey(item.CollectionKeys, collectionKey)
		default:
			include = item.Featured || item.Recommended || len(item.CollectionKeys) > 0
		}
		if include {
			result = append(result, toListItem(item))
		}
	}

	return dto.DiscoveryResponse{Mode: mode, Items: result, Count: len(result)}, nil
}

func hasCollectionKey(keys []string, expected string) bool {
	for _, key := range keys {
		if normalizeValue(key) == expected {
			return true
		}
	}
	return false
}
