package feedbin

import "fmt"

type APIError struct {
	StatusCode int
	Body       string
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error: %d %s - %s", e.StatusCode, e.Message, e.Body)
}