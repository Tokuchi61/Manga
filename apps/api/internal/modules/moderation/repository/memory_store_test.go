package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStoreCreateGetAndSourceRefLookup(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 12, 20, 0, 0, 0, time.UTC)
	sourceRef := uuid.NewString()
	moderationCase := entity.Case{
		ID:               uuid.NewString(),
		Source:           entity.CaseSourceSupportReport,
		SourceRefID:      &sourceRef,
		TargetType:       entity.TargetTypeComment,
		TargetID:         uuid.NewString(),
		Status:           entity.CaseStatusQueued,
		AssignmentStatus: entity.AssignmentStatusUnassigned,
		EscalationStatus: entity.EscalationStatusNotEscalated,
		ActionResult:     entity.ActionResultNone,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	err := store.CreateCase(context.Background(), moderationCase)
	require.NoError(t, err)

	persisted, err := store.GetCaseByID(context.Background(), moderationCase.ID)
	require.NoError(t, err)
	require.Equal(t, moderationCase.ID, persisted.ID)

	bySourceRef, err := store.GetCaseBySourceRef(context.Background(), entity.CaseSourceSupportReport, sourceRef)
	require.NoError(t, err)
	require.Equal(t, moderationCase.ID, bySourceRef.ID)
}

func TestMemoryStoreRejectsDuplicateSourceRef(t *testing.T) {
	store := NewMemoryStore()
	sourceRef := uuid.NewString()
	base := entity.Case{
		Source:           entity.CaseSourceSupportReport,
		SourceRefID:      &sourceRef,
		TargetType:       entity.TargetTypeComment,
		TargetID:         uuid.NewString(),
		Status:           entity.CaseStatusQueued,
		AssignmentStatus: entity.AssignmentStatusUnassigned,
		EscalationStatus: entity.EscalationStatusNotEscalated,
		ActionResult:     entity.ActionResultNone,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	first := base
	first.ID = uuid.NewString()
	second := base
	second.ID = uuid.NewString()

	require.NoError(t, store.CreateCase(context.Background(), first))
	err := store.CreateCase(context.Background(), second)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrConflict)
}

func TestMemoryStoreListQueueFilters(t *testing.T) {
	store := NewMemoryStore()
	assignee := uuid.NewString()
	now := time.Now().UTC()

	items := []entity.Case{
		{
			ID:               uuid.NewString(),
			Source:           entity.CaseSourceSupportReport,
			TargetType:       entity.TargetTypeComment,
			TargetID:         uuid.NewString(),
			Status:           entity.CaseStatusQueued,
			AssignmentStatus: entity.AssignmentStatusUnassigned,
			EscalationStatus: entity.EscalationStatusNotEscalated,
			ActionResult:     entity.ActionResultNone,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:                      uuid.NewString(),
			Source:                  entity.CaseSourceSupportReport,
			TargetType:              entity.TargetTypeManga,
			TargetID:                uuid.NewString(),
			Status:                  entity.CaseStatusAssigned,
			AssignmentStatus:        entity.AssignmentStatusAssigned,
			AssignedModeratorUserID: &assignee,
			EscalationStatus:        entity.EscalationStatusNotEscalated,
			ActionResult:            entity.ActionResultNone,
			CreatedAt:               now.Add(1 * time.Minute),
			UpdatedAt:               now.Add(1 * time.Minute),
		},
	}
	for _, item := range items {
		require.NoError(t, store.CreateCase(context.Background(), item))
	}

	result, err := store.ListQueue(context.Background(), QueueQuery{Status: "assigned", AssignedModeratorUserID: assignee})
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, entity.CaseStatusAssigned, result[0].Status)
	require.Equal(t, entity.TargetTypeManga, result[0].TargetType)
}
