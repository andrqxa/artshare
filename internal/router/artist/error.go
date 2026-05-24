package artist

import (
	routererror "github.com/andrqxa/artshare/internal/router/error"
)

const (
	ErrCodeArtistNotFound = 100001 + iota
)

var (
	ErrArtistNotFound = routererror.ErrorResponse{
		Code:    ErrCodeArtistNotFound,
		Message: "The requested artist was not found.",
	}
)
