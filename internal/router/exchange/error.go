package exchange

import (
	routererror "github.com/andrqxa/artshare/internal/router/error"
)

type (
	FieldError    = routererror.FieldError
	ErrorResponse = routererror.ErrorResponse
)

const (
	ErrCodeExchangeNotFound = 100400 + iota
	ErrCodeExchangeConflict
	ErrCodeExchangeForbidden
)

var (
	ErrExchangeNotFound = ErrorResponse{
		Code:    ErrCodeExchangeNotFound,
		Message: "The exchange request was not found.",
	}
	ErrExchangeConflict = ErrorResponse{
		Code:    ErrCodeExchangeConflict,
		Message: "The exchange request is not in a state that allows this transition.",
	}
	ErrExchangeForbidden = ErrorResponse{
		Code:    ErrCodeExchangeForbidden,
		Message: "You don't have permission to modify this exchange request.",
	}
)
