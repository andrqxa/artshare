package exchange

import (
	"fmt"

	"github.com/andrqxa/artshare/internal/controller/apperror"
)

var (
	ErrExchangeNotFound       = fmt.Errorf("%w: exchange request not found", apperror.ErrNotFound)
	ErrExchangeNotPending     = fmt.Errorf("%w: exchange request is not pending", apperror.ErrConflict)
	ErrExchangeNotPermitted   = fmt.Errorf("%w: exchange request is owned by another user", apperror.ErrForbidden)
)
