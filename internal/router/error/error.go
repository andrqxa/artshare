package error

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

const (
	ErrCodeBadRequest          = 100000
	ErrCodeValidation          = 100001
	ErrCodeMultipleJSONValues  = 100002
	ErrCodeInternal            = 100003
	ErrCodeNotFound            = 100004
	ErrCodeConflict            = 100005
	ErrCodeUnauthorized        = 100006
	ErrCodeForbidden           = 100007
)

var (
	ErrBadRequest = ErrorResponse{
		Code:    ErrCodeBadRequest,
		Message: "invalid request body",
	}
	ErrMultipleJSONValues = ErrorResponse{
		Code:    ErrCodeMultipleJSONValues,
		Message: "request body must contain a single JSON object",
	}
	ErrInternal = ErrorResponse{
		Code:    ErrCodeInternal,
		Message: "internal server error",
	}
	ErrNotFound = ErrorResponse{
		Code:    ErrCodeNotFound,
		Message: "requested resource was not found",
	}
	ErrConflict = ErrorResponse{
		Code:    ErrCodeConflict,
		Message: "request conflicts with current resource state",
	}
	ErrUnauthorized = ErrorResponse{
		Code:    ErrCodeUnauthorized,
		Message: "authentication failed",
	}
	ErrForbidden = ErrorResponse{
		Code:    ErrCodeForbidden,
		Message: "permission denied",
	}
)

func Validation(details []FieldError) ErrorResponse {
	return ErrorResponse{
		Code:    ErrCodeValidation,
		Message: "request validation failed",
		Details: details,
	}
}
