package elleven

import (
	"encoding/json"
	"fmt"
)

// APIError represents an error returned by the ERPVoalle API.
// It carries the HTTP status code, parsed API messages, and the raw response body.
type APIError struct {
	// StatusCode is the HTTP status code returned by the server.
	StatusCode int
	// Messages are the API-level messages parsed from the response body.
	Messages []APIMessage
	// RawBody is the raw response body for debugging.
	RawBody []byte
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if len(e.Messages) > 0 {
		return fmt.Sprintf("elleven: API error %d: %s", e.StatusCode, e.Messages[0].Message)
	}
	return fmt.Sprintf("elleven: API error %d: %s", e.StatusCode, string(e.RawBody))
}

// IsNotFound returns true if the error is an HTTP 404 Not Found.
func IsNotFound(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.StatusCode == 404
	}
	return false
}

// IsUnauthorized returns true if the error is an HTTP 401 Unauthorized.
func IsUnauthorized(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.StatusCode == 401
	}
	return false
}

// IsForbidden returns true if the error is an HTTP 403 Forbidden.
func IsForbidden(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.StatusCode == 403
	}
	return false
}

// IsRateLimited returns true if the error is an HTTP 429 Too Many Requests.
func IsRateLimited(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.StatusCode == 429
	}
	return false
}

// IsServerError returns true if the error is an HTTP 5xx server error.
func IsServerError(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.StatusCode >= 500
	}
	return false
}

// parseAPIError attempts to parse an API error from a raw response body.
func parseAPIError(statusCode int, body []byte) *APIError {
	apiErr := &APIError{
		StatusCode: statusCode,
		RawBody:    body,
	}

	var envelope struct {
		Success  bool         `json:"success"`
		Messages []APIMessage `json:"messages"`
	}
	if err := json.Unmarshal(body, &envelope); err == nil && len(envelope.Messages) > 0 {
		apiErr.Messages = envelope.Messages
	}

	return apiErr
}
