package dto

// IngestMissionProgressRequest ingests mission progress.
type IngestMissionProgressRequest struct {
	ActorUserID   string `json:"-" validate:"required,uuid4"`
	MissionID     string `json:"-" validate:"required,max=128"`
	Delta         int    `json:"delta" validate:"required,min=1,max=100000"`
	SourceType    string `json:"source_type" validate:"required,max=64"`
	SourceRef     string `json:"source_ref,omitempty" validate:"omitempty,max=128"`
	RequestID     string `json:"request_id,omitempty" validate:"omitempty,max=128"`
	CorrelationID string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
}

// IngestMissionProgressResponse returns mission progress ingest result.
type IngestMissionProgressResponse struct {
	Status    string                      `json:"status"`
	Completed bool                        `json:"completed"`
	Mission   MissionProgressItemResponse `json:"mission"`
}
