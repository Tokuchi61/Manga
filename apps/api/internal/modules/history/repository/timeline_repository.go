package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/entity"
)

func (s *MemoryStore) ListTimeline(_ context.Context, query TimelineQuery) ([]entity.TimelineEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID := normalizeValue(query.UserID)
	eventFilter := normalizeValue(query.Event)

	rawItems := s.timelineByUser[userID]
	items := make([]entity.TimelineEvent, 0, len(rawItems))
	for _, event := range rawItems {
		if eventFilter != "" && normalizeValue(event.Event) != eventFilter {
			continue
		}
		items = append(items, cloneTimelineEvent(event))
	}

	sortTimeline(items, query.SortBy)
	return applyOffsetLimitTimeline(items, query.Offset, query.Limit), nil
}
