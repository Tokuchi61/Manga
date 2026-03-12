package service

import (
	"context"
	"errors"

	commentcontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/contract"
	commentrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/repository"
)

// GetTargetRelation exposes comment-owned target relation for support/moderation.
func (s *CommentService) GetTargetRelation(ctx context.Context, commentID string) (commentcontract.TargetRelation, error) {
	parsedID, err := parseID(commentID, "comment_id")
	if err != nil {
		return commentcontract.TargetRelation{}, err
	}

	comment, err := s.store.GetCommentByID(ctx, parsedID)
	if err != nil {
		if errors.Is(err, commentrepository.ErrNotFound) {
			return commentcontract.TargetRelation{}, ErrCommentNotFound
		}
		return commentcontract.TargetRelation{}, err
	}

	return commentcontract.TargetRelation{
		CommentID:        comment.ID,
		TargetType:       string(comment.TargetType),
		TargetID:         comment.TargetID,
		ParentCommentID:  comment.ParentCommentID,
		RootCommentID:    comment.RootCommentID,
		ModerationStatus: string(comment.ModerationStatus),
		Deleted:          comment.DeletedAt != nil,
		UpdatedAt:        comment.UpdatedAt,
	}, nil
}

// TargetExists exposes comment target existence checks for consumer modules.
func (s *CommentService) TargetExists(ctx context.Context, commentID string) (bool, error) {
	parsedID, err := parseID(commentID, "comment_id")
	if err != nil {
		return false, nil
	}

	_, err = s.store.GetCommentByID(ctx, parsedID)
	if err != nil {
		if errors.Is(err, commentrepository.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// BuildModerationSignal creates stable comment signal payload.
func (s *CommentService) BuildModerationSignal(commentID string, targetType string, targetID string, event string, requestID string, correlationID string) commentcontract.ModerationSignal {
	if event == "" {
		event = commentcontract.EventCommentModerated
	}
	return commentcontract.ModerationSignal{
		Event:         event,
		CommentID:     commentID,
		TargetType:    targetType,
		TargetID:      targetID,
		OccurredAt:    s.now().UTC(),
		RequestID:     requestID,
		CorrelationID: correlationID,
	}
}
