// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// APIError represents an error returned by the Feedbin API
type APIError struct {
	// StatusCode is the HTTP status code returned by the API
	StatusCode int

	// Message is the error message
	Message string

	// Response is the HTTP response that caused the error
	Response *http.Response
}

// Error returns the error message
func (e *APIError) Error() string {
	return fmt.Sprintf("feedbin: API error: %d %s", e.StatusCode, e.Message)
}

// NewAPIError creates a new APIError from an HTTP response
func NewAPIError(resp *http.Response) *APIError {
	var message string

	// Try to read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		message = fmt.Sprintf("failed to read error response body: %v", err)
	} else if len(body) > 0 {
		message = string(body)
	} else {
		message = resp.Status
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    message,
		Response:   resp,
	}
}

// IsNotFound returns true if the error is a 404 Not Found error
func IsNotFound(err error) bool {
	apiErr, ok := err.(*APIError)
	return ok && apiErr.StatusCode == http.StatusNotFound
}

// IsUnauthorized returns true if the error is a 401 Unauthorized error
func IsUnauthorized(err error) bool {
	apiErr, ok := err.(*APIError)
	return ok && apiErr.StatusCode == http.StatusUnauthorized
}

// IsForbidden returns true if the error is a 403 Forbidden error
func IsForbidden(err error) bool {
	apiErr, ok := err.(*APIError)
	return ok && apiErr.StatusCode == http.StatusForbidden
}
