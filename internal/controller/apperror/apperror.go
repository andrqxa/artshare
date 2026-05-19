package apperror

import (
	"errors"
	"fmt"
)

var (
	ErrConflict     = errors.New("resource conflict")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInternal     = errors.New("internal error")
)

func WrapInternal(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %w", ErrInternal, err)
}
