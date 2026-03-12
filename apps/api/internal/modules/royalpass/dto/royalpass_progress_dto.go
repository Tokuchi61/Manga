package dto

// IngestRoyalPassProgressRequest ingests royalpass progress delta.
type IngestRoyalPassProgressRequest struct {
	ActorUserID   string `json:"-" validate:"required,uuid4"`
	SeasonID      string `json:"season_id,omitempty" validate:"omitempty,max=64"`
	Delta         int    `json:"delta" validate:"required,min=1,max=1000000"`
	SourceType    string `json:"source_type" validate:"required,max=64"`
	SourceRef     string `json:"source_ref,omitempty" validate:"omitempty,max=128"`
	RequestID     string `json:"request_id,omitempty" validate:"omitempty,max=128"`
	CorrelationID string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
}

// IngestRoyalPassProgressResponse returns progress ingest result payload.
type IngestRoyalPassProgressResponse struct {
	Status   string `json:"status"`
	SeasonID string `json:"season_id"`
	Points   int    `json:"points"`
}
