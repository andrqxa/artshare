package artist

import (
	"fmt"

	"github.com/andrqxa/artshare/internal/controller/apperror"
)

var (
	ErrArtistNotFound      = fmt.Errorf("%w: artist not found", apperror.ErrNotFound)
	ErrArtistAlreadyExists = fmt.Errorf("%w: artist already exists", apperror.ErrConflict)
)
