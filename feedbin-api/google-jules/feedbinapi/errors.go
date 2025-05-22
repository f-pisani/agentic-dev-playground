package feedbinapi

import (
	"fmt"
	"net/http"
	// "encoding/json" // Uncomment if parsing JSON error bodies
)

// APIError represents an error returned by the Feedbin API.
type APIError struct {
	Response *http.Response // The HTTP response that caused this error.
	Message  string         // Error message derived from status or API response.
	Body     []byte         // Raw response body, useful for debugging.
	// Example: if Feedbin returns structured errors like {"errors": ["message1", "message2"]}
	// APIMessages []string `json:"errors"`
}

// Error returns a string representation of the APIError.
func (e *APIError) Error() string {
	if e.Response != nil {
		return fmt.Sprintf("API error: status_code=%d, message=%s", e.Response.StatusCode, e.Message)
	}
	return fmt.Sprintf("API error: %s", e.Message)
}

// IsStatus returns true if the APIError has the given HTTP status code.
func (e *APIError) IsStatus(statusCode int) bool {
	return e.Response != nil && e.Response.StatusCode == statusCode
}

// newAPIError creates a new APIError.
// It populates the Message field based on the HTTP status.
// If the API returns a JSON error body (e.g., {"error": "message"}),
// this function could be enhanced to parse it.
func newAPIError(resp *http.Response, body []byte) *APIError {
	// Default message based on status code
	message := fmt.Sprintf("request failed with status %s (code %d)", http.StatusText(resp.StatusCode), resp.StatusCode)

	// Placeholder for parsing structured JSON error from body:
	// var errResp struct { ErrorMsg string `json:"error"` } // Or more complex if needed
	// if err := json.Unmarshal(body, &errResp); err == nil && errResp.ErrorMsg != "" {
	//    message = errResp.ErrorMsg
	// }
	// Or if it's like {"errors": ["msg1", "msg2"]}
	// var errResp struct { Errors []string `json:"errors"` }
	// if err := json.Unmarshal(body, &errResp); err == nil && len(errResp.Errors) > 0 {
	//    message = strings.Join(errResp.Errors, ", ")
	// }

	return &APIError{
		Response: resp,
		Message:  message,
		Body:     body,
	}
}

// Predefined error types for common HTTP status codes.
// These can be used with errors.Is for checking specific error conditions.
// Note: These are general purpose; APIError.IsStatus(code) is more direct for checking response codes.

// ErrBadRequest indicates a 400 Bad Request error.
type ErrBadRequest struct{ *APIError }

func (e ErrBadRequest) Error() string { return e.APIError.Error() }

// ErrUnauthorized indicates a 401 Unauthorized error.
type ErrUnauthorized struct{ *APIError }

func (e ErrUnauthorized) Error() string { return e.APIError.Error() }

// ErrForbidden indicates a 403 Forbidden error.
type ErrForbidden struct{ *APIError }

func (e ErrForbidden) Error() string { return e.APIError.Error() }

// ErrNotFound indicates a 404 Not Found error.
type ErrNotFound struct{ *APIError }

func (e ErrNotFound) Error() string { return e.APIError.Error() }

// ErrUnsupportedMedia indicates a 415 Unsupported Media Type error.
type ErrUnsupportedMedia struct{ *APIError }

func (e ErrUnsupportedMedia) Error() string { return e.APIError.Error() }

// ErrInternalServerError indicates a 500 Internal Server Error.
type ErrInternalServerError struct{ *APIError }

func (e ErrInternalServerError) Error() string { return e.APIError.Error() }

// Specific error creation helpers (optional, could be used in `do` method)
func newErrBadRequest(apiErr *APIError) error   { return ErrBadRequest{apiErr} }
func newErrUnauthorized(apiErr *APIError) error { return ErrUnauthorized{apiErr} }

// ... and so on for other specific errors if needed for type assertion by callers.
