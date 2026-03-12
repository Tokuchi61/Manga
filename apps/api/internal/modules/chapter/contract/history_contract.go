package contract

import "time"

const (
	EventReadStarted   = "chapter.read.started"
	EventReadCheckpoint = "chapter.read.checkpoint"
	EventReadFinished  = "chapter.read.finished"
)

// ResumeAnchor is chapter-owned anchor surface consumed by history module.
type ResumeAnchor struct {
	ChapterID   string
	MangaID     string
	PageNumber  int
	PageCount   int
	UpdatedAt   time.Time
}

// ReadSignal is chapter->history progress signal contract.
type ReadSignal struct {
	Event       string
	ChapterID   string
	MangaID     string
	PageNumber  int
	PageCount   int
	OccurredAt  time.Time
	RequestID   string
	CorrelationID string
}
