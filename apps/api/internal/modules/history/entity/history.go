package entity

import "time"

// ReadingStatus defines library reading lifecycle state.
type ReadingStatus string

const (
	ReadingStatusInProgress ReadingStatus = "in_progress"
	ReadingStatusCompleted  ReadingStatus = "completed"
	ReadingStatusDropped    ReadingStatus = "dropped"
)

// LibraryEntry is user-owned manga reading state.
type LibraryEntry struct {
	ID             string
	UserID         string
	MangaID        string
	LastChapterID  string
	LastPageNumber int
	PageCount      int
	Status         ReadingStatus
	Bookmarked     bool
	Favorited      bool
	SharePublic    bool
	LastReadAt     time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TimelineEvent is immutable history timeline record.
type TimelineEvent struct {
	ID            string
	UserID        string
	MangaID       string
	ChapterID     string
	Event         string
	PageNumber    int
	PageCount     int
	RequestID     string
	CorrelationID string
	OccurredAt    time.Time
	CreatedAt     time.Time
}

// Checkpoint carries chapter->history ingestion data.
type Checkpoint struct {
	UserID        string
	MangaID       string
	ChapterID     string
	Event         string
	PageNumber    int
	PageCount     int
	RequestID     string
	CorrelationID string
	OccurredAt    time.Time
}

// RuntimeConfig stores module-level operational controls.
type RuntimeConfig struct {
	ContinueReadingEnabled bool
	LibraryEnabled         bool
	TimelineEnabled        bool
	BookmarkWriteEnabled   bool
	UpdatedAt              time.Time
}
