package feedbin

import (
	"fmt"
	"net/http"
)

// AuthenticationService handles communication with the authentication related
// methods of the Feedbin API.
type AuthenticationService service

// Check verifies the user's credentials.
// It returns nil if credentials are valid (HTTP 200 OK), otherwise returns an error.
// Specifically, it returns *APIError with StatusCode 401 for invalid credentials.
func (s *AuthenticationService) Check() error {
	req, err := s.client.newRequest(http.MethodGet, "authentication.json", nil)
	if err != nil {
		return err
	}

	resp, err := s.client.do(req, nil) // No response body expected on success
	if err != nil {
		if apiErr, ok := err.(*APIError); ok {
			if apiErr.StatusCode == http.StatusUnauthorized {
				return apiErr // Return specific 401 error
			}
		}
		return err // Return other errors
	}

	// Check for unexpected success status codes (should be 200)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil // Credentials are valid
}
