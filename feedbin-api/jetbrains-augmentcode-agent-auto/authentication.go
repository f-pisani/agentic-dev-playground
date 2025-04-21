package feedbin

import (
	"net/http"
)

// AuthenticationService handles communication with the authentication related
// methods of the Feedbin API.
type AuthenticationService struct {
	client *Client
}

// Verify checks if the user credentials are valid.
// It returns nil if the credentials are valid, or an error if they are not.
func (s *AuthenticationService) Verify() error {
	req, err := s.client.NewRequest(http.MethodGet, "/v2/authentication.json", nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return &ErrorResponse{
			Response: resp,
			Message:  "Authentication failed",
		}
	}

	return nil
}
