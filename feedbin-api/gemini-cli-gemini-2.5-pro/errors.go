package feedbin

import (
	"fmt"
	"net/http"
)

// APIError represents an error returned from the Feedbin API.
type APIError struct {
	Response *http.Response
	Message  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %s (status code: %d, url: %s)", e.Message, e.Response.StatusCode, e.Response.Request.URL)
}
