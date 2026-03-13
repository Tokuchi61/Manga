package dto

// IngestChapterSignalRequest ingests chapter->history signal.
type IngestChapterSignalRequest struct {
	UserID        string `json:"-" validate:"required,uuid4"`
	ChapterID     string `json:"chapter_id" validate:"required,uuid4"`
	Event         string `json:"event,omitempty" validate:"omitempty,max=64"`
	PageNumber    int    `json:"page_number" validate:"omitempty,min=0,max=10000"`
	RequestID     string `json:"request_id,omitempty" validate:"omitempty,max=128"`
	CorrelationID string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
}

// IngestChapterSignalResponse is intake result payload.
type IngestChapterSignalResponse struct {
	LibraryEntryID  string `json:"library_entry_id"`
	TimelineEventID string `json:"timeline_event_id,omitempty"`
	MangaID         string `json:"manga_id"`
	ChapterID       string `json:"chapter_id"`
	Event           string `json:"event"`
	Status          string `json:"status"`
	Created         bool   `json:"created"`
}
