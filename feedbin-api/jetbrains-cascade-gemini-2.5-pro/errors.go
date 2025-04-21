package feedbin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError represents an error response from the Feedbin API.
type APIError struct {
	StatusCode int            // HTTP status code
	Response   *http.Response // HTTP response that caused this error
	Body       string         // Raw response body
	Message    string         // Human-readable message (may be empty)
}

func (e *APIError) Error() string {
	msg := fmt.Sprintf("HTTP %d", e.StatusCode)
	if e.Message != "" {
		msg += ": " + e.Message
	} else if e.Body != "" {
		msg += ": " + e.Body // Use body if message is empty
	}
	return msg
}

// MultipleChoicesError represents a 300 Multiple Choices response, typically from CreateSubscription.
type MultipleChoicesError struct {
	APIError              // Embed APIError
	Choices  []FeedChoice // Choices returned by the API
}

// FeedChoice represents one choice when multiple feeds are found for a URL.
type FeedChoice struct {
	FeedURL string `json:"feed_url"`
	Title   string `json:"title"`
}

// NewMultipleChoicesError creates a MultipleChoicesError from an APIError.
func NewMultipleChoicesError(apiErr *APIError) *MultipleChoicesError {
	mcErr := &MultipleChoicesError{APIError: *apiErr}
	// Attempt to unmarshal the body into choices
	_ = json.Unmarshal([]byte(apiErr.Body), &mcErr.Choices)
	return mcErr
}

func (e *MultipleChoicesError) Error() string {
	return fmt.Sprintf("Multiple choices found (HTTP %d): %s", e.StatusCode, e.Body)
}
