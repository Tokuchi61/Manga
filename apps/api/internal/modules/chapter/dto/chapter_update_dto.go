package dto

// UpdateChapterRequest updates chapter metadata and page structures.
type UpdateChapterRequest struct {
	ChapterID      string               `json:"-" validate:"required,uuid4"`
	Slug           *string              `json:"slug,omitempty" validate:"omitempty,max=180"`
	Title          *string              `json:"title,omitempty" validate:"omitempty,min=1,max=180"`
	Summary        *string              `json:"summary,omitempty" validate:"omitempty,max=5000"`
	SequenceNo     *int                 `json:"sequence_no,omitempty" validate:"omitempty,min=1,max=1000000"`
	DisplayNumber  *string              `json:"display_number,omitempty" validate:"omitempty,max=64"`
	Pages          *[]ChapterPageRequest `json:"pages,omitempty" validate:"omitempty,min=1,max=500,dive"`
}

// UpdatePublishStateRequest updates chapter publish lifecycle.
type UpdatePublishStateRequest struct {
	ChapterID   string  `json:"-" validate:"required,uuid4"`
	Action      string  `json:"action" validate:"required,oneof=draft schedule publish archive unpublish"`
	ScheduledAt *string `json:"scheduled_at,omitempty"`
}

// UpdateAccessRequest updates chapter read/early access fields.
type UpdateAccessRequest struct {
	ChapterID                  string  `json:"-" validate:"required,uuid4"`
	ReadAccessLevel            *string `json:"read_access_level,omitempty" validate:"omitempty,oneof=guest authenticated vip"`
	InheritAccessFromManga     *bool   `json:"inherit_access_from_manga,omitempty"`
	VIPOnly                    *bool   `json:"vip_only,omitempty"`
	EarlyAccessEnabled         *bool   `json:"early_access_enabled,omitempty"`
	EarlyAccessLevel           *string `json:"early_access_level,omitempty" validate:"omitempty,oneof=none vip"`
	EarlyAccessStartAt         *string `json:"early_access_start_at,omitempty"`
	EarlyAccessEndAt           *string `json:"early_access_end_at,omitempty"`
	EarlyAccessFallbackAccess  *string `json:"early_access_fallback_access,omitempty" validate:"omitempty,oneof=guest authenticated vip"`
	PreviewEnabled             *bool   `json:"preview_enabled,omitempty"`
	PreviewPageCount           *int    `json:"preview_page_count,omitempty" validate:"omitempty,min=0,max=300"`
}

// ReorderChapterRequest updates chapter sequence for navigation.
type ReorderChapterRequest struct {
	ChapterID  string `json:"-" validate:"required,uuid4"`
	SequenceNo int    `json:"sequence_no" validate:"required,min=1,max=1000000"`
}

// UpdateMediaHealthRequest updates media health signal.
type UpdateMediaHealthRequest struct {
	ChapterID          string `json:"-" validate:"required,uuid4"`
	MediaHealthStatus  string `json:"media_health_status" validate:"required,oneof=healthy degraded broken"`
}

// UpdateIntegrityRequest updates chapter integrity signal.
type UpdateIntegrityRequest struct {
	ChapterID       string `json:"-" validate:"required,uuid4"`
	IntegrityStatus string `json:"integrity_status" validate:"required,oneof=unknown passed failed"`
}

// SoftDeleteChapterRequest soft-deletes chapter.
type SoftDeleteChapterRequest struct {
	ChapterID string `json:"-" validate:"required,uuid4"`
}

// RestoreChapterRequest restores soft-deleted chapter.
type RestoreChapterRequest struct {
	ChapterID string `json:"-" validate:"required,uuid4"`
}

// OperationResponse is a simple operation payload.
type OperationResponse struct {
	Status string `json:"status"`
}
