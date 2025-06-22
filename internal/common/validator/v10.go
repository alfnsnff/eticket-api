package validator

import (
	"eticket-api/config"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Struct(s any) error
}

type PlaygroundValidator struct {
	v *validator.Validate
}

func NewValidator(cfg *config.Config) Validator {
	return &PlaygroundValidator{
		v: validator.New(),
	}
}

func (p *PlaygroundValidator) Struct(s any) error {
	return p.v.Struct(s)
}
