package validator

import (
	"github.com/go-playground/validator/v10"
)

func ParseErrors(err error) map[string]string {
	errs := make(map[string]string)
	if verrs, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range verrs {
			errs[fe.Field()] = msg(fe)
		}
	}
	return errs
}

func msg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return "must be at least " + fe.Param()
	case "max":
		return "must be at most " + fe.Param()
	case "gte":
		return "must be ≥ " + fe.Param()
	case "lte":
		return "must be ≤ " + fe.Param()
	case "len":
		return "must be exactly " + fe.Param() + " characters"
	case "oneof":
		return "must be one of: " + fe.Param()
	}
	return "is invalid"
}
