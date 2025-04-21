package feedbin

import (
	"net/http"
)

// AuthenticationService handles communication with the authentication related
// endpoints of the Feedbin API
type AuthenticationService struct {
	client *Client
}

// Verify checks if the credentials are valid
// It returns true if credentials are valid, false otherwise
func (s *AuthenticationService) Verify() (bool, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "authentication.json", nil)
	if err != nil {
		return false, nil, err
	}

	resp, err := s.client.client.Do(req)
	if err != nil {
		return false, resp, err
	}
	defer resp.Body.Close()

	// 200 OK means credentials are valid
	// 401 Unauthorized means credentials are invalid
	return resp.StatusCode == http.StatusOK, resp, nil
}
