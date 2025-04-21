package feedbin

import (
	"net/http"
)

// AuthenticationService handles authentication-related operations
type AuthenticationService struct {
	client *Client
}

// Verify checks if the credentials are valid
// Returns true if credentials are valid, false otherwise
func (s *AuthenticationService) Verify() (bool, error) {
	req, err := s.client.NewRequest(http.MethodGet, "authentication.json", nil)
	if err != nil {
		return false, err
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		// Check if it's an authentication error (401 Unauthorized)
		if apiErr, ok := err.(*Error); ok && apiErr.StatusCode == http.StatusUnauthorized {
			return false, nil
		}
		return false, err
	}

	// If we get here, authentication was successful
	return true, nil
}
