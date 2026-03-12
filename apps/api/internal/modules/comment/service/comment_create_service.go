package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/entity"
	commentrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/repository"
	"github.com/google/uuid"
)

func (s *CommentService) CreateComment(ctx context.Context, request dto.CreateCommentRequest) (dto.CreateCommentResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CreateCommentResponse{}, err
	}

	targetType, err := parseTargetType(request.TargetType)
	if err != nil {
		return dto.CreateCommentResponse{}, err
	}
	targetID, err := parseID(request.TargetID, "target_id")
	if err != nil {
		return dto.CreateCommentResponse{}, err
	}
	if err := s.ensureTargetExists(ctx, targetType, targetID); err != nil {
		return dto.CreateCommentResponse{}, err
	}
	authorUserID, err := parseID(request.AuthorUserID, "author_user_id")
	if err != nil {
		return dto.CreateCommentResponse{}, err
	}
	content, err := sanitizeContent(request.Content)
	if err != nil {
		return dto.CreateCommentResponse{}, err
	}
	attachments := normalizeAttachments(request.Attachments)

	now := s.now().UTC()
	latestByAuthor, err := s.store.GetLatestCommentByAuthor(ctx, authorUserID)
	if err == nil {
		if now.Sub(latestByAuthor.CreatedAt) < s.writeCooldown {
			return dto.CreateCommentResponse{}, ErrRateLimited
		}
	} else if !errors.Is(err, commentrepository.ErrNotFound) {
		return dto.CreateCommentResponse{}, err
	}

	var parentComment *entity.Comment
	var parentCommentID *string
	var rootCommentID *string
	depth := 0

	if request.ParentCommentID != nil {
		parsedParentID, parseErr := parseID(*request.ParentCommentID, "parent_comment_id")
		if parseErr != nil {
			return dto.CreateCommentResponse{}, parseErr
		}
		parent, parentErr := s.store.GetCommentByID(ctx, parsedParentID)
		if parentErr != nil {
			if errors.Is(parentErr, commentrepository.ErrNotFound) {
				return dto.CreateCommentResponse{}, ErrCommentNotFound
			}
			return dto.CreateCommentResponse{}, parentErr
		}
		if parent.TargetType != targetType || parent.TargetID != targetID {
			return dto.CreateCommentResponse{}, fmt.Errorf("%w: parent comment target mismatch", ErrValidation)
		}
		if parent.Locked {
			return dto.CreateCommentResponse{}, ErrCommentLocked
		}
		if parent.DeletedAt != nil {
			return dto.CreateCommentResponse{}, ErrInvalidStateTransition
		}

		depth = parent.Depth + 1
		if depth > s.maxReplyDepth {
			return dto.CreateCommentResponse{}, ErrReplyDepthExceeded
		}

		parentComment = &parent
		parentCommentID = &parsedParentID
		resolvedRootID := parent.ID
		if parent.RootCommentID != nil {
			resolvedRootID = *parent.RootCommentID
		}
		rootCommentID = &resolvedRootID
	}

	comment := entity.Comment{
		ID:               uuid.NewString(),
		TargetType:       targetType,
		TargetID:         targetID,
		AuthorUserID:     authorUserID,
		ParentCommentID:  parentCommentID,
		RootCommentID:    rootCommentID,
		Depth:            depth,
		Content:          request.Content,
		SanitizedContent: content,
		Attachments:      attachments,
		Spoiler:          request.Spoiler,
		ModerationStatus: entity.ModerationStatusVisible,
		SpamRiskScore:    spamRiskScore(content, attachments),
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.store.CreateComment(ctx, comment); err != nil {
		if errors.Is(err, commentrepository.ErrConflict) {
			return dto.CreateCommentResponse{}, ErrCommentAlreadyExists
		}
		return dto.CreateCommentResponse{}, err
	}

	if parentComment != nil {
		parentComment.ReplyCount++
		parentComment.UpdatedAt = now
		if err := s.store.UpdateComment(ctx, *parentComment); err != nil {
			return dto.CreateCommentResponse{}, err
		}
	}

	return dto.CreateCommentResponse{
		CommentID:       comment.ID,
		TargetType:      string(comment.TargetType),
		TargetID:        comment.TargetID,
		ParentCommentID: comment.ParentCommentID,
		RootCommentID:   comment.RootCommentID,
		Depth:           comment.Depth,
		ModerationState: string(comment.ModerationStatus),
	}, nil
}

func (s *CommentService) ensureTargetExists(ctx context.Context, targetType entity.TargetType, targetID string) error {
	switch targetType {
	case entity.TargetTypeManga:
		if s.mangaLookup == nil {
			return nil
		}
		exists, err := s.mangaLookup.TargetExists(ctx, targetID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrTargetNotFound
		}
	case entity.TargetTypeChapter:
		if s.chapterLookup == nil {
			return nil
		}
		exists, err := s.chapterLookup.TargetExists(ctx, targetID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrTargetNotFound
		}
	}
	return nil
}
