package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/entity"
	commentrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/repository"
)

func (s *CommentService) ListCommentsByTarget(ctx context.Context, request dto.ListCommentsRequest) (dto.ListCommentsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListCommentsResponse{}, err
	}

	targetType, err := parseTargetType(request.TargetType)
	if err != nil {
		return dto.ListCommentsResponse{}, err
	}
	targetID, err := parseID(request.TargetID, "target_id")
	if err != nil {
		return dto.ListCommentsResponse{}, err
	}

	sortBy := parseSortBy(request.SortBy, "newest")
	limit := request.Limit
	if limit <= 0 {
		limit = 50
	}

	comments, err := s.store.ListCommentsByTarget(ctx, commentrepository.ListQuery{
		TargetType:     string(targetType),
		TargetID:       targetID,
		ParentOnly:     true,
		SortBy:         sortBy,
		Limit:          limit,
		Offset:         request.Offset,
		IncludeHidden:  request.IncludeHidden,
		IncludeDeleted: true,
	})
	if err != nil {
		return dto.ListCommentsResponse{}, err
	}

	items := make([]dto.CommentListItemResponse, 0, len(comments))
	for _, comment := range comments {
		items = append(items, toListItem(comment))
	}

	return dto.ListCommentsResponse{Items: items, Count: len(items)}, nil
}

func (s *CommentService) GetCommentDetail(ctx context.Context, request dto.GetCommentDetailRequest) (dto.CommentDetailResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CommentDetailResponse{}, err
	}

	commentID, err := parseID(request.CommentID, "comment_id")
	if err != nil {
		return dto.CommentDetailResponse{}, err
	}

	comment, err := s.store.GetCommentByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, commentrepository.ErrNotFound) {
			return dto.CommentDetailResponse{}, ErrCommentNotFound
		}
		return dto.CommentDetailResponse{}, err
	}
	if !isCommentVisible(comment, request.IncludeHidden) {
		return dto.CommentDetailResponse{}, ErrCommentNotVisible
	}

	return toDetail(comment), nil
}

func (s *CommentService) GetCommentThread(ctx context.Context, request dto.GetCommentThreadRequest) (dto.CommentThreadResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CommentThreadResponse{}, err
	}

	commentID, err := parseID(request.CommentID, "comment_id")
	if err != nil {
		return dto.CommentThreadResponse{}, err
	}

	current, err := s.store.GetCommentByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, commentrepository.ErrNotFound) {
			return dto.CommentThreadResponse{}, ErrCommentNotFound
		}
		return dto.CommentThreadResponse{}, err
	}

	rootID := current.ID
	if current.RootCommentID != nil {
		rootID = *current.RootCommentID
	}
	root, err := s.store.GetCommentByID(ctx, rootID)
	if err != nil {
		if errors.Is(err, commentrepository.ErrNotFound) {
			return dto.CommentThreadResponse{}, ErrCommentNotFound
		}
		return dto.CommentThreadResponse{}, err
	}
	if !isCommentVisible(root, request.IncludeHidden) {
		return dto.CommentThreadResponse{}, ErrCommentNotVisible
	}

	sortBy := parseSortBy(request.SortBy, "oldest")
	limit := request.Limit
	if limit <= 0 {
		limit = 100
	}

	replies, err := s.store.ListCommentsByRoot(ctx, commentrepository.ThreadQuery{
		RootCommentID:  root.ID,
		SortBy:         sortBy,
		Limit:          limit,
		Offset:         request.Offset,
		IncludeHidden:  request.IncludeHidden,
		IncludeDeleted: true,
	})
	if err != nil {
		return dto.CommentThreadResponse{}, err
	}

	replyItems := make([]dto.CommentListItemResponse, 0, len(replies))
	for _, reply := range replies {
		replyItems = append(replyItems, toListItem(reply))
	}

	return dto.CommentThreadResponse{
		Root:    toDetail(root),
		Replies: replyItems,
		Count:   len(replyItems),
	}, nil
}

func isCommentVisible(comment entity.Comment, includeHidden bool) bool {
	if includeHidden {
		return true
	}
	if comment.ModerationStatus == entity.ModerationStatusHidden {
		return false
	}
	if comment.Shadowbanned {
		return false
	}
	return true
}

func toListItem(comment entity.Comment) dto.CommentListItemResponse {
	return dto.CommentListItemResponse{
		CommentID:        comment.ID,
		TargetType:       string(comment.TargetType),
		TargetID:         comment.TargetID,
		AuthorUserID:     comment.AuthorUserID,
		ParentCommentID:  comment.ParentCommentID,
		RootCommentID:    comment.RootCommentID,
		Depth:            comment.Depth,
		Content:          toSafeContent(comment),
		Spoiler:          comment.Spoiler,
		Pinned:           comment.Pinned,
		Locked:           comment.Locked,
		ModerationStatus: string(comment.ModerationStatus),
		Shadowbanned:     comment.Shadowbanned,
		Deleted:          comment.DeletedAt != nil,
		ReplyCount:       comment.ReplyCount,
		LikeCount:        comment.LikeCount,
		EditCount:        comment.EditCount,
		EditedAt:         comment.EditedAt,
		CreatedAt:        comment.CreatedAt,
		UpdatedAt:        comment.UpdatedAt,
	}
}

func toDetail(comment entity.Comment) dto.CommentDetailResponse {
	return dto.CommentDetailResponse{
		CommentID:        comment.ID,
		TargetType:       string(comment.TargetType),
		TargetID:         comment.TargetID,
		AuthorUserID:     comment.AuthorUserID,
		ParentCommentID:  comment.ParentCommentID,
		RootCommentID:    comment.RootCommentID,
		Depth:            comment.Depth,
		Content:          toSafeContent(comment),
		Attachments:      append([]string(nil), comment.Attachments...),
		Spoiler:          comment.Spoiler,
		Pinned:           comment.Pinned,
		Locked:           comment.Locked,
		ModerationStatus: string(comment.ModerationStatus),
		Shadowbanned:     comment.Shadowbanned,
		SpamRiskScore:    comment.SpamRiskScore,
		Deleted:          comment.DeletedAt != nil,
		DeleteReason:     comment.DeleteReason,
		ReplyCount:       comment.ReplyCount,
		LikeCount:        comment.LikeCount,
		EditCount:        comment.EditCount,
		EditedAt:         comment.EditedAt,
		DeletedAt:        comment.DeletedAt,
		CreatedAt:        comment.CreatedAt,
		UpdatedAt:        comment.UpdatedAt,
	}
}
