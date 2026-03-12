package dto

import "time"

// ListMangaRequest defines listing/search/filter controls for public surfaces.
type ListMangaRequest struct {
	Search         string `json:"-"`
	Genre          string `json:"-"`
	Tag            string `json:"-"`
	Theme          string `json:"-"`
	ContentWarning string `json:"-"`
	SortBy         string `json:"-" validate:"omitempty,oneof=newest title popular updated"`
	Limit          int    `json:"-" validate:"omitempty,min=1,max=100"`
	Offset         int    `json:"-" validate:"omitempty,min=0"`
}

// DiscoveryRequest defines editorial/discovery listing filters.
type DiscoveryRequest struct {
	Mode          string `json:"-" validate:"omitempty,oneof=all recommended featured collection"`
	CollectionKey string `json:"-"`
	SortBy        string `json:"-" validate:"omitempty,oneof=newest popular"`
	Limit         int    `json:"-" validate:"omitempty,min=1,max=100"`
}

// GetMangaDetailRequest resolves a single public detail surface.
type GetMangaDetailRequest struct {
	MangaID string `json:"-" validate:"required,uuid4"`
}

// MangaListItemResponse is the public listing payload.
type MangaListItemResponse struct {
	MangaID          string     `json:"manga_id"`
	Slug             string     `json:"slug"`
	Title            string     `json:"title"`
	ShortSummary     string     `json:"short_summary"`
	CoverImageURL    string     `json:"cover_image_url"`
	Genres           []string   `json:"genres"`
	Tags             []string   `json:"tags"`
	PublishState     string     `json:"publish_state"`
	Visibility       string     `json:"visibility"`
	Featured         bool       `json:"featured"`
	Recommended      bool       `json:"recommended"`
	ChapterCount     int64      `json:"chapter_count"`
	CommentCount     int64      `json:"comment_count"`
	ViewCount        int64      `json:"view_count"`
	ReleaseSchedule  string     `json:"release_schedule"`
	TranslationGroup string     `json:"translation_group"`
	PublishedAt      *time.Time `json:"published_at,omitempty"`
}

// MangaDetailResponse is the public detail payload.
type MangaDetailResponse struct {
	MangaID            string     `json:"manga_id"`
	Slug               string     `json:"slug"`
	Title              string     `json:"title"`
	AlternativeTitles  []string   `json:"alternative_titles"`
	Summary            string     `json:"summary"`
	ShortSummary       string     `json:"short_summary"`
	CoverImageURL      string     `json:"cover_image_url"`
	BannerImageURL     string     `json:"banner_image_url"`
	SEOTitle           string     `json:"seo_title"`
	SEODescription     string     `json:"seo_description"`
	Genres             []string   `json:"genres"`
	Tags               []string   `json:"tags"`
	Themes             []string   `json:"themes"`
	ContentWarnings    []string   `json:"content_warnings"`
	PublishState       string     `json:"publish_state"`
	Visibility         string     `json:"visibility"`
	Featured           bool       `json:"featured"`
	Recommended        bool       `json:"recommended"`
	CollectionKeys     []string   `json:"collection_keys"`
	DefaultReadAccess  string     `json:"default_read_access"`
	EarlyAccessEnabled bool       `json:"early_access_enabled"`
	EarlyAccessLevel   string     `json:"early_access_level"`
	ReleaseSchedule    string     `json:"release_schedule"`
	TranslationGroup   string     `json:"translation_group"`
	ChapterCount       int64      `json:"chapter_count"`
	CommentCount       int64      `json:"comment_count"`
	ViewCount          int64      `json:"view_count"`
	ContentVersion     int        `json:"content_version"`
	ScheduledAt        *time.Time `json:"scheduled_at,omitempty"`
	PublishedAt        *time.Time `json:"published_at,omitempty"`
	ArchivedAt         *time.Time `json:"archived_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// ListMangaResponse wraps listing results.
type ListMangaResponse struct {
	Items []MangaListItemResponse `json:"items"`
	Count int                     `json:"count"`
}

// DiscoveryResponse wraps discovery results.
type DiscoveryResponse struct {
	Mode  string                  `json:"mode"`
	Items []MangaListItemResponse `json:"items"`
	Count int                     `json:"count"`
}
