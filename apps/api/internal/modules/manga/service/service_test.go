package service

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/dto"
	mangarepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/stretchr/testify/require"
)

func newTestService(nowRef *time.Time) *MangaService {
	svc := New(mangarepository.NewMemoryStore(), validation.New())
	svc.now = func() time.Time { return nowRef.UTC() }
	return svc
}

func TestCreatePublishAndPublicDetailFlow(t *testing.T) {
	now := time.Date(2026, 3, 12, 13, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	created, err := svc.CreateManga(ctx, dto.CreateMangaRequest{
		Title:   "Moonblade Chronicles",
		Summary: "A fantasy action saga",
		Genres:  []string{"action", "fantasy"},
		Tags:    []string{"adventure"},
	})
	require.NoError(t, err)
	require.NotEmpty(t, created.MangaID)
	require.Equal(t, "draft", created.PublishState)

	_, err = svc.GetPublicMangaDetail(ctx, dto.GetMangaDetailRequest{MangaID: created.MangaID})
	require.ErrorIs(t, err, ErrMangaNotVisible)

	_, err = svc.UpdatePublishState(ctx, dto.UpdatePublishStateRequest{MangaID: created.MangaID, Action: "publish"})
	require.NoError(t, err)

	detail, err := svc.GetPublicMangaDetail(ctx, dto.GetMangaDetailRequest{MangaID: created.MangaID})
	require.NoError(t, err)
	require.Equal(t, "moonblade-chronicles", detail.Slug)
	require.Equal(t, "published", detail.PublishState)
}

func TestListingDiscoveryCountersAndLifecycleFlow(t *testing.T) {
	now := time.Date(2026, 3, 12, 13, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	first, err := svc.CreateManga(ctx, dto.CreateMangaRequest{
		Title:   "Nebula Academy",
		Summary: "School and mystery",
		Genres:  []string{"drama"},
		Tags:    []string{"school"},
	})
	require.NoError(t, err)
	_, err = svc.UpdatePublishState(ctx, dto.UpdatePublishStateRequest{MangaID: first.MangaID, Action: "publish"})
	require.NoError(t, err)

	second, err := svc.CreateManga(ctx, dto.CreateMangaRequest{
		Title:   "Blade Frontier",
		Summary: "Battle focused",
		Genres:  []string{"action"},
		Tags:    []string{"battle"},
	})
	require.NoError(t, err)
	_, err = svc.UpdatePublishState(ctx, dto.UpdatePublishStateRequest{MangaID: second.MangaID, Action: "publish"})
	require.NoError(t, err)

	viewCount := int64(120)
	_, err = svc.SyncCounters(ctx, dto.SyncCountersRequest{MangaID: second.MangaID, ViewCount: &viewCount})
	require.NoError(t, err)

	recommended := true
	collections := []string{"editorial_spring"}
	_, err = svc.UpdateEditorial(ctx, dto.UpdateEditorialRequest{MangaID: second.MangaID, Recommended: &recommended, CollectionKeys: &collections})
	require.NoError(t, err)

	listed, err := svc.ListPublicManga(ctx, dto.ListMangaRequest{SortBy: "popular"})
	require.NoError(t, err)
	require.Len(t, listed.Items, 2)
	require.Equal(t, second.MangaID, listed.Items[0].MangaID)

	discovery, err := svc.ListDiscovery(ctx, dto.DiscoveryRequest{Mode: "recommended"})
	require.NoError(t, err)
	require.Len(t, discovery.Items, 1)
	require.Equal(t, second.MangaID, discovery.Items[0].MangaID)

	_, err = svc.SoftDeleteManga(ctx, dto.SoftDeleteMangaRequest{MangaID: second.MangaID})
	require.NoError(t, err)

	_, err = svc.GetPublicMangaDetail(ctx, dto.GetMangaDetailRequest{MangaID: second.MangaID})
	require.ErrorIs(t, err, ErrMangaNotVisible)

	_, err = svc.RestoreManga(ctx, dto.RestoreMangaRequest{MangaID: second.MangaID})
	require.NoError(t, err)
}

func TestCreateMangaSlugConflict(t *testing.T) {
	now := time.Date(2026, 3, 12, 13, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	_, err := svc.CreateManga(ctx, dto.CreateMangaRequest{
		Title:   "Echo Realm",
		Slug:    "echo-realm",
		Summary: "summary",
		Genres:  []string{"fantasy"},
	})
	require.NoError(t, err)

	_, err = svc.CreateManga(ctx, dto.CreateMangaRequest{
		Title:   "Echo Realm Two",
		Slug:    "echo-realm",
		Summary: "summary",
		Genres:  []string{"fantasy"},
	})
	require.ErrorIs(t, err, ErrMangaAlreadyExists)
}
