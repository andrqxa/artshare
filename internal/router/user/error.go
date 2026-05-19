package user

import (
	routererror "github.com/andrqxa/artshare/internal/router/error"
)

type (
	FieldError    = routererror.FieldError
	ErrorResponse = routererror.ErrorResponse
)

const (
	ErrCodeUserConflict = 100500 + iota
	ErrCodeUserUnauthorized
)

var (
	ErrUserAlreadyExists = ErrorResponse{
		Code:    ErrCodeUserConflict,
		Message: "A user with this email already exists.",
	}
	ErrUserUnauthorized = ErrorResponse{
		Code:    ErrCodeUserUnauthorized,
		Message: "Invalid email or password.",
	}
)
