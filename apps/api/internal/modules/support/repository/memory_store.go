package repository

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/entity"
)

// MemoryStore is a stage-10 bootstrap persistence for support flows.
type MemoryStore struct {
	mu sync.RWMutex

	casesByID      map[string]entity.SupportCase
	requestIDIndex map[string]map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		casesByID:      make(map[string]entity.SupportCase),
		requestIDIndex: make(map[string]map[string]string),
	}
}

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func cloneStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	v := *value
	return &v
}

func cloneTimePtr(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	v := value.UTC()
	return &v
}

func cloneTargetTypePtr(value *entity.SupportTargetType) *entity.SupportTargetType {
	if value == nil {
		return nil
	}
	v := *value
	return &v
}

func cloneAttachments(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	result := make([]string, len(values))
	copy(result, values)
	return result
}

func cloneReplies(values []entity.SupportReply) []entity.SupportReply {
	if len(values) == 0 {
		return nil
	}
	result := make([]entity.SupportReply, len(values))
	copy(result, values)
	return result
}

func cloneSupportCase(in entity.SupportCase) entity.SupportCase {
	out := in
	out.TargetType = cloneTargetTypePtr(in.TargetType)
	out.TargetID = cloneStringPtr(in.TargetID)
	out.DuplicateOfSupportID = cloneStringPtr(in.DuplicateOfSupportID)
	out.AssigneeUserID = cloneStringPtr(in.AssigneeUserID)
	out.ReviewedByUserID = cloneStringPtr(in.ReviewedByUserID)
	out.ResolvedAt = cloneTimePtr(in.ResolvedAt)
	out.ClosedAt = cloneTimePtr(in.ClosedAt)
	out.ModerationHandoffRequestedAt = cloneTimePtr(in.ModerationHandoffRequestedAt)
	out.LinkedModerationCaseID = cloneStringPtr(in.LinkedModerationCaseID)
	out.Attachments = cloneAttachments(in.Attachments)
	out.Replies = cloneReplies(in.Replies)
	return out
}

func priorityRank(priority entity.SupportPriority) int {
	switch priority {
	case entity.SupportPriorityUrgent:
		return 4
	case entity.SupportPriorityHigh:
		return 3
	case entity.SupportPriorityNormal:
		return 2
	case entity.SupportPriorityLow:
		return 1
	default:
		return 0
	}
}

func sortSupports(items []entity.SupportCase, sortBy string) {
	sortKey := normalizeValue(sortBy)
	switch sortKey {
	case "oldest":
		sort.Slice(items, func(i, j int) bool {
			if items[i].CreatedAt.Equal(items[j].CreatedAt) {
				return items[i].ID < items[j].ID
			}
			return items[i].CreatedAt.Before(items[j].CreatedAt)
		})
	case "priority":
		sort.Slice(items, func(i, j int) bool {
			leftRank := priorityRank(items[i].Priority)
			rightRank := priorityRank(items[j].Priority)
			if leftRank == rightRank {
				if items[i].CreatedAt.Equal(items[j].CreatedAt) {
					return items[i].ID < items[j].ID
				}
				return items[i].CreatedAt.Before(items[j].CreatedAt)
			}
			return leftRank > rightRank
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

func applyOffsetLimit(items []entity.SupportCase, offset int, limit int) []entity.SupportCase {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(items) {
		return []entity.SupportCase{}
	}
	if limit <= 0 {
		limit = 50
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return append([]entity.SupportCase(nil), items[offset:end]...)
}

func defaultQueueStatus(status entity.SupportStatus) bool {
	switch status {
	case entity.SupportStatusOpen, entity.SupportStatusTriaged, entity.SupportStatusWaitingTeam, entity.SupportStatusWaitingUser:
		return true
	default:
		return false
	}
}

func similarCaseMatch(candidate entity.SupportCase, requesterUserID string, kind entity.SupportKind, targetType *entity.SupportTargetType, targetID *string, reasonCode string, since time.Time) bool {
	if normalizeValue(candidate.RequesterUserID) != normalizeValue(requesterUserID) {
		return false
	}
	if candidate.Kind != kind {
		return false
	}
	if normalizeValue(candidate.ReasonCode) != normalizeValue(reasonCode) {
		return false
	}
	if candidate.CreatedAt.Before(since) {
		return false
	}

	if targetType == nil {
		if candidate.TargetType != nil {
			return false
		}
	} else {
		if candidate.TargetType == nil {
			return false
		}
		if normalizeValue(string(*candidate.TargetType)) != normalizeValue(string(*targetType)) {
			return false
		}
	}

	if targetID == nil {
		if candidate.TargetID != nil {
			return false
		}
	} else {
		if candidate.TargetID == nil {
			return false
		}
		if normalizeValue(*candidate.TargetID) != normalizeValue(*targetID) {
			return false
		}
	}

	return true
}
