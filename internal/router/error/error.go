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
