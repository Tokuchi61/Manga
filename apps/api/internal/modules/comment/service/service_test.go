package service

import (
	"context"
	"testing"
	"time"

	commentcontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/contract"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/dto"
	commentrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func newTestService(nowRef *time.Time) *CommentService {
	svc := New(commentrepository.NewMemoryStore(), validation.New())
	svc.now = func() time.Time { return nowRef.UTC() }
	return svc
}

func TestCreateListThreadDepthAndRateLimitFlow(t *testing.T) {
	current := time.Date(2026, 3, 12, 20, 0, 0, 0, time.UTC)
	svc := newTestService(&current)
	ctx := context.Background()
	targetID := uuid.NewString()

	authorA := uuid.NewString()
	rootA, err := svc.CreateComment(ctx, dto.CreateCommentRequest{
		TargetType:   "manga",
		TargetID:     targetID,
		AuthorUserID: authorA,
		Content:      "Root A",
	})
	require.NoError(t, err)
	require.Equal(t, 0, rootA.Depth)

	_, err = svc.CreateComment(ctx, dto.CreateCommentRequest{
		TargetType:   "manga",
		TargetID:     targetID,
		AuthorUserID: authorA,
		Content:      "Root A second",
	})
	require.ErrorIs(t, err, ErrRateLimited)

	current = current.Add(20 * time.Second)
	rootB, err := svc.CreateComment(ctx, dto.CreateCommentRequest{
		TargetType:   "manga",
		TargetID:     targetID,
		AuthorUserID: authorA,
		Content:      "Root B",
	})
	require.NoError(t, err)
	require.NotEmpty(t, rootB.CommentID)

	reply1, err := svc.CreateComment(ctx, dto.CreateCommentRequest{
		TargetType:      "manga",
		TargetID:        targetID,
		AuthorUserID:    uuid.NewString(),
		ParentCommentID: &rootA.CommentID,
		Content:         "Reply 1",
	})
	require.NoError(t, err)
	require.Equal(t, 1, reply1.Depth)

	reply2, err := svc.CreateComment(ctx, dto.CreateCommentRequest{
		TargetType:      "manga",
		TargetID:        targetID,
		AuthorUserID:    uuid.NewString(),
		ParentCommentID: &reply1.CommentID,
		Content:         "Reply 2",
	})
	require.NoError(t, err)
	require.Equal(t, 2, reply2.Depth)

	reply3, err := svc.CreateComment(ctx, dto.CreateCommentRequest{
		TargetType:      "manga",
		TargetID:        targetID,
		AuthorUserID:    uuid.NewString(),
		ParentCommentID: &reply2.CommentID,
		Content:         "Reply 3",
	})
	require.NoError(t, err)
	require.Equal(t, 3, reply3.Depth)

	_, err = svc.CreateComment(ctx, dto.CreateCommentRequest{
		TargetType:      "manga",
		TargetID:        targetID,
		AuthorUserID:    uuid.NewString(),
		ParentCommentID: &reply3.CommentID,
		Content:         "Reply 4",
	})
	require.ErrorIs(t, err, ErrReplyDepthExceeded)

	listing, err := svc.ListCommentsByTarget(ctx, dto.ListCommentsRequest{
		TargetType: "manga",
		TargetID:   targetID,
		SortBy:     "newest",
	})
	require.NoError(t, err)
	require.Equal(t, 2, listing.Count)

	thread, err := svc.GetCommentThread(ctx, dto.GetCommentThreadRequest{CommentID: reply2.CommentID})
	require.NoError(t, err)
	require.Equal(t, rootA.CommentID, thread.Root.CommentID)
	require.Equal(t, 3, thread.Count)
}

func TestEditDeleteRestoreModerationVisibilityFlow(t *testing.T) {
	current := time.Date(2026, 3, 12, 20, 30, 0, 0, time.UTC)
	svc := newTestService(&current)
	ctx := context.Background()
	targetID := uuid.NewString()
	authorID := uuid.NewString()

	created, err := svc.CreateComment(ctx, dto.CreateCommentRequest{
		TargetType:   "chapter",
		TargetID:     targetID,
		AuthorUserID: authorID,
		Content:      "Initial",
	})
	require.NoError(t, err)

	newContent := "Edited <content>"
	attachments := []string{"https://cdn.example.com/asset-1.png"}
	spoiler := true
	_, err = svc.UpdateComment(ctx, dto.UpdateCommentRequest{
		CommentID:   created.CommentID,
		ActorUserID: authorID,
		Content:     &newContent,
		Attachments: &attachments,
		Spoiler:     &spoiler,
	})
	require.NoError(t, err)

	current = current.Add(31 * time.Minute)
	secondContent := "Edit after window"
	_, err = svc.UpdateComment(ctx, dto.UpdateCommentRequest{
		CommentID:   created.CommentID,
		ActorUserID: authorID,
		Content:     &secondContent,
	})
	require.ErrorIs(t, err, ErrEditWindowExpired)

	statusHidden := "hidden"
	_, err = svc.UpdateModeration(ctx, dto.UpdateModerationRequest{
		CommentID:        created.CommentID,
		ModerationStatus: &statusHidden,
	})
	require.NoError(t, err)

	_, err = svc.GetCommentDetail(ctx, dto.GetCommentDetailRequest{CommentID: created.CommentID})
	require.ErrorIs(t, err, ErrCommentNotVisible)

	detailHidden, err := svc.GetCommentDetail(ctx, dto.GetCommentDetailRequest{CommentID: created.CommentID, IncludeHidden: true})
	require.NoError(t, err)
	require.Equal(t, "hidden", detailHidden.ModerationStatus)

	_, err = svc.DeleteComment(ctx, dto.DeleteCommentRequest{CommentID: created.CommentID, ActorUserID: uuid.NewString()})
	require.ErrorIs(t, err, ErrForbiddenAction)

	_, err = svc.DeleteComment(ctx, dto.DeleteCommentRequest{CommentID: created.CommentID, ActorUserID: authorID, Reason: "cleanup"})
	require.NoError(t, err)

	detailDeleted, err := svc.GetCommentDetail(ctx, dto.GetCommentDetailRequest{CommentID: created.CommentID, IncludeHidden: true})
	require.NoError(t, err)
	require.True(t, detailDeleted.Deleted)
	require.Equal(t, "[deleted]", detailDeleted.Content)

	current = current.Add(80 * time.Hour)
	_, err = svc.RestoreComment(ctx, dto.RestoreCommentRequest{CommentID: created.CommentID, ActorUserID: authorID})
	require.ErrorIs(t, err, ErrRestoreWindowExpired)

	current = current.Add(10 * time.Second)
	freshAuthor := uuid.NewString()
	fresh, err := svc.CreateComment(ctx, dto.CreateCommentRequest{
		TargetType:   "chapter",
		TargetID:     targetID,
		AuthorUserID: freshAuthor,
		Content:      "Fresh",
	})
	require.NoError(t, err)

	_, err = svc.DeleteComment(ctx, dto.DeleteCommentRequest{CommentID: fresh.CommentID, ActorUserID: freshAuthor})
	require.NoError(t, err)
	_, err = svc.RestoreComment(ctx, dto.RestoreCommentRequest{CommentID: fresh.CommentID, ActorUserID: freshAuthor})
	require.NoError(t, err)
}

func TestContractSurfaceAndLockedReplyFlow(t *testing.T) {
	current := time.Date(2026, 3, 12, 21, 0, 0, 0, time.UTC)
	svc := newTestService(&current)
	ctx := context.Background()
	targetID := uuid.NewString()

	root, err := svc.CreateComment(ctx, dto.CreateCommentRequest{
		TargetType:   "manga",
		TargetID:     targetID,
		AuthorUserID: uuid.NewString(),
		Content:      "Root comment",
	})
	require.NoError(t, err)

	locked := true
	_, err = svc.UpdateModeration(ctx, dto.UpdateModerationRequest{
		CommentID: root.CommentID,
		Locked:    &locked,
	})
	require.NoError(t, err)

	_, err = svc.CreateComment(ctx, dto.CreateCommentRequest{
		TargetType:      "manga",
		TargetID:        targetID,
		AuthorUserID:    uuid.NewString(),
		ParentCommentID: &root.CommentID,
		Content:         "Reply should fail",
	})
	require.ErrorIs(t, err, ErrCommentLocked)

	relation, err := svc.GetTargetRelation(ctx, root.CommentID)
	require.NoError(t, err)
	require.Equal(t, root.CommentID, relation.CommentID)
	require.Equal(t, "manga", relation.TargetType)
	require.Equal(t, targetID, relation.TargetID)

	signal := svc.BuildModerationSignal(root.CommentID, "manga", targetID, "", "req-1", "corr-1")
	require.Equal(t, commentcontract.EventCommentModerated, signal.Event)
	require.Equal(t, "req-1", signal.RequestID)
	require.Equal(t, "corr-1", signal.CorrelationID)
}
