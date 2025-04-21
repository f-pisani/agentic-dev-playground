// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"fmt"
	"net/http"
)

// APIError represents an error returned by the Feedbin API.
type APIError struct {
	// StatusCode is the HTTP status code returned by the API.
	StatusCode int

	// Status is the HTTP status text returned by the API.
	Status string

	// Message is the error message.
	Message string

	// Response is the HTTP response that caused this error.
	Response *http.Response
}

// Error returns a string representation of the error.
func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("feedbin: %s: %s", e.Status, e.Message)
	}
	return fmt.Sprintf("feedbin: %s", e.Status)
}

// NewAPIError creates a new APIError from an HTTP response.
func NewAPIError(resp *http.Response) *APIError {
	return &APIError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Response:   resp,
	}
}

// IsNotFound returns true if the error is a 404 Not Found error.
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsUnauthorized returns true if the error is a 401 Unauthorized error.
func IsUnauthorized(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return false
}

// IsForbidden returns true if the error is a 403 Forbidden error.
func IsForbidden(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusForbidden
	}
	return false
}

// IsMultipleChoices returns true if the error is a 300 Multiple Choices error.
func IsMultipleChoices(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusMultipleChoices
	}
	return false
}
