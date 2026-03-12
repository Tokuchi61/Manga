package repository

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/entity"
)

// MemoryStore is a stage-9 bootstrap persistence for comment flows.
type MemoryStore struct {
	mu sync.RWMutex

	commentsByID map[string]entity.Comment
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		commentsByID: make(map[string]entity.Comment),
	}
}

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func cloneTimePtr(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	v := value.UTC()
	return &v
}

func cloneStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	v := *value
	return &v
}

func cloneStringSlice(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	result := make([]string, len(values))
	copy(result, values)
	return result
}

func cloneComment(in entity.Comment) entity.Comment {
	out := in
	out.ParentCommentID = cloneStringPtr(in.ParentCommentID)
	out.RootCommentID = cloneStringPtr(in.RootCommentID)
	out.Attachments = cloneStringSlice(in.Attachments)
	out.EditedAt = cloneTimePtr(in.EditedAt)
	out.DeletedAt = cloneTimePtr(in.DeletedAt)
	return out
}

func commentPopularityScore(comment entity.Comment) int64 {
	riskPenalty := int64(comment.SpamRiskScore)
	if riskPenalty < 0 {
		riskPenalty = 0
	}
	return comment.LikeCount + comment.ReplyCount - riskPenalty
}

func sortComments(items []entity.Comment, sortBy string) {
	sortKey := normalizeValue(sortBy)
	switch sortKey {
	case "oldest":
		sort.Slice(items, func(i, j int) bool {
			if items[i].CreatedAt.Equal(items[j].CreatedAt) {
				return items[i].ID < items[j].ID
			}
			return items[i].CreatedAt.Before(items[j].CreatedAt)
		})
	case "popular":
		sort.Slice(items, func(i, j int) bool {
			leftScore := commentPopularityScore(items[i])
			rightScore := commentPopularityScore(items[j])
			if leftScore == rightScore {
				if items[i].CreatedAt.Equal(items[j].CreatedAt) {
					return items[i].ID < items[j].ID
				}
				return items[i].CreatedAt.After(items[j].CreatedAt)
			}
			return leftScore > rightScore
		})
	default:
		sort.Slice(items, func(i, j int) bool {
			if items[i].CreatedAt.Equal(items[j].CreatedAt) {
				return items[i].ID < items[j].ID
			}
			return items[i].CreatedAt.After(items[j].CreatedAt)
		})
	}
}

func applyOffsetLimit(items []entity.Comment, offset int, limit int) []entity.Comment {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(items) {
		return []entity.Comment{}
	}
	if limit <= 0 {
		limit = 50
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return append([]entity.Comment(nil), items[offset:end]...)
}

func commentVisible(comment entity.Comment, includeHidden bool, includeDeleted bool) bool {
	if !includeDeleted && comment.DeletedAt != nil {
		return false
	}
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
