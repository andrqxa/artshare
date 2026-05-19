package user

import (
	"fmt"

	"github.com/andrqxa/artshare/internal/controller/apperror"
)

var (
	ErrUserAlreadyExists = fmt.Errorf("%w: user already exists", apperror.ErrConflict)
	ErrInvalidCredentials = fmt.Errorf("%w: invalid credentials", apperror.ErrUnauthorized)
)
