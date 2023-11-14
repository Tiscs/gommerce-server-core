package validator

import (
	"context"
)

type validateAller interface {
	ValidateAll() error
}

type validator interface {
	Validate(all bool) error
}

type validatorLegacy interface {
	Validate() error
}

func Validate(_ context.Context, data any) error {
	switch v := data.(type) {
	case validateAller:
		return v.ValidateAll()
	case validator:
		return v.Validate(false)
	case validatorLegacy:
		return v.Validate()
	default:
		return nil
	}
}
