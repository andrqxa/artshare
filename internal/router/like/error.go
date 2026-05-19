package like

import (
	routererror "github.com/andrqxa/artshare/internal/router/error"
)

type (
	FieldError    = routererror.FieldError
	ErrorResponse = routererror.ErrorResponse
)

const (
	ErrCodeLikeNotFound = 100300 + iota
	ErrCodeLikeConflict
)

var (
	ErrLikeNotFound = ErrorResponse{
		Code:    ErrCodeLikeNotFound,
		Message: "The like was not found.",
	}
	ErrLikeConflict = ErrorResponse{
		Code:    ErrCodeLikeConflict,
		Message: "The artwork is already liked.",
	}
)
