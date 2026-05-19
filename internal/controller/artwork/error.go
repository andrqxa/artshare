package artwork

import (
	"fmt"

	"github.com/andrqxa/artshare/internal/controller/apperror"
)

var (
	ErrArtworkNotFound  = fmt.Errorf("%w: artwork not found", apperror.ErrNotFound)
	ErrArtworkForbidden = fmt.Errorf("%w: artwork is owned by another user", apperror.ErrForbidden)
)
