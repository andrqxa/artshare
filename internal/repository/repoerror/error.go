package repoerror

import "errors"

var (
	ErrConflict = errors.New("resource conflict")
	ErrNotFound = errors.New("resource not found")
)
