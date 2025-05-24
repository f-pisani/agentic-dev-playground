package feedbin

import (
	"fmt"
	"net/http"
)

// APIError represents an error returned by the Feedbin API
type APIError struct {
	StatusCode int
	Status     string
	Message    string
	Response   *http.Response
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("feedbin API error: %s (status: %d %s)", e.Message, e.StatusCode, e.Status)
	}
	return fmt.Sprintf("feedbin API error: %d %s", e.StatusCode, e.Status)
}

// AuthenticationError represents a 401 Unauthorized error
type AuthenticationError struct {
	*APIError
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("authentication failed: %s", e.APIError.Error())
}

// NotFoundError represents a 404 Not Found error
type NotFoundError struct {
	*APIError
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("resource not found: %s", e.APIError.Error())
}

// ValidationError represents a 400 Bad Request error
type ValidationError struct {
	*APIError
	Errors map[string][]string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.APIError.Error())
}

// RateLimitError represents a rate limit error
type RateLimitError struct {
	*APIError
	RetryAfter int
}

func (e *RateLimitError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("rate limit exceeded, retry after %d seconds: %s", e.RetryAfter, e.APIError.Error())
	}
	return fmt.Sprintf("rate limit exceeded: %s", e.APIError.Error())
}

// ServerError represents a 5xx server error
type ServerError struct {
	*APIError
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("server error: %s", e.APIError.Error())
}

// newAPIError creates an appropriate error type based on the HTTP status code
func newAPIError(resp *http.Response, message string) error {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Message:    message,
		Response:   resp,
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return &AuthenticationError{APIError: apiErr}
	case http.StatusNotFound:
		return &NotFoundError{APIError: apiErr}
	case http.StatusBadRequest:
		return &ValidationError{APIError: apiErr}
	case http.StatusTooManyRequests:
		retryAfter := 0
		if val := resp.Header.Get("Retry-After"); val != "" {
			fmt.Sscanf(val, "%d", &retryAfter)
		}
		return &RateLimitError{APIError: apiErr, RetryAfter: retryAfter}
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return &ServerError{APIError: apiErr}
	default:
		return apiErr
	}
}
