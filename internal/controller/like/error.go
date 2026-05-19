package like

import (
	"fmt"

	"github.com/andrqxa/artshare/internal/controller/apperror"
)

var (
	ErrLikeNotFound = fmt.Errorf("%w: like not found", apperror.ErrNotFound)
	ErrLikeConflict = fmt.Errorf("%w: artwork is already liked", apperror.ErrConflict)
)
