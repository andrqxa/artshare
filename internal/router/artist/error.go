package artist

import (
	routererror "github.com/andrqxa/artshare/internal/router/error"
)

type (
	FieldError    = routererror.FieldError
	ErrorResponse = routererror.ErrorResponse
)

const (
	ErrCodeArtistNotFound = 100100 + iota
	ErrCodeArtistConflict
	ErrCodeArtistDisplayNameInvalid
)

var (
	ErrArtistNotFound = ErrorResponse{
		Code:    ErrCodeArtistNotFound,
		Message: "The requested artist was not found.",
	}
	ErrArtistConflict = ErrorResponse{
		Code:    ErrCodeArtistConflict,
		Message: "An artist profile already exists for this user.",
	}
)

func errDisplayNameInvalid() ErrorResponse {
	return routererror.Validation([]FieldError{{Field: "displayName", Message: "must be between 2 and 80 characters"}})
}
