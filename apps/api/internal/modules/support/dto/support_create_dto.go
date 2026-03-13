package dto

// CreateCommunicationRequest creates general communication support intake.
type CreateCommunicationRequest struct {
	RequesterUserID string   `json:"-" validate:"required,uuid4"`
	Category        string   `json:"category" validate:"required,min=1,max=64"`
	Priority        string   `json:"priority,omitempty" validate:"omitempty,oneof=low normal high urgent"`
	ReasonCode      string   `json:"reason_code,omitempty" validate:"omitempty,max=64"`
	ReasonText      string   `json:"reason_text" validate:"required,min=1,max=5000"`
	Attachments     []string `json:"attachments,omitempty" validate:"omitempty,max=10,dive,url,max=1500"`
	RequestID       string   `json:"request_id" validate:"required,max=128"`
}

// CreateTicketRequest creates support ticket intake.
type CreateTicketRequest struct {
	RequesterUserID string   `json:"-" validate:"required,uuid4"`
	Category        string   `json:"category" validate:"required,min=1,max=64"`
	Priority        string   `json:"priority,omitempty" validate:"omitempty,oneof=low normal high urgent"`
	ReasonCode      string   `json:"reason_code,omitempty" validate:"omitempty,max=64"`
	ReasonText      string   `json:"reason_text" validate:"required,min=1,max=5000"`
	Attachments     []string `json:"attachments,omitempty" validate:"omitempty,max=10,dive,url,max=1500"`
	RequestID       string   `json:"request_id" validate:"required,max=128"`
}

// CreateReportRequest creates target-bound report intake.
type CreateReportRequest struct {
	RequesterUserID string   `json:"-" validate:"required,uuid4"`
	Category        string   `json:"category" validate:"required,min=1,max=64"`
	Priority        string   `json:"priority,omitempty" validate:"omitempty,oneof=low normal high urgent"`
	ReasonCode      string   `json:"reason_code,omitempty" validate:"omitempty,max=64"`
	ReasonText      string   `json:"reason_text" validate:"required,min=1,max=5000"`
	TargetType      string   `json:"target_type" validate:"required,oneof=manga chapter comment"`
	TargetID        string   `json:"target_id" validate:"required,uuid4"`
	Attachments     []string `json:"attachments,omitempty" validate:"omitempty,max=10,dive,url,max=1500"`
	RequestID       string   `json:"request_id" validate:"required,max=128"`
}

// CreateSupportResponse returns stable create fields.
type CreateSupportResponse struct {
	SupportID            string  `json:"support_id"`
	SupportKind          string  `json:"support_kind"`
	Status               string  `json:"status"`
	DuplicateOfSupportID *string `json:"duplicate_of_support_id,omitempty"`
	TargetType           *string `json:"target_type,omitempty"`
	TargetID             *string `json:"target_id,omitempty"`
}
