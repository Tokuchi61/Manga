package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/entity"
)

func (s *MemoryStore) CreateComment(_ context.Context, comment entity.Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.commentsByID[comment.ID]; exists {
		return ErrConflict
	}

	stored := cloneComment(comment)
	stored.TargetType = entity.TargetType(normalizeValue(string(comment.TargetType)))
	stored.TargetID = normalizeValue(comment.TargetID)
	stored.AuthorUserID = normalizeValue(comment.AuthorUserID)

	s.commentsByID[stored.ID] = stored
	return nil
}

func (s *MemoryStore) GetCommentByID(_ context.Context, commentID string) (entity.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	comment, ok := s.commentsByID[commentID]
	if !ok {
		return entity.Comment{}, ErrNotFound
	}
	return cloneComment(comment), nil
}

func (s *MemoryStore) GetLatestCommentByAuthor(_ context.Context, authorUserID string) (entity.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	authorKey := normalizeValue(authorUserID)
	var latest entity.Comment
	found := false
	for _, comment := range s.commentsByID {
		if normalizeValue(comment.AuthorUserID) != authorKey {
			continue
		}
		if !found || comment.CreatedAt.After(latest.CreatedAt) {
			latest = comment
			found = true
		}
	}
	if !found {
		return entity.Comment{}, ErrNotFound
	}
	return cloneComment(latest), nil
}

func (s *MemoryStore) ListCommentsByTarget(_ context.Context, query ListQuery) ([]entity.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	targetTypeKey := normalizeValue(query.TargetType)
	targetIDKey := normalizeValue(query.TargetID)

	result := make([]entity.Comment, 0)
	for _, comment := range s.commentsByID {
		if normalizeValue(string(comment.TargetType)) != targetTypeKey {
			continue
		}
		if normalizeValue(comment.TargetID) != targetIDKey {
			continue
		}
		if query.ParentOnly && comment.ParentCommentID != nil {
			continue
		}
		if !commentVisible(comment, query.IncludeHidden, query.IncludeDeleted) {
			continue
		}
		result = append(result, cloneComment(comment))
	}

	sortComments(result, query.SortBy)
	return applyOffsetLimit(result, query.Offset, query.Limit), nil
}

func (s *MemoryStore) ListCommentsByRoot(_ context.Context, query ThreadQuery) ([]entity.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rootIDKey := normalizeValue(query.RootCommentID)
	result := make([]entity.Comment, 0)
	for _, comment := range s.commentsByID {
		if comment.RootCommentID == nil {
			continue
		}
		if normalizeValue(*comment.RootCommentID) != rootIDKey {
			continue
		}
		if !commentVisible(comment, query.IncludeHidden, query.IncludeDeleted) {
			continue
		}
		result = append(result, cloneComment(comment))
	}

	sortComments(result, query.SortBy)
	return applyOffsetLimit(result, query.Offset, query.Limit), nil
}

func (s *MemoryStore) UpdateComment(_ context.Context, comment entity.Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.commentsByID[comment.ID]; !exists {
		return ErrNotFound
	}

	stored := cloneComment(comment)
	stored.TargetType = entity.TargetType(normalizeValue(string(comment.TargetType)))
	stored.TargetID = normalizeValue(comment.TargetID)
	stored.AuthorUserID = normalizeValue(comment.AuthorUserID)

	s.commentsByID[stored.ID] = stored
	return nil
}
