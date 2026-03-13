package dto

import "time"

// UpdateContinueReadingStateRequest updates continue-reading surface state.
type UpdateContinueReadingStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateLibraryStateRequest updates library surface state.
type UpdateLibraryStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateTimelineStateRequest updates timeline surface state.
type UpdateTimelineStateRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateBookmarkWriteStateRequest updates bookmark-write surface state.
type UpdateBookmarkWriteStateRequest struct {
	Enabled bool `json:"enabled"`
}

// RuntimeConfigResponse is admin runtime payload.
type RuntimeConfigResponse struct {
	ContinueReadingEnabled bool      `json:"continue_reading_enabled"`
	LibraryEnabled         bool      `json:"library_enabled"`
	TimelineEnabled        bool      `json:"timeline_enabled"`
	BookmarkWriteEnabled   bool      `json:"bookmark_write_enabled"`
	UpdatedAt              time.Time `json:"updated_at"`
}
