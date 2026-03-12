package service

import (
	"context"
	"testing"
	"time"

	chaptercontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/contract"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/dto"
	chapterrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func newTestService(nowRef *time.Time) *ChapterService {
	svc := New(chapterrepository.NewMemoryStore(), validation.New())
	svc.now = func() time.Time { return nowRef.UTC() }
	return svc
}

func testPages(count int) []dto.ChapterPageRequest {
	items := make([]dto.ChapterPageRequest, 0, count)
	for i := 1; i <= count; i++ {
		items = append(items, dto.ChapterPageRequest{
			PageNumber: i,
			MediaURL:   "https://cdn.example.com/chapter/page.jpg",
			Width:      1200,
			Height:     1800,
		})
	}
	return items
}

func boolPtr(v bool) *bool {
	return &v
}

func intPtr(v int) *int {
	return &v
}

func strPtr(v string) *string {
	return &v
}

func TestCreatePublishReadAndDetailFlow(t *testing.T) {
	now := time.Date(2026, 3, 12, 16, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	earlyStart := now.Add(-30 * time.Minute)
	earlyEnd := now.Add(30 * time.Minute)

	created, err := svc.CreateChapter(ctx, dto.CreateChapterRequest{
		MangaID:                   uuid.NewString(),
		Title:                     "Chapter 1",
		SequenceNo:                1,
		PreviewEnabled:            true,
		PreviewPageCount:          2,
		EarlyAccessEnabled:        true,
		EarlyAccessLevel:          "vip",
		EarlyAccessFallbackAccess: "authenticated",
		EarlyAccessStartAt:        &earlyStart,
		EarlyAccessEndAt:          &earlyEnd,
		Pages:                     testPages(4),
	})
	require.NoError(t, err)
	require.NotEmpty(t, created.ChapterID)
	require.Equal(t, "draft", created.PublishState)

	_, err = svc.GetChapterDetail(ctx, dto.GetChapterDetailRequest{ChapterID: created.ChapterID})
	require.ErrorIs(t, err, ErrChapterNotVisible)

	_, err = svc.UpdatePublishState(ctx, dto.UpdatePublishStateRequest{
		ChapterID: created.ChapterID,
		Action:    "publish",
	})
	require.NoError(t, err)

	detail, err := svc.GetChapterDetail(ctx, dto.GetChapterDetailRequest{ChapterID: created.ChapterID})
	require.NoError(t, err)
	require.Equal(t, "published", detail.PublishState)
	require.Equal(t, 4, detail.PageCount)

	previewRead, err := svc.ReadChapter(ctx, dto.ReadChapterRequest{
		ChapterID: created.ChapterID,
		Mode:      "preview",
		At:        &now,
	})
	require.NoError(t, err)
	require.Equal(t, 2, previewRead.PageCount)
	require.True(t, previewRead.EarlyAccessActive)

	endedAt := now.Add(2 * time.Hour)
	fullRead, err := svc.ReadChapter(ctx, dto.ReadChapterRequest{
		ChapterID: created.ChapterID,
		Mode:      "full",
		At:        &endedAt,
	})
	require.NoError(t, err)
	require.Equal(t, 4, fullRead.PageCount)
	require.False(t, fullRead.EarlyAccessActive)
}

func TestNavigationReorderAndLifecycleFlow(t *testing.T) {
	now := time.Date(2026, 3, 12, 16, 30, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()
	mangaID := uuid.NewString()

	createPublished := func(title string, sequence int) string {
		created, err := svc.CreateChapter(ctx, dto.CreateChapterRequest{
			MangaID:      mangaID,
			Title:        title,
			SequenceNo:   sequence,
			PublishState: "published",
			Pages:        testPages(2),
		})
		require.NoError(t, err)
		return created.ChapterID
	}

	chapter1ID := createPublished("Chapter 1", 1)
	chapter2ID := createPublished("Chapter 2", 2)
	chapter3ID := createPublished("Chapter 3", 3)

	navigation, err := svc.GetNavigation(ctx, dto.NavigationRequest{ChapterID: chapter2ID})
	require.NoError(t, err)
	require.NotNil(t, navigation.PreviousChapterID)
	require.Equal(t, chapter1ID, *navigation.PreviousChapterID)
	require.NotNil(t, navigation.NextChapterID)
	require.Equal(t, chapter3ID, *navigation.NextChapterID)

	_, err = svc.ReorderChapter(ctx, dto.ReorderChapterRequest{ChapterID: chapter2ID, SequenceNo: 4})
	require.NoError(t, err)

	navigationAfterReorder, err := svc.GetNavigation(ctx, dto.NavigationRequest{ChapterID: chapter2ID})
	require.NoError(t, err)
	require.NotNil(t, navigationAfterReorder.PreviousChapterID)
	require.Equal(t, chapter3ID, *navigationAfterReorder.PreviousChapterID)
	require.Nil(t, navigationAfterReorder.NextChapterID)

	_, err = svc.SoftDeleteChapter(ctx, dto.SoftDeleteChapterRequest{ChapterID: chapter2ID})
	require.NoError(t, err)
	_, err = svc.ReadChapter(ctx, dto.ReadChapterRequest{ChapterID: chapter2ID, Mode: "full"})
	require.ErrorIs(t, err, ErrChapterNotVisible)

	_, err = svc.RestoreChapter(ctx, dto.RestoreChapterRequest{ChapterID: chapter2ID})
	require.NoError(t, err)
	_, err = svc.UpdatePublishState(ctx, dto.UpdatePublishStateRequest{ChapterID: chapter2ID, Action: "publish"})
	require.NoError(t, err)

	_, err = svc.ReadChapter(ctx, dto.ReadChapterRequest{ChapterID: chapter2ID, Mode: "full"})
	require.NoError(t, err)
}

func TestAccessMediaIntegrityAndHistoryContractFlow(t *testing.T) {
	now := time.Date(2026, 3, 12, 17, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()
	mangaID := uuid.NewString()

	created, err := svc.CreateChapter(ctx, dto.CreateChapterRequest{
		MangaID:      mangaID,
		Title:        "VIP Chapter",
		SequenceNo:   1,
		PublishState: "published",
		Pages:        testPages(3),
	})
	require.NoError(t, err)

	_, err = svc.UpdateAccess(ctx, dto.UpdateAccessRequest{
		ChapterID:          created.ChapterID,
		VIPOnly:            boolPtr(true),
		EarlyAccessEnabled: boolPtr(true),
		PreviewEnabled:     boolPtr(true),
		PreviewPageCount:   intPtr(10),
	})
	require.NoError(t, err)

	read, err := svc.ReadChapter(ctx, dto.ReadChapterRequest{ChapterID: created.ChapterID, Mode: "full"})
	require.NoError(t, err)
	require.True(t, read.VIPOnly)
	require.Equal(t, "vip", read.ReadAccessLevel)
	require.False(t, read.EarlyAccessEnabled)

	_, err = svc.UpdateMediaHealth(ctx, dto.UpdateMediaHealthRequest{ChapterID: created.ChapterID, MediaHealthStatus: "broken"})
	require.NoError(t, err)
	_, err = svc.UpdateIntegrity(ctx, dto.UpdateIntegrityRequest{ChapterID: created.ChapterID, IntegrityStatus: "failed"})
	require.NoError(t, err)

	detail, err := svc.GetChapterDetail(ctx, dto.GetChapterDetailRequest{ChapterID: created.ChapterID})
	require.NoError(t, err)
	require.Equal(t, "broken", detail.MediaHealthStatus)
	require.Equal(t, "failed", detail.IntegrityStatus)
	require.Equal(t, 3, detail.PreviewPageCount)

	anchor, err := svc.GetResumeAnchor(ctx, created.ChapterID, 99)
	require.NoError(t, err)
	require.Equal(t, created.ChapterID, anchor.ChapterID)
	require.Equal(t, mangaID, anchor.MangaID)
	require.Equal(t, 3, anchor.PageNumber)
	require.Equal(t, 3, anchor.PageCount)

	signal := svc.BuildReadSignal(created.ChapterID, mangaID, -1, -2, "", "req-1", "corr-1")
	require.Equal(t, chaptercontract.EventReadCheckpoint, signal.Event)
	require.Equal(t, 0, signal.PageNumber)
	require.Equal(t, 0, signal.PageCount)
	require.Equal(t, "req-1", signal.RequestID)
	require.Equal(t, "corr-1", signal.CorrelationID)

	_, err = svc.UpdateAccess(ctx, dto.UpdateAccessRequest{
		ChapterID:          created.ChapterID,
		EarlyAccessStartAt: strPtr("2026-03-12T18:00:00Z"),
		EarlyAccessEndAt:   strPtr("2026-03-12T17:00:00Z"),
		EarlyAccessEnabled: boolPtr(true),
		EarlyAccessLevel:   strPtr("vip"),
	})
	require.ErrorIs(t, err, ErrValidation)
}
