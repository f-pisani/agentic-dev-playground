package feedbin

import (
	"fmt"
	"io"
	"net/http"
)

// APIError represents an error returned by the Feedbin API
type APIError struct {
	StatusCode int
	Body       string
	Message    string
}

// NewAPIError creates a new APIError from an HTTP response
func NewAPIError(resp *http.Response) *APIError {
	body, _ := io.ReadAll(resp.Body)

	return &APIError{
		StatusCode: resp.StatusCode,
		Body:       string(body),
		Message:    fmt.Sprintf("API error: %s", resp.Status),
	}
}

// Error returns the error message
func (e *APIError) Error() string {
	return e.Message
}
