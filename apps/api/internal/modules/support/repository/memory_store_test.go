package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func targetTypePtr(v entity.SupportTargetType) *entity.SupportTargetType {
	return &v
}

func stringPtr(v string) *string {
	return &v
}

func TestMemoryStoreRequestIDIndexAndGetByRequester(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 12, 22, 0, 0, 0, time.UTC)
	requesterID := uuid.NewString()

	support := entity.SupportCase{
		ID:              uuid.NewString(),
		RequesterUserID: requesterID,
		Kind:            entity.SupportKindTicket,
		Category:        "operations",
		Priority:        entity.SupportPriorityNormal,
		ReasonCode:      "payment_issue",
		ReasonText:      "Payment mismatch",
		Status:          entity.SupportStatusOpen,
		RequestID:       "req-100",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	require.NoError(t, store.CreateCase(context.Background(), support))

	loadedByRequestID, err := store.FindCaseByRequesterRequestID(context.Background(), requesterID, "req-100")
	require.NoError(t, err)
	require.Equal(t, support.ID, loadedByRequestID.ID)

	err = store.CreateCase(context.Background(), support)
	require.ErrorIs(t, err, ErrConflict)
}

func TestMemoryStoreQueueAndRecentSimilarLookup(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 12, 22, 30, 0, 0, time.UTC)
	requesterID := uuid.NewString()
	targetID := uuid.NewString()

	older := entity.SupportCase{
		ID:              uuid.NewString(),
		RequesterUserID: requesterID,
		Kind:            entity.SupportKindReport,
		Category:        "content",
		Priority:        entity.SupportPriorityHigh,
		ReasonCode:      "spoiler",
		ReasonText:      "Spoiler report",
		TargetType:      targetTypePtr(entity.SupportTargetTypeComment),
		TargetID:        stringPtr(targetID),
		Status:          entity.SupportStatusOpen,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	newer := older
	newer.ID = uuid.NewString()
	newer.Priority = entity.SupportPriorityUrgent
	newer.CreatedAt = now.Add(2 * time.Minute)
	newer.UpdatedAt = now.Add(2 * time.Minute)

	resolved := older
	resolved.ID = uuid.NewString()
	resolved.Priority = entity.SupportPriorityNormal
	resolved.Status = entity.SupportStatusResolved
	resolved.CreatedAt = now.Add(1 * time.Minute)
	resolved.UpdatedAt = now.Add(1 * time.Minute)

	require.NoError(t, store.CreateCase(context.Background(), older))
	require.NoError(t, store.CreateCase(context.Background(), newer))
	require.NoError(t, store.CreateCase(context.Background(), resolved))

	queue, err := store.ListReviewQueue(context.Background(), QueueQuery{Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.Len(t, queue, 2)
	require.Equal(t, newer.ID, queue[0].ID)
	require.Equal(t, older.ID, queue[1].ID)

	recent, err := store.FindRecentSimilarCase(
		context.Background(),
		requesterID,
		entity.SupportKindReport,
		targetTypePtr(entity.SupportTargetTypeComment),
		stringPtr(targetID),
		"spoiler",
		now.Add(-1*time.Minute),
	)
	require.NoError(t, err)
	require.Equal(t, newer.ID, recent.ID)
}
