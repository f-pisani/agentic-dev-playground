package feedbin

import (
	"fmt"
	"io"
	"net/http"
)

// APIError represents an error returned by the Feedbin API
type APIError struct {
	StatusCode int
	Status     string
	Message    string
	Body       string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("feedbin api error %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("feedbin api error %d: %s", e.StatusCode, e.Status)
}

// IsNotFound returns true if the error is a 404 Not Found
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsUnauthorized returns true if the error is a 401 Unauthorized
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// IsForbidden returns true if the error is a 403 Forbidden
func (e *APIError) IsForbidden() bool {
	return e.StatusCode == http.StatusForbidden
}

// IsUnsupportedMediaType returns true if the error is a 415 Unsupported Media Type
func (e *APIError) IsUnsupportedMediaType() bool {
	return e.StatusCode == http.StatusUnsupportedMediaType
}

// IsMultipleChoices returns true if the error is a 300 Multiple Choices
func (e *APIError) IsMultipleChoices() bool {
	return e.StatusCode == http.StatusMultipleChoices
}

// NewAPIError creates a new API error from an HTTP response
func NewAPIError(resp *http.Response) *APIError {
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	err := &APIError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       string(body),
	}

	// Try to extract a meaningful message based on status code
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		err.Message = "invalid credentials"
	case http.StatusForbidden:
		err.Message = "access denied - you don't own this resource"
	case http.StatusNotFound:
		err.Message = "resource not found"
	case http.StatusUnsupportedMediaType:
		err.Message = "content-type must be application/json; charset=utf-8"
	case http.StatusMultipleChoices:
		err.Message = "multiple feeds found at the specified URL"
	default:
		if len(body) > 0 && len(body) < 200 {
			err.Message = string(body)
		}
	}

	return err
}

// AuthenticationError represents an authentication failure
type AuthenticationError struct {
	Underlying error
}

func (e *AuthenticationError) Error() string {
	if e.Underlying != nil {
		return fmt.Sprintf("authentication failed: %v", e.Underlying)
	}
	return "authentication failed"
}

func (e *AuthenticationError) Unwrap() error {
	return e.Underlying
}

// ValidationError represents a client-side validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// NetworkError represents a network-related error
type NetworkError struct {
	Operation  string
	URL        string
	Underlying error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error during %s to %s: %v", e.Operation, e.URL, e.Underlying)
}

func (e *NetworkError) Unwrap() error {
	return e.Underlying
}

// checkResponse examines an HTTP response and returns an error if the status indicates failure
func checkResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return NewAPIError(resp)
}
