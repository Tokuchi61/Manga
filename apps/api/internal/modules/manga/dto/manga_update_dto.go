package dto

// UpdateMangaMetadataRequest updates owner metadata and taxonomy fields.
type UpdateMangaMetadataRequest struct {
	MangaID           string    `json:"-" validate:"required,uuid4"`
	Title             *string   `json:"title,omitempty" validate:"omitempty,min=1,max=180"`
	Slug              *string   `json:"slug,omitempty" validate:"omitempty,max=180"`
	AlternativeTitles *[]string `json:"alternative_titles,omitempty" validate:"omitempty,max=10,dive,min=1,max=180"`
	Summary           *string   `json:"summary,omitempty" validate:"omitempty,max=5000"`
	ShortSummary      *string   `json:"short_summary,omitempty" validate:"omitempty,max=500"`
	CoverImageURL     *string   `json:"cover_image_url,omitempty" validate:"omitempty,url,max=1000"`
	BannerImageURL    *string   `json:"banner_image_url,omitempty" validate:"omitempty,url,max=1000"`
	SEOTitle          *string   `json:"seo_title,omitempty" validate:"omitempty,max=180"`
	SEODescription    *string   `json:"seo_description,omitempty" validate:"omitempty,max=500"`
	Genres            *[]string `json:"genres,omitempty" validate:"omitempty,max=10,dive,min=1,max=64"`
	Tags              *[]string `json:"tags,omitempty" validate:"omitempty,max=20,dive,min=1,max=64"`
	Themes            *[]string `json:"themes,omitempty" validate:"omitempty,max=10,dive,min=1,max=64"`
	ContentWarnings   *[]string `json:"content_warnings,omitempty" validate:"omitempty,max=10,dive,min=1,max=64"`
	ReleaseSchedule   *string   `json:"release_schedule,omitempty" validate:"omitempty,max=128"`
	TranslationGroup  *string   `json:"translation_group,omitempty" validate:"omitempty,max=128"`
}

// UpdatePublishStateRequest updates manga publication lifecycle.
type UpdatePublishStateRequest struct {
	MangaID     string  `json:"-" validate:"required,uuid4"`
	Action      string  `json:"action" validate:"required,oneof=draft schedule publish archive unpublish"`
	ScheduledAt *string `json:"scheduled_at,omitempty"`
}

// UpdateVisibilityRequest updates public visibility state.
type UpdateVisibilityRequest struct {
	MangaID    string `json:"-" validate:"required,uuid4"`
	Visibility string `json:"visibility" validate:"required,oneof=public hidden"`
}

// UpdateEditorialRequest updates featured/recommended and collection marks.
type UpdateEditorialRequest struct {
	MangaID        string    `json:"-" validate:"required,uuid4"`
	Featured       *bool     `json:"featured,omitempty"`
	Recommended    *bool     `json:"recommended,omitempty"`
	CollectionKeys *[]string `json:"collection_keys,omitempty" validate:"omitempty,max=20,dive,min=1,max=64"`
}

// SyncCountersRequest applies denormalized counter snapshots.
type SyncCountersRequest struct {
	MangaID       string `json:"-" validate:"required,uuid4"`
	ChapterCount  *int64 `json:"chapter_count,omitempty"`
	CommentCount  *int64 `json:"comment_count,omitempty"`
	ViewCount     *int64 `json:"view_count,omitempty"`
	SignalVersion *int   `json:"signal_version,omitempty"`
}

// SoftDeleteMangaRequest performs soft delete.
type SoftDeleteMangaRequest struct {
	MangaID string `json:"-" validate:"required,uuid4"`
}

// RestoreMangaRequest restores soft-deleted manga.
type RestoreMangaRequest struct {
	MangaID string `json:"-" validate:"required,uuid4"`
}

// OperationResponse is a simple lifecycle operation response.
type OperationResponse struct {
	Status string `json:"status"`
}
