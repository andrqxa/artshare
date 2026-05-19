package artwork

import (
	routererror "github.com/andrqxa/artshare/internal/router/error"
)

type (
	FieldError    = routererror.FieldError
	ErrorResponse = routererror.ErrorResponse
)

const (
	ErrCodeArtworkNotFound = 100200 + iota
	ErrCodeArtworkForbidden
)

var (
	ErrArtworkNotFound = ErrorResponse{
		Code:    ErrCodeArtworkNotFound,
		Message: "The requested artwork was not found.",
	}
	ErrArtworkForbidden = ErrorResponse{
		Code:    ErrCodeArtworkForbidden,
		Message: "You don't have permission to modify this artwork.",
	}
)
