package surepay

import (
	"errors"

	"github.com/surepay-one/surepay-go-sdk/internal/core"
)

// APIError is returned for all non-2xx HTTP responses from the SurePay API.
type APIError = core.APIError

// API error codes returned in the "code" field of error responses.
const (
	CodeValidationError     = "validation_error"
	CodeUnauthorized        = "unauthorized"
	CodeSignatureInvalid    = "signature_invalid"
	CodePermissionDenied    = "permission_denied"
	CodeIPNotAllowed        = "ip_not_allowed"
	CodeNotFound            = "not_found"
	CodeDuplicateRequest    = "duplicate_request"
	CodeInsufficientBalance = "insufficient_balance"
	CodeInvalidTransition   = "invalid_state_transition"
	CodeRateLimitExceeded   = "rate_limit_exceeded"
	CodeNoChannelAvailable  = "no_channel_available"
	CodeInternalError       = "internal_error"
)

// IsNotFound reports whether err is a 404 not_found API error.
func IsNotFound(err error) bool {
	var e *APIError
	return errors.As(err, &e) && e.Code == CodeNotFound
}

// IsRateLimit reports whether err is a 429 rate_limit_exceeded error.
func IsRateLimit(err error) bool {
	var e *APIError
	return errors.As(err, &e) && e.Code == CodeRateLimitExceeded
}

// IsInsufficientBalance reports whether err is a 422 insufficient_balance error.
func IsInsufficientBalance(err error) bool {
	var e *APIError
	return errors.As(err, &e) && e.Code == CodeInsufficientBalance
}

// IsDuplicate reports whether err is a 409 duplicate_request error.
func IsDuplicate(err error) bool {
	var e *APIError
	return errors.As(err, &e) && e.Code == CodeDuplicateRequest
}
