package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/dto"
	commentrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/repository"
)

func (s *CommentService) UpdateComment(ctx context.Context, request dto.UpdateCommentRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	commentID, err := parseID(request.CommentID, "comment_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}

	comment, err := s.store.GetCommentByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, commentrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrCommentNotFound
		}
		return dto.OperationResponse{}, err
	}

	if comment.AuthorUserID != actorUserID {
		return dto.OperationResponse{}, ErrForbiddenAction
	}
	if comment.Locked {
		return dto.OperationResponse{}, ErrCommentLocked
	}
	if comment.DeletedAt != nil {
		return dto.OperationResponse{}, ErrInvalidStateTransition
	}
	if s.now().UTC().Sub(comment.CreatedAt) > s.editWindow {
		return dto.OperationResponse{}, ErrEditWindowExpired
	}

	updated := false
	if request.Content != nil {
		sanitized, sanitizeErr := sanitizeContent(*request.Content)
		if sanitizeErr != nil {
			return dto.OperationResponse{}, sanitizeErr
		}
		comment.Content = *request.Content
		comment.SanitizedContent = sanitized
		comment.SpamRiskScore = spamRiskScore(sanitized, comment.Attachments)
		updated = true
	}
	if request.Attachments != nil {
		comment.Attachments = normalizeAttachments(*request.Attachments)
		comment.SpamRiskScore = spamRiskScore(comment.SanitizedContent, comment.Attachments)
		updated = true
	}
	if request.Spoiler != nil {
		comment.Spoiler = *request.Spoiler
		updated = true
	}
	if !updated {
		return dto.OperationResponse{}, errors.Join(ErrValidation, errors.New("no update field provided"))
	}

	now := s.now().UTC()
	comment.EditCount++
	comment.EditedAt = &now
	comment.UpdatedAt = now

	if err := s.store.UpdateComment(ctx, comment); err != nil {
		if errors.Is(err, commentrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrCommentNotFound
		}
		return dto.OperationResponse{}, err
	}

	return dto.OperationResponse{Status: "comment_updated"}, nil
}
