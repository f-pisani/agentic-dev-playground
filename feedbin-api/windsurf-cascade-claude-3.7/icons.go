package feedbin

import (
	"net/http"
)

// IconsService handles operations related to feed icons
type IconsService struct {
	client *Client
}

// List returns all feed icons for the user's subscriptions
func (s *IconsService) List() ([]Icon, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/icons.json", nil)
	if err != nil {
		return nil, err
	}

	var icons []Icon
	_, err = s.client.Do(req, &icons)
	if err != nil {
		return nil, err
	}

	return icons, nil
}
