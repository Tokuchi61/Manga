package entity

import "time"

// PublishState defines manga publication lifecycle states.
type PublishState string

const (
	PublishStateDraft       PublishState = "draft"
	PublishStateScheduled   PublishState = "scheduled"
	PublishStatePublished   PublishState = "published"
	PublishStateArchived    PublishState = "archived"
	PublishStateUnpublished PublishState = "unpublished"
)

// VisibilityState defines public visibility for manga surfaces.
type VisibilityState string

const (
	VisibilityPublic VisibilityState = "public"
	VisibilityHidden VisibilityState = "hidden"
)

// EarlyAccessLevel defines chapter default early-access audience.
type EarlyAccessLevel string

const (
	EarlyAccessNone EarlyAccessLevel = "none"
	EarlyAccessVIP  EarlyAccessLevel = "vip"
)

// ReadAccessLevel defines chapter default read baseline.
type ReadAccessLevel string

const (
	ReadAccessGuest         ReadAccessLevel = "guest"
	ReadAccessAuthenticated ReadAccessLevel = "authenticated"
)

// Manga is the owner aggregate for stage-7 content metadata and discovery surfaces.
type Manga struct {
	ID                 string
	Slug               string
	Title              string
	AlternativeTitles  []string
	Summary            string
	ShortSummary       string
	CoverImageURL      string
	BannerImageURL     string
	SEOTitle           string
	SEODescription     string
	Genres             []string
	Tags               []string
	Themes             []string
	ContentWarnings    []string
	PublishState       PublishState
	Visibility         VisibilityState
	Featured           bool
	Recommended        bool
	CollectionKeys     []string
	DefaultReadAccess  ReadAccessLevel
	EarlyAccessEnabled bool
	EarlyAccessLevel   EarlyAccessLevel
	ReleaseSchedule    string
	TranslationGroup   string
	ViewCount          int64
	CommentCount       int64
	ChapterCount       int64
	ContentVersion     int
	ScheduledAt        *time.Time
	PublishedAt        *time.Time
	ArchivedAt         *time.Time
	DeletedAt          *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
