package feedbin

import (
	"net/http"
)

// AuthenticationService handles authentication-related operations
type AuthenticationService struct {
	client *Client
}

// Verify checks if the provided credentials are valid
func (s *AuthenticationService) Verify() (bool, error) {
	req, err := s.client.newRequest(http.MethodGet, "/authentication.json", nil)
	if err != nil {
		return false, err
	}

	resp, err := s.client.do(req, nil)
	if err != nil {
		// Check if the error is due to invalid credentials
		if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode == http.StatusUnauthorized {
			return false, nil
		}
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}
