package apperror

import "errors"

var (
	ErrConflict     = errors.New("resource conflict")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
)
