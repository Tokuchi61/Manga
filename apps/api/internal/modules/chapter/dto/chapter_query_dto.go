package dto

import "time"

// ListChapterRequest defines manga-based listing query controls.
type ListChapterRequest struct {
	MangaID          string `json:"-" validate:"required,uuid4"`
	SortBy           string `json:"-" validate:"omitempty,oneof=sequence_newest sequence_oldest published_newest published_oldest"`
	Limit            int    `json:"-" validate:"omitempty,min=1,max=200"`
	Offset           int    `json:"-" validate:"omitempty,min=0"`
	IncludeUnpublished bool `json:"-"`
}

// GetChapterDetailRequest resolves chapter detail.
type GetChapterDetailRequest struct {
	ChapterID string `json:"-" validate:"required,uuid4"`
}

// ReadChapterRequest resolves preview/full read payload.
type ReadChapterRequest struct {
	ChapterID string `json:"-" validate:"required,uuid4"`
	Mode      string `json:"-" validate:"omitempty,oneof=preview full"`
	At        *time.Time `json:"-"`
}

// NavigationRequest resolves previous/next/first/last navigation.
type NavigationRequest struct {
	ChapterID string `json:"-" validate:"required,uuid4"`
}

// ChapterListItemResponse is the chapter listing item payload.
type ChapterListItemResponse struct {
	ChapterID         string     `json:"chapter_id"`
	MangaID           string     `json:"manga_id"`
	Slug              string     `json:"slug"`
	Title             string     `json:"title"`
	SequenceNo        int        `json:"sequence_no"`
	DisplayNumber     string     `json:"display_number"`
	PublishState      string     `json:"publish_state"`
	ReadAccessLevel   string     `json:"read_access_level"`
	VIPOnly           bool       `json:"vip_only"`
	EarlyAccessEnabled bool      `json:"early_access_enabled"`
	PreviewEnabled    bool       `json:"preview_enabled"`
	PreviewPageCount  int        `json:"preview_page_count"`
	PageCount         int        `json:"page_count"`
	PublishedAt       *time.Time `json:"published_at,omitempty"`
}

// ChapterDetailResponse is the chapter detail payload.
type ChapterDetailResponse struct {
	ChapterID                    string     `json:"chapter_id"`
	MangaID                      string     `json:"manga_id"`
	Slug                         string     `json:"slug"`
	Title                        string     `json:"title"`
	Summary                      string     `json:"summary"`
	SequenceNo                   int        `json:"sequence_no"`
	DisplayNumber                string     `json:"display_number"`
	PublishState                 string     `json:"publish_state"`
	ReadAccessLevel              string     `json:"read_access_level"`
	InheritAccessFromManga       bool       `json:"inherit_access_from_manga"`
	VIPOnly                      bool       `json:"vip_only"`
	EarlyAccessEnabled           bool       `json:"early_access_enabled"`
	EarlyAccessLevel             string     `json:"early_access_level"`
	EarlyAccessStartAt           *time.Time `json:"early_access_start_at,omitempty"`
	EarlyAccessEndAt             *time.Time `json:"early_access_end_at,omitempty"`
	EarlyAccessFallbackAccess    string     `json:"early_access_fallback_access"`
	PreviewEnabled               bool       `json:"preview_enabled"`
	PreviewPageCount             int        `json:"preview_page_count"`
	MediaHealthStatus            string     `json:"media_health_status"`
	IntegrityStatus              string     `json:"integrity_status"`
	PageCount                    int        `json:"page_count"`
	ScheduledAt                  *time.Time `json:"scheduled_at,omitempty"`
	PublishedAt                  *time.Time `json:"published_at,omitempty"`
	ArchivedAt                   *time.Time `json:"archived_at,omitempty"`
	CreatedAt                    time.Time  `json:"created_at"`
	UpdatedAt                    time.Time  `json:"updated_at"`
}

// ChapterPageResponse is the read payload page item.
type ChapterPageResponse struct {
	PageNumber int    `json:"page_number"`
	MediaURL   string `json:"media_url"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	LongStrip  bool   `json:"long_strip"`
}

// ReadChapterResponse is the chapter read payload.
type ReadChapterResponse struct {
	ChapterID         string               `json:"chapter_id"`
	Mode              string               `json:"mode"`
	PublishState      string               `json:"publish_state"`
	ReadAccessLevel   string               `json:"read_access_level"`
	VIPOnly           bool                 `json:"vip_only"`
	EarlyAccessEnabled bool                `json:"early_access_enabled"`
	EarlyAccessLevel  string               `json:"early_access_level"`
	EarlyAccessActive bool                 `json:"early_access_active"`
	EarlyAccessFallbackAccess string       `json:"early_access_fallback_access"`
	Pages             []ChapterPageResponse `json:"pages"`
	PageCount         int                  `json:"page_count"`
}

// NavigationResponse is chapter navigation payload.
type NavigationResponse struct {
	CurrentChapterID string  `json:"current_chapter_id"`
	PreviousChapterID *string `json:"previous_chapter_id,omitempty"`
	NextChapterID    *string `json:"next_chapter_id,omitempty"`
	FirstChapterID   *string `json:"first_chapter_id,omitempty"`
	LastChapterID    *string `json:"last_chapter_id,omitempty"`
}

// ListChapterResponse wraps chapter list payload.
type ListChapterResponse struct {
	Items []ChapterListItemResponse `json:"items"`
	Count int                       `json:"count"`
}
