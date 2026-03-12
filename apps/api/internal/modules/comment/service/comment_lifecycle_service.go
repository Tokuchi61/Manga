package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/dto"
	commentrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/repository"
)

func (s *CommentService) DeleteComment(ctx context.Context, request dto.DeleteCommentRequest) (dto.OperationResponse, error) {
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
	if comment.DeletedAt != nil {
		return dto.OperationResponse{Status: "already_deleted"}, nil
	}

	now := s.now().UTC()
	comment.DeletedAt = &now
	comment.DeleteReason = request.Reason
	comment.UpdatedAt = now

	if err := s.store.UpdateComment(ctx, comment); err != nil {
		if errors.Is(err, commentrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrCommentNotFound
		}
		return dto.OperationResponse{}, err
	}

	return dto.OperationResponse{Status: "soft_deleted"}, nil
}

func (s *CommentService) RestoreComment(ctx context.Context, request dto.RestoreCommentRequest) (dto.OperationResponse, error) {
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
	if comment.DeletedAt == nil {
		return dto.OperationResponse{Status: "already_active"}, nil
	}
	if s.now().UTC().Sub(comment.DeletedAt.UTC()) > s.restoreWindow {
		return dto.OperationResponse{}, ErrRestoreWindowExpired
	}

	comment.DeletedAt = nil
	comment.DeleteReason = ""
	comment.UpdatedAt = s.now().UTC()

	if err := s.store.UpdateComment(ctx, comment); err != nil {
		if errors.Is(err, commentrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrCommentNotFound
		}
		return dto.OperationResponse{}, err
	}

	return dto.OperationResponse{Status: "restored"}, nil
}
