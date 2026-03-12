package validator

import "github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"

// ValidateStruct validates moderation DTO payloads.
func ValidateStruct(v *validation.Validator, payload any) error {
	if v == nil {
		return nil
	}
	return v.Struct(payload)
}
