package providers

import (
	"fmt"
)

const (
	ErrProviderNotFound = "ERR_PROVIDER_NOT_FOUND"
	ErrInvalidConfig    = "ERR_INVALID_CONFIG"
	ErrAPIError         = "ERR_API_ERROR"
	ErrNoResponse       = "ERR_NO_RESPONSE"
	ErrInvalidRequest   = "ERR_INVALID_REQUEST"
	ErrParseError       = "ERR_PARSE_ERROR"
)

type ProviderError struct {
	Code    string
	Message string
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewError(code, message string) *ProviderError {
	return &ProviderError{
		Code:    code,
		Message: message,
	}
}
