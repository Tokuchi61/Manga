package validation

import validatorv10 "github.com/go-playground/validator/v10"

// Validator wraps go-playground/validator as the canonical input validator.
type Validator struct {
	engine *validatorv10.Validate
}

func New() *Validator {
	return &Validator{engine: validatorv10.New()}
}

func (v *Validator) Struct(input any) error {
	if v == nil || v.engine == nil {
		return nil
	}
	return v.engine.Struct(input)
}
