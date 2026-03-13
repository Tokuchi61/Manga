package dto

import "time"

// ListContinueReadingRequest resolves own continue-reading surface.
type ListContinueReadingRequest struct {
	UserID string `json:"-" validate:"required,uuid4"`
	Limit  int    `json:"-" validate:"omitempty,min=1,max=200"`
	Offset int    `json:"-" validate:"omitempty,min=0"`
	SortBy string `json:"-" validate:"omitempty,oneof=newest oldest"`
}

// ListTimelineRequest resolves own history timeline surface.
type ListTimelineRequest struct {
	UserID string `json:"-" validate:"required,uuid4"`
	Event  string `json:"-" validate:"omitempty,max=64"`
	Limit  int    `json:"-" validate:"omitempty,min=1,max=200"`
	Offset int    `json:"-" validate:"omitempty,min=0"`
	SortBy string `json:"-" validate:"omitempty,oneof=newest oldest"`
}

// HistoryEntryResponse is library/continue-reading payload item.
type HistoryEntryResponse struct {
	LibraryEntryID string    `json:"library_entry_id"`
	UserID         string    `json:"user_id,omitempty"`
	MangaID        string    `json:"manga_id"`
	LastChapterID  string    `json:"last_chapter_id,omitempty"`
	LastPageNumber int       `json:"last_page_number"`
	PageCount      int       `json:"page_count"`
	Status         string    `json:"status"`
	Bookmarked     bool      `json:"bookmarked"`
	Favorited      bool      `json:"favorited"`
	SharePublic    bool      `json:"share_public"`
	LastReadAt     time.Time `json:"last_read_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ListContinueReadingResponse wraps continue-reading payload.
type ListContinueReadingResponse struct {
	Items []HistoryEntryResponse `json:"items"`
	Count int                    `json:"count"`
}

// TimelineItemResponse is timeline item payload.
type TimelineItemResponse struct {
	TimelineEventID string    `json:"timeline_event_id"`
	MangaID         string    `json:"manga_id"`
	ChapterID       string    `json:"chapter_id"`
	Event           string    `json:"event"`
	PageNumber      int       `json:"page_number"`
	PageCount       int       `json:"page_count"`
	RequestID       string    `json:"request_id,omitempty"`
	CorrelationID   string    `json:"correlation_id,omitempty"`
	OccurredAt      time.Time `json:"occurred_at"`
}

// ListTimelineResponse wraps timeline payload.
type ListTimelineResponse struct {
	Items []TimelineItemResponse `json:"items"`
	Count int                    `json:"count"`
}

// OperationResponse is generic operation result payload.
type OperationResponse struct {
	Status string `json:"status"`
}
