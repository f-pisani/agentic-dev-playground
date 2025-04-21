package feedbin

import (
	"net/http"
)

// AuthenticationService handles communication with the authentication related
// endpoints of the Feedbin API
type AuthenticationService struct {
	client *Client
}

// Validate checks if the provided credentials are valid
func (s *AuthenticationService) Validate() (bool, error) {
	req, err := s.client.NewRequest(http.MethodGet, "authentication.json", nil)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		// Check if the error is a 401 Unauthorized
		if errResp, ok := err.(*ErrorResponse); ok && errResp.Response.StatusCode == http.StatusUnauthorized {
			return false, nil
		}
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}
