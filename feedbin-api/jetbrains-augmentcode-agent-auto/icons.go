package feedbin

import (
	"net/http"
)

// IconsService handles communication with the icons related
// methods of the Feedbin API.
type IconsService struct {
	client *Client
}

// List returns all feed icons.
func (s *IconsService) List() ([]*Icon, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/v2/icons.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var icons []*Icon
	resp, err := s.client.Do(req, &icons)
	if err != nil {
		return nil, resp, err
	}

	return icons, resp, nil
}
