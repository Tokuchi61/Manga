package dto

import "time"

// ChapterPageRequest is input page payload for create/update operations.
type ChapterPageRequest struct {
	PageNumber int    `json:"page_number" validate:"required,min=1"`
	MediaURL   string `json:"media_url" validate:"required,url,max=1500"`
	Width      int    `json:"width" validate:"required,min=1,max=10000"`
	Height     int    `json:"height" validate:"required,min=1,max=10000"`
	LongStrip  bool   `json:"long_strip"`
	Checksum   string `json:"checksum,omitempty" validate:"omitempty,max=256"`
}

// CreateChapterRequest creates a new chapter owner record.
type CreateChapterRequest struct {
	MangaID                       string             `json:"manga_id" validate:"required,uuid4"`
	Slug                          string             `json:"slug,omitempty" validate:"omitempty,max=180"`
	Title                         string             `json:"title" validate:"required,min=1,max=180"`
	Summary                       string             `json:"summary,omitempty" validate:"omitempty,max=5000"`
	SequenceNo                    int                `json:"sequence_no" validate:"required,min=1,max=1000000"`
	DisplayNumber                 string             `json:"display_number,omitempty" validate:"omitempty,max=64"`
	ReadAccessLevel               string             `json:"read_access_level,omitempty" validate:"omitempty,oneof=guest authenticated vip"`
	InheritAccessFromManga        bool               `json:"inherit_access_from_manga"`
	VIPOnly                       bool               `json:"vip_only"`
	EarlyAccessEnabled            bool               `json:"early_access_enabled"`
	EarlyAccessLevel              string             `json:"early_access_level,omitempty" validate:"omitempty,oneof=none vip"`
	EarlyAccessStartAt            *time.Time         `json:"early_access_start_at,omitempty"`
	EarlyAccessEndAt              *time.Time         `json:"early_access_end_at,omitempty"`
	EarlyAccessFallbackAccess     string             `json:"early_access_fallback_access,omitempty" validate:"omitempty,oneof=guest authenticated vip"`
	PreviewEnabled                bool               `json:"preview_enabled"`
	PreviewPageCount              int                `json:"preview_page_count,omitempty" validate:"omitempty,min=0,max=300"`
	PublishState                  string             `json:"publish_state,omitempty" validate:"omitempty,oneof=draft scheduled published"`
	ScheduledAt                   *time.Time         `json:"scheduled_at,omitempty"`
	Pages                         []ChapterPageRequest `json:"pages" validate:"required,min=1,max=500,dive"`
}

// CreateChapterResponse returns stable creation fields.
type CreateChapterResponse struct {
	ChapterID     string `json:"chapter_id"`
	MangaID       string `json:"manga_id"`
	SequenceNo    int    `json:"sequence_no"`
	PublishState  string `json:"publish_state"`
	ReadAccessLevel string `json:"read_access_level"`
}
