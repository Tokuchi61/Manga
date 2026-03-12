package dto

import "time"

// CreateMangaRequest creates a new manga owner record.
type CreateMangaRequest struct {
	Title                     string     `json:"title" validate:"required,min=1,max=180"`
	Slug                      string     `json:"slug,omitempty" validate:"omitempty,max=180"`
	AlternativeTitles         []string   `json:"alternative_titles,omitempty" validate:"omitempty,max=10,dive,min=1,max=180"`
	Summary                   string     `json:"summary" validate:"required,max=5000"`
	ShortSummary              string     `json:"short_summary,omitempty" validate:"omitempty,max=500"`
	CoverImageURL             string     `json:"cover_image_url,omitempty" validate:"omitempty,url,max=1000"`
	BannerImageURL            string     `json:"banner_image_url,omitempty" validate:"omitempty,url,max=1000"`
	SEOTitle                  string     `json:"seo_title,omitempty" validate:"omitempty,max=180"`
	SEODescription            string     `json:"seo_description,omitempty" validate:"omitempty,max=500"`
	Genres                    []string   `json:"genres" validate:"required,min=1,max=10,dive,min=1,max=64"`
	Tags                      []string   `json:"tags,omitempty" validate:"omitempty,max=20,dive,min=1,max=64"`
	Themes                    []string   `json:"themes,omitempty" validate:"omitempty,max=10,dive,min=1,max=64"`
	ContentWarnings           []string   `json:"content_warnings,omitempty" validate:"omitempty,max=10,dive,min=1,max=64"`
	PublishState              string     `json:"publish_state,omitempty" validate:"omitempty,oneof=draft scheduled published"`
	ScheduledAt               *time.Time `json:"scheduled_at,omitempty"`
	DefaultReadAccessLevel    string     `json:"default_read_access_level,omitempty" validate:"omitempty,oneof=guest authenticated"`
	DefaultEarlyAccessEnabled bool       `json:"default_early_access_enabled"`
	DefaultEarlyAccessLevel   string     `json:"default_early_access_level,omitempty" validate:"omitempty,oneof=none vip"`
	ReleaseSchedule           string     `json:"release_schedule,omitempty" validate:"omitempty,max=128"`
	TranslationGroup          string     `json:"translation_group,omitempty" validate:"omitempty,max=128"`
}

// CreateMangaResponse returns stable creation fields.
type CreateMangaResponse struct {
	MangaID      string `json:"manga_id"`
	Slug         string `json:"slug"`
	PublishState string `json:"publish_state"`
	Visibility   string `json:"visibility"`
}
