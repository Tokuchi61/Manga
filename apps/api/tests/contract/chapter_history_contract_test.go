package contract_test

import (
	"testing"
	"time"

	chaptercontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/contract"
	"github.com/stretchr/testify/require"
)

func TestChapterResumeAnchorContractShape(t *testing.T) {
	now := time.Date(2026, 3, 12, 18, 0, 0, 0, time.UTC)
	anchor := chaptercontract.ResumeAnchor{
		ChapterID:  "9f8538ab-d6eb-4d91-8e8f-5d90e45c052f",
		MangaID:    "4de443dd-2644-4f2b-9039-5e8df0538f8f",
		PageNumber: 7,
		PageCount:  19,
		UpdatedAt:  now,
	}

	require.NotEmpty(t, anchor.ChapterID)
	require.NotEmpty(t, anchor.MangaID)
	require.Equal(t, 7, anchor.PageNumber)
	require.Equal(t, 19, anchor.PageCount)
	require.Equal(t, now, anchor.UpdatedAt)
}

func TestChapterReadSignalContractShape(t *testing.T) {
	now := time.Date(2026, 3, 12, 18, 5, 0, 0, time.UTC)
	signal := chaptercontract.ReadSignal{
		Event:         chaptercontract.EventReadFinished,
		ChapterID:     "9f8538ab-d6eb-4d91-8e8f-5d90e45c052f",
		MangaID:       "4de443dd-2644-4f2b-9039-5e8df0538f8f",
		PageNumber:    19,
		PageCount:     19,
		OccurredAt:    now,
		RequestID:     "req-stage8",
		CorrelationID: "corr-stage8",
	}

	require.Equal(t, chaptercontract.EventReadFinished, signal.Event)
	require.Equal(t, 19, signal.PageNumber)
	require.Equal(t, 19, signal.PageCount)
	require.Equal(t, "req-stage8", signal.RequestID)
	require.Equal(t, "corr-stage8", signal.CorrelationID)
	require.Equal(t, now, signal.OccurredAt)
}
