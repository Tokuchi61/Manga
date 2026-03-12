package service

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/dto"
	moderationrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/repository"
	supportcontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/contract"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type supportLookupStub struct {
	reference supportcontract.ModerationHandoffReference
	err       error
}

func (s supportLookupStub) GetModerationHandoffReference(_ context.Context, _ string, _ string, _ string) (supportcontract.ModerationHandoffReference, error) {
	if s.err != nil {
		return supportcontract.ModerationHandoffReference{}, s.err
	}
	return s.reference, nil
}

type supportLinkerStub struct {
	linked map[string]string
}

func (s *supportLinkerStub) LinkModerationCase(_ context.Context, supportID string, moderationCaseID string) error {
	if s.linked == nil {
		s.linked = make(map[string]string)
	}
	s.linked[supportID] = moderationCaseID
	return nil
}

func TestCreateCaseFromSupportHandoffCreatesAndLinks(t *testing.T) {
	store := moderationrepository.NewMemoryStore()
	svc := New(store, validation.New())
	svc.now = func() time.Time { return time.Date(2026, 3, 12, 22, 0, 0, 0, time.UTC) }

	supportID := uuid.NewString()
	targetID := uuid.NewString()
	linker := &supportLinkerStub{}
	svc.SetSupportContracts(supportLookupStub{reference: supportcontract.ModerationHandoffReference{
		SupportID:     supportID,
		SupportKind:   "report",
		TargetType:    "comment",
		TargetID:      targetID,
		ReasonCode:    "abuse",
		RequestedAt:   svc.now(),
		RequestID:     "req-stage11-support",
		CorrelationID: "corr-stage11-support",
	}}, linker)

	response, err := svc.CreateCaseFromSupportHandoff(context.Background(), dto.CreateCaseFromSupportHandoffRequest{
		SupportID:   supportID,
		ActorUserID: uuid.NewString(),
	})
	require.NoError(t, err)
	require.True(t, response.Created)
	require.Equal(t, "queued", response.Case.Status)
	require.NotEmpty(t, response.Case.CaseID)
	require.Equal(t, response.Case.CaseID, linker.linked[supportID])
}

func TestCreateCaseFromSupportHandoffIsIdempotent(t *testing.T) {
	store := moderationrepository.NewMemoryStore()
	svc := New(store, validation.New())

	supportID := uuid.NewString()
	svc.SetSupportContracts(supportLookupStub{reference: supportcontract.ModerationHandoffReference{
		SupportID:   supportID,
		SupportKind: "report",
		TargetType:  "comment",
		TargetID:    uuid.NewString(),
	}}, &supportLinkerStub{})

	first, err := svc.CreateCaseFromSupportHandoff(context.Background(), dto.CreateCaseFromSupportHandoffRequest{
		SupportID:   supportID,
		ActorUserID: uuid.NewString(),
	})
	require.NoError(t, err)
	require.True(t, first.Created)

	second, err := svc.CreateCaseFromSupportHandoff(context.Background(), dto.CreateCaseFromSupportHandoffRequest{
		SupportID:   supportID,
		ActorUserID: uuid.NewString(),
	})
	require.NoError(t, err)
	require.False(t, second.Created)
	require.Equal(t, first.Case.CaseID, second.Case.CaseID)
}

func TestModerationLifecycleAssignNoteActionEscalate(t *testing.T) {
	store := moderationrepository.NewMemoryStore()
	svc := New(store, validation.New())

	supportID := uuid.NewString()
	svc.SetSupportContracts(supportLookupStub{reference: supportcontract.ModerationHandoffReference{
		SupportID:   supportID,
		SupportKind: "report",
		TargetType:  "comment",
		TargetID:    uuid.NewString(),
	}}, &supportLinkerStub{})

	created, err := svc.CreateCaseFromSupportHandoff(context.Background(), dto.CreateCaseFromSupportHandoffRequest{
		SupportID:   supportID,
		ActorUserID: uuid.NewString(),
	})
	require.NoError(t, err)
	caseID := created.Case.CaseID
	moderatorID := uuid.NewString()

	_, err = svc.AssignCase(context.Background(), dto.AssignCaseRequest{
		CaseID:      caseID,
		ActorUserID: moderatorID,
	})
	require.NoError(t, err)

	_, err = svc.AddModeratorNote(context.Background(), dto.AddModeratorNoteRequest{
		CaseID:      caseID,
		ActorUserID: moderatorID,
		Body:        "Need additional review data",
	})
	require.NoError(t, err)

	_, err = svc.ApplyAction(context.Background(), dto.ApplyActionRequest{
		CaseID:      caseID,
		ActorUserID: moderatorID,
		ActionType:  "hide",
		ReasonCode:  "abuse",
		Summary:     "hidden during review",
	})
	require.NoError(t, err)

	_, err = svc.EscalateCase(context.Background(), dto.EscalateCaseRequest{
		CaseID:      caseID,
		ActorUserID: moderatorID,
		Reason:      "requires admin final decision",
	})
	require.NoError(t, err)

	_, err = svc.EscalateCase(context.Background(), dto.EscalateCaseRequest{
		CaseID:      caseID,
		ActorUserID: moderatorID,
		Reason:      "retry escalation",
	})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrAlreadyEscalated)
}
