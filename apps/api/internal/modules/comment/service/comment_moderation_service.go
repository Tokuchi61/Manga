package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/dto"
	commentrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/repository"
)

func (s *CommentService) UpdateModeration(ctx context.Context, request dto.UpdateModerationRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	commentID, err := parseID(request.CommentID, "comment_id")
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

	updated := false
	if request.ModerationStatus != nil {
		status, parseErr := parseModerationStatus(*request.ModerationStatus)
		if parseErr != nil {
			return dto.OperationResponse{}, parseErr
		}
		comment.ModerationStatus = status
		updated = true
	}
	if request.Pinned != nil {
		comment.Pinned = *request.Pinned
		updated = true
	}
	if request.Locked != nil {
		comment.Locked = *request.Locked
		updated = true
	}
	if request.Shadowbanned != nil {
		comment.Shadowbanned = *request.Shadowbanned
		updated = true
	}
	if request.Spoiler != nil {
		comment.Spoiler = *request.Spoiler
		updated = true
	}
	if !updated {
		return dto.OperationResponse{}, errors.Join(ErrValidation, errors.New("no moderation field provided"))
	}

	comment.UpdatedAt = s.now().UTC()
	if err := s.store.UpdateComment(ctx, comment); err != nil {
		if errors.Is(err, commentrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrCommentNotFound
		}
		return dto.OperationResponse{}, err
	}

	return dto.OperationResponse{Status: "moderation_updated"}, nil
}
