package validator

import "github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"

// ValidateStruct validates request DTOs with canonical validator.
func ValidateStruct(v *validation.Validator, payload any) error {
	if v == nil {
		return nil
	}
	return v.Struct(payload)
}
