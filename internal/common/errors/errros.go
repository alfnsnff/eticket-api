package errors

import "errors"

var (
	ErrNotFound        = errors.New("resource not found")
	ErrConflict        = errors.New("resource already exists or conflict")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrValidation      = errors.New("validation error")
	ErrInternal        = errors.New("internal error")
	ErrExpired         = errors.New("resource expired")
	ErrBadRequest      = errors.New("bad request")
	ErrExternalTimeout = errors.New("external service timeout")
	ErrExternalDown    = errors.New("external service unavailable")
)
