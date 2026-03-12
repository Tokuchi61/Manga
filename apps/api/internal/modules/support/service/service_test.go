package service

import (
	"context"
	"testing"
	"time"

	supportcontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/contract"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/dto"
	supportrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func newTestService(nowRef *time.Time) *SupportService {
	svc := New(supportrepository.NewMemoryStore(), validation.New())
	svc.now = func() time.Time { return nowRef.UTC() }
	return svc
}

func strPtr(v string) *string {
	return &v
}

func TestCreateOwnListDetailAndDuplicateFlow(t *testing.T) {
	now := time.Date(2026, 3, 12, 22, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	requesterID := uuid.NewString()
	targetID := uuid.NewString()

	communication, err := svc.CreateCommunication(ctx, dto.CreateCommunicationRequest{
		RequesterUserID: requesterID,
		Category:        "communication",
		Priority:        "normal",
		ReasonCode:      "general",
		ReasonText:      "Need help for account visibility",
		RequestID:       "req-comm-1",
	})
	require.NoError(t, err)
	require.Equal(t, "communication", communication.SupportKind)

	communicationRetry, err := svc.CreateCommunication(ctx, dto.CreateCommunicationRequest{
		RequesterUserID: requesterID,
		Category:        "communication",
		Priority:        "normal",
		ReasonCode:      "general",
		ReasonText:      "Different payload but same idempotency key",
		RequestID:       "req-comm-1",
	})
	require.NoError(t, err)
	require.Equal(t, communication.SupportID, communicationRetry.SupportID)

	reportA, err := svc.CreateReport(ctx, dto.CreateReportRequest{
		RequesterUserID: requesterID,
		Category:        "content",
		Priority:        "high",
		ReasonCode:      "spoiler",
		ReasonText:      "This comment contains spoiler text",
		TargetType:      "comment",
		TargetID:        targetID,
		RequestID:       "req-report-1",
	})
	require.NoError(t, err)
	require.Equal(t, "report", reportA.SupportKind)
	require.NotNil(t, reportA.TargetType)
	require.NotNil(t, reportA.TargetID)

	now = now.Add(2 * time.Minute)
	reportB, err := svc.CreateReport(ctx, dto.CreateReportRequest{
		RequesterUserID: requesterID,
		Category:        "content",
		Priority:        "high",
		ReasonCode:      "spoiler",
		ReasonText:      "Same spoiler issue reported again",
		TargetType:      "comment",
		TargetID:        targetID,
		RequestID:       "req-report-2",
	})
	require.NoError(t, err)
	require.NotNil(t, reportB.DuplicateOfSupportID)
	require.Equal(t, reportA.SupportID, *reportB.DuplicateOfSupportID)

	listing, err := svc.ListOwnSupport(ctx, dto.ListOwnSupportRequest{
		RequesterUserID: requesterID,
		SortBy:          "newest",
	})
	require.NoError(t, err)
	require.Equal(t, 3, listing.Count)

	detail, err := svc.GetSupportDetail(ctx, dto.GetSupportDetailRequest{
		SupportID:       reportB.SupportID,
		RequesterUserID: requesterID,
	})
	require.NoError(t, err)
	require.Equal(t, "report", detail.SupportKind)
	require.Equal(t, "content", detail.Category)
	require.Equal(t, "spoiler", detail.ReasonCode)

	_, err = svc.GetSupportDetail(ctx, dto.GetSupportDetailRequest{
		SupportID:       reportB.SupportID,
		RequesterUserID: uuid.NewString(),
	})
	require.ErrorIs(t, err, ErrForbiddenAction)
}

func TestReplyStatusResolveAndQueueFlow(t *testing.T) {
	now := time.Date(2026, 3, 12, 22, 30, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	requesterID := uuid.NewString()
	agentID := uuid.NewString()
	reviewerID := uuid.NewString()

	ticket, err := svc.CreateTicket(ctx, dto.CreateTicketRequest{
		RequesterUserID: requesterID,
		Category:        "operations",
		Priority:        "normal",
		ReasonCode:      "payment_issue",
		ReasonText:      "Payment has not been reflected",
		RequestID:       "req-ticket-1",
	})
	require.NoError(t, err)

	_, err = svc.AddReply(ctx, dto.AddSupportReplyRequest{
		SupportID:   ticket.SupportID,
		ActorUserID: requesterID,
		Message:     "Any update?",
		Visibility:  "public_to_requester",
	})
	require.NoError(t, err)

	afterRequesterReply, err := svc.GetSupportDetail(ctx, dto.GetSupportDetailRequest{
		SupportID:       ticket.SupportID,
		RequesterUserID: requesterID,
	})
	require.NoError(t, err)
	require.Equal(t, "waiting_team", afterRequesterReply.Status)

	_, err = svc.AddReply(ctx, dto.AddSupportReplyRequest{
		SupportID:   ticket.SupportID,
		ActorUserID: requesterID,
		Message:     "Internal note should fail",
		Visibility:  "internal_only",
	})
	require.ErrorIs(t, err, ErrForbiddenAction)

	_, err = svc.AddReply(ctx, dto.AddSupportReplyRequest{
		SupportID:   ticket.SupportID,
		ActorUserID: agentID,
		ActorIsTeam: true,
		Message:     "We are checking this now",
		Visibility:  "public_to_requester",
	})
	require.NoError(t, err)

	afterAgentReply, err := svc.GetSupportDetail(ctx, dto.GetSupportDetailRequest{
		SupportID:       ticket.SupportID,
		RequesterUserID: requesterID,
	})
	require.NoError(t, err)
	require.Equal(t, "waiting_user", afterAgentReply.Status)

	_, err = svc.UpdateStatus(ctx, dto.UpdateSupportStatusRequest{
		SupportID:        ticket.SupportID,
		Status:           "triaged",
		AssigneeUserID:   strPtr(agentID),
		ReviewedByUserID: strPtr(reviewerID),
	})
	require.NoError(t, err)

	_, err = svc.Resolve(ctx, dto.ResolveSupportRequest{
		SupportID:        ticket.SupportID,
		ReviewedByUserID: reviewerID,
		ResolutionNote:   "Issue confirmed and resolved.",
	})
	require.NoError(t, err)

	queue, err := svc.ListReviewQueue(ctx, dto.ListReviewQueueRequest{Status: "resolved", Limit: 10})
	require.NoError(t, err)
	require.Equal(t, 1, queue.Count)
	require.Equal(t, ticket.SupportID, queue.Items[0].SupportID)

	_, err = svc.AddReply(ctx, dto.AddSupportReplyRequest{
		SupportID:   ticket.SupportID,
		ActorUserID: agentID,
		ActorIsTeam: true,
		Message:     "Internal follow-up note",
		Visibility:  "internal_only",
	})
	require.NoError(t, err)

	detailWithoutInternal, err := svc.GetSupportDetail(ctx, dto.GetSupportDetailRequest{
		SupportID:       ticket.SupportID,
		RequesterUserID: requesterID,
	})
	require.NoError(t, err)
	require.Equal(t, 2, len(detailWithoutInternal.Replies))

	detailWithInternal, err := svc.GetSupportDetail(ctx, dto.GetSupportDetailRequest{
		SupportID:       ticket.SupportID,
		RequesterUserID: requesterID,
		IncludeInternal: true,
	})
	require.NoError(t, err)
	require.Equal(t, 3, len(detailWithInternal.Replies))
}

func TestModerationHandoffAndContractSignals(t *testing.T) {
	now := time.Date(2026, 3, 12, 23, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	requesterID := uuid.NewString()
	targetID := uuid.NewString()

	communication, err := svc.CreateCommunication(ctx, dto.CreateCommunicationRequest{
		RequesterUserID: requesterID,
		Category:        "communication",
		ReasonText:      "General question",
		RequestID:       "req-comm-contract",
	})
	require.NoError(t, err)

	_, err = svc.RequestModerationHandoff(ctx, dto.RequestModerationHandoffRequest{SupportID: communication.SupportID})
	require.ErrorIs(t, err, ErrModerationHandoffNotAllowed)

	report, err := svc.CreateReport(ctx, dto.CreateReportRequest{
		RequesterUserID: requesterID,
		Category:        "content",
		Priority:        "urgent",
		ReasonCode:      "abuse",
		ReasonText:      "Abusive comment reported",
		TargetType:      "comment",
		TargetID:        targetID,
		RequestID:       "req-report-contract",
	})
	require.NoError(t, err)

	_, err = svc.RequestModerationHandoff(ctx, dto.RequestModerationHandoffRequest{SupportID: report.SupportID})
	require.NoError(t, err)

	_, err = svc.RequestModerationHandoff(ctx, dto.RequestModerationHandoffRequest{SupportID: report.SupportID})
	require.ErrorIs(t, err, ErrAlreadyHandedOff)

	handoff, err := svc.GetModerationHandoffReference(ctx, report.SupportID, "", "corr-support-1")
	require.NoError(t, err)
	require.Equal(t, report.SupportID, handoff.SupportID)
	require.Equal(t, "report", handoff.SupportKind)
	require.Equal(t, "comment", handoff.TargetType)
	require.Equal(t, targetID, handoff.TargetID)
	require.Equal(t, "req-report-contract", handoff.RequestID)
	require.Equal(t, "corr-support-1", handoff.CorrelationID)

	notification, err := svc.BuildNotificationSignal(ctx, report.SupportID, supportcontract.EventSupportResolved, "", "corr-notify-1")
	require.NoError(t, err)
	require.Equal(t, supportcontract.EventSupportResolved, notification.Event)
	require.Equal(t, report.SupportID, notification.SupportID)
	require.Equal(t, requesterID, notification.RequesterUserID)
	require.Equal(t, "triaged", notification.Status)
	require.Equal(t, "req-report-contract", notification.RequestID)
	require.Equal(t, "corr-notify-1", notification.CorrelationID)
}
