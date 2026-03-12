package repository

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/entity"
)

// MemoryStore is stage-11 bootstrap persistence for moderation flows.
type MemoryStore struct {
	mu sync.RWMutex

	casesByID      map[string]entity.Case
	sourceRefIndex map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		casesByID:      make(map[string]entity.Case),
		sourceRefIndex: make(map[string]string),
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

func cloneNotes(notes []entity.ModeratorNote) []entity.ModeratorNote {
	if len(notes) == 0 {
		return nil
	}
	result := make([]entity.ModeratorNote, len(notes))
	copy(result, notes)
	return result
}

func cloneActions(actions []entity.CaseAction) []entity.CaseAction {
	if len(actions) == 0 {
		return nil
	}
	result := make([]entity.CaseAction, len(actions))
	copy(result, actions)
	return result
}

func cloneCase(in entity.Case) entity.Case {
	out := in
	out.SourceRefID = cloneStringPtr(in.SourceRefID)
	out.ReporterUserID = cloneStringPtr(in.ReporterUserID)
	out.AssignedModeratorUserID = cloneStringPtr(in.AssignedModeratorUserID)
	out.EscalatedAt = cloneTimePtr(in.EscalatedAt)
	out.LastActionAt = cloneTimePtr(in.LastActionAt)
	out.Notes = cloneNotes(in.Notes)
	out.Actions = cloneActions(in.Actions)
	return out
}

func sourceRefKey(source entity.CaseSource, sourceRefID string) string {
	return normalizeValue(string(source)) + "|" + normalizeValue(sourceRefID)
}

func sortCases(items []entity.Case, sortBy string) {
	sortKey := normalizeValue(sortBy)
	switch sortKey {
	case "oldest":
		sort.Slice(items, func(i, j int) bool {
			if items[i].CreatedAt.Equal(items[j].CreatedAt) {
				return items[i].ID < items[j].ID
			}
			return items[i].CreatedAt.Before(items[j].CreatedAt)
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

func applyOffsetLimit(items []entity.Case, offset int, limit int) []entity.Case {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(items) {
		return []entity.Case{}
	}
	if limit <= 0 {
		limit = 50
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return append([]entity.Case(nil), items[offset:end]...)
}
