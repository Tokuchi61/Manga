package entity

import "time"

// PublishState defines chapter publication lifecycle states.
type PublishState string

const (
	PublishStateDraft       PublishState = "draft"
	PublishStateScheduled   PublishState = "scheduled"
	PublishStatePublished   PublishState = "published"
	PublishStateArchived    PublishState = "archived"
	PublishStateUnpublished PublishState = "unpublished"
)

// ReadAccessLevel defines chapter-level read baseline.
type ReadAccessLevel string

const (
	ReadAccessGuest         ReadAccessLevel = "guest"
	ReadAccessAuthenticated ReadAccessLevel = "authenticated"
	ReadAccessVIP           ReadAccessLevel = "vip"
)

// EarlyAccessLevel defines chapter early-access audience.
type EarlyAccessLevel string

const (
	EarlyAccessNone EarlyAccessLevel = "none"
	EarlyAccessVIP  EarlyAccessLevel = "vip"
)

// MediaHealthStatus describes page media readiness state.
type MediaHealthStatus string

const (
	MediaHealthHealthy MediaHealthStatus = "healthy"
	MediaHealthDegraded MediaHealthStatus = "degraded"
	MediaHealthBroken   MediaHealthStatus = "broken"
)

// IntegrityStatus describes page checksum validation state.
type IntegrityStatus string

const (
	IntegrityUnknown IntegrityStatus = "unknown"
	IntegrityPassed  IntegrityStatus = "passed"
	IntegrityFailed  IntegrityStatus = "failed"
)

// ChapterPage keeps page/media metadata owned by chapter module.
type ChapterPage struct {
	PageNumber  int
	MediaURL    string
	Width       int
	Height      int
	LongStrip   bool
	Checksum    string
	CDNHealthy  bool
	Missing     bool
	Broken      bool
	UpdatedAt   time.Time
}

// Chapter is the owner aggregate for stage-8 read and release surfaces.
type Chapter struct {
	ID                          string
	MangaID                     string
	Slug                        string
	Title                       string
	Summary                     string
	SequenceNo                  int
	DisplayNumber               string
	PublishState                PublishState
	ReadAccessLevel             ReadAccessLevel
	InheritAccessFromManga      bool
	VIPOnly                     bool
	EarlyAccessEnabled          bool
	EarlyAccessLevel            EarlyAccessLevel
	EarlyAccessStartAt          *time.Time
	EarlyAccessEndAt            *time.Time
	EarlyAccessFallbackAccess   ReadAccessLevel
	PreviewEnabled              bool
	PreviewPageCount            int
	MediaHealthStatus           MediaHealthStatus
	IntegrityStatus             IntegrityStatus
	PageCount                   int
	Pages                       []ChapterPage
	ScheduledAt                 *time.Time
	PublishedAt                 *time.Time
	ArchivedAt                  *time.Time
	DeletedAt                   *time.Time
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
}
